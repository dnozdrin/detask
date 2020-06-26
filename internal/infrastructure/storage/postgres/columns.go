package postgres

import (
	"database/sql"
	"github.com/dnozdrin/detask/internal/app"
	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

// ColumnDAO is a data access object for columns
type ColumnDAO struct {
	db  *sql.DB
	log app.Logger
}

// NewColumnDAO represents a ColumnDAO constructor
func NewColumnDAO(db *sql.DB, log app.Logger) *ColumnDAO {
	return &ColumnDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided column into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (c ColumnDAO) Save(column *models.Column) (*models.Column, error) {
	empty := &models.Column{}
	if column.ID > 0 {
		return empty, errors.Errorf("can not create a new record with an existing given ID")
	}

	stmt, err := c.db.Prepare(`
		insert into columns (name, board, position)
		values ($1, $2, $3)
		returning id, created_at, updated_at, name, board, position;`,
	)
	if err != nil {
		c.log.Error(err)
		return empty, err
	}
	defer deferred(c.log, stmt.Close)
	if err = stmt.QueryRow(column.Name, column.BoardID, column.Position).Scan(
		&column.ID,
		&column.CreatedAt,
		&column.UpdatedAt,
		&column.Name,
		&column.BoardID,
		&column.Position,
	); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			err = services.ErrRecordAlreadyExist
		} else {
			c.log.Error(err)
		}

		return empty, err
	}

	return column, nil
}

// FindOneById will return a pointer to a column with the provided ID or
// a pointer to an empty column and an error
func (c ColumnDAO) FindOneById(ID uint) (*models.Column, error) {
	column := &models.Column{}
	err := c.db.QueryRow(`
		select id, created_at, updated_at, name, board, position
		from columns
		where id = $1
		`, ID).
		Scan(
			&column.ID,
			&column.CreatedAt,
			&column.UpdatedAt,
			&column.Name,
			&column.BoardID,
			&column.Position,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		}
		return nil, err
	}

	return column, nil
}

// Find will return all found columns or an error
func (c ColumnDAO) Find() ([]*models.Column, error) {
	columns := make([]*models.Column, 0)

	rows, err := c.db.Query(`
		select id, created_at, updated_at, name, board, position
		from columns
		order by position
		`)
	if err != nil {
		return nil, err
	}
	defer deferred(c.log, rows.Close)

	for rows.Next() {
		column := &models.Column{}
		if err := rows.Scan(
			&column.ID,
			&column.CreatedAt,
			&column.UpdatedAt,
			&column.Name,
			&column.BoardID,
			&column.Position,
		); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

// Update will update the name of the persistent representation
// of the column. Returns pointer to a updated column or to a empty column
// entity and an error
func (c ColumnDAO) Update(column *models.Column) (*models.Column, error) {
	stmt, err := c.db.Prepare(`
		update columns
		set updated_at = $1, name = $2
		where id = $3
		returning id, created_at, updated_at, name, board, position
	`)
	if err != nil {
		return nil, err
	}
	defer deferred(c.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), column.Name, column.ID).Scan(
		&column.ID,
		&column.CreatedAt,
		&column.UpdatedAt,
		&column.Name,
		&column.BoardID,
		&column.Position,
	); err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		}
		return nil, err
	}

	return column, nil
}

// Delete will the column with the provided ID. The last column cannot be deleted.
// When a column is deleted, its tasks are moved to the column to the left of the
// current or to the right of the current if the curring is the leftmost
func (c ColumnDAO) Delete(ID uint) error {
	tx, err := c.db.Begin()
	if err != nil {
		c.log.Error(err)
		return err
	}

	{
		var (
			position     float64
			num, boardID uint
		)

		err = tx.QueryRow(`
		select count(c1.id), c1.board, c1.position
		from "columns" c1 left join "columns" c2 on c1.board = c2.board
		where c1.id = $1
		group by c1.board, c1.position`, ID).Scan(&num, &boardID, &position)
		if err != nil {
			c.log.Error(err)
			return err
		}

		if num <= 1 {
			return services.ErrLastColumn
		}

	}

	{
		var (
			prev, next sql.NullInt64
			target     uint
		)
		err = tx.QueryRow(`
		select prev, next from
			(select
			        id,
			        lag(id) over (order by position) as prev,
			        lead(id) over (order by position) as next
			from "columns" ) sub
		where id = $1`, ID).Scan(&prev, &next)
		if err != nil {
			c.log.Error(err)
			return err
		}

		if prev.Valid && prev.Int64 > 0 {
			target = uint(prev.Int64)
		} else if next.Valid && next.Int64 > 0 {
			target = uint(next.Int64)
		} else {
			err = errors.Errorf("target column for tasks transfer not found, column ID: %d", ID)
			c.log.Error(err)
			return err
		}

		_, err = tx.Exec(`update tasks set "column" = $1 where "column" = $2`, target, ID)
		if err != nil {
			c.log.Error(err)
			return err
		}
	}

	{
		res, err := tx.Exec(`delete from "columns" where id = $1`, ID)
		if err != nil {
			c.log.Error(err)
			return err
		}

		rowsNum, err := res.RowsAffected()
		if err != nil {
			c.log.Error(err)
			return err
		}

		if rowsNum != 1 {
			err = errors.Errorf("tried to delete %d rows, want 1", rowsNum)
			c.log.Error(err)
			return err
		}
	}

	return tx.Commit()
}
