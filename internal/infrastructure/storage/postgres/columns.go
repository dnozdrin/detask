package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dnozdrin/detask/internal/app/log"
	"github.com/dnozdrin/detask/internal/domain/models"
	sv "github.com/dnozdrin/detask/internal/domain/services"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"time"
)

// ColumnDAO is a data access object for columns
type ColumnDAO struct {
	db  querier
	log log.Logger
}

// NewColumnDAO represents a ColumnDAO constructor
func NewColumnDAO(db querier, log log.Logger) ColumnDAO {
	return ColumnDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided column into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (dao ColumnDAO) Save(column *models.Column) (*models.Column, error) {
	if column == nil {
		dao.log.Error("columns storage: nil pointer given")
		return nil, errors.New("nil column pointer given")
	}
	if column.ID > 0 {
		dao.log.Warnf("columns storage: %v, ID: %d", sv.ErrRecordAlreadyExist, column.ID)
		return nil, sv.ErrRecordAlreadyExist
	}

	stmt, err := dao.db.Prepare(`
		insert into columns (name, board, position)
		values ($1, $2, $3)
		returning id, created_at, updated_at, name, board, position;`,
	)
	if err != nil {
		dao.log.Errorf("columns storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
	if err = stmt.QueryRow(column.Name, column.BoardID, column.Position).Scan(
		&column.ID,
		&column.CreatedAt,
		&column.UpdatedAt,
		&column.Name,
		&column.BoardID,
		&column.Position,
	); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			switch pgErr.Constraint {
			case "columns_name_board_key":
				err = sv.ErrNameDuplicate
			case "columns_position_board_key":
				err = sv.ErrPositionDuplicate
			case "columns_board_fkey":
				err = sv.ErrBoardRelation
			default:
				dao.log.Errorf("columns storage: integrity constraint violation: %v", err)
			}
		} else {
			dao.log.Errorf("columns storage: error while querying a row: %v", err)
		}

		return nil, err
	}

	return column, nil
}

// FindOneById will return a pointer to a column with the provided ID or
// a pointer to an empty column and an error
func (dao ColumnDAO) FindOneById(ID uint) (*models.Column, error) {
	column := &models.Column{}
	err := dao.db.QueryRow(`
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
		if err != sql.ErrNoRows {
			dao.log.Errorf("columns storage: error while querying a row: %v", err)
			return nil, err
		}

		return nil, sv.ErrRecordNotFound
	}

	return column, nil
}

// Find will return all found columns or an error
func (dao ColumnDAO) Find(demand sv.ColumnDemand) ([]*models.Column, error) {
	const querySelect = "id, created_at, updated_at, name, board, position"
	columns := make([]*models.Column, 0)
	where := "1=1"
	if taskID, ok := demand["board"]; ok {
		where = where + fmt.Sprintf(" and board = %d", taskID)
	}

	rows, err := dao.db.Query(fmt.Sprintf(`select %s from columns where %s order by position;`, querySelect, where))
	if err != nil {
		dao.log.Errorf("columns storage: error while querying rows: %v", err)
		return nil, err
	}
	defer deferred(dao.log, rows.Close)

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
			dao.log.Errorf("columns storage: error while querying next row: %v", err)
			return nil, err
		}
		columns = append(columns, column)
	}

	if err := rows.Err(); err != nil {
		dao.log.Errorf("columns storage: rows query error: %v", err)
		return nil, err
	}

	return columns, nil
}

// Update will update the name of the persistent representation
// of the column. Returns pointer to a updated column or to a empty column
// entity and an error
func (dao ColumnDAO) Update(column *models.Column) (*models.Column, error) {
	if column == nil {
		dao.log.Error("columns storage: nil pointer given")
		return nil, errors.New("nil column pointer given")
	}
	stmt, err := dao.db.Prepare(`
		update columns
		set updated_at = $1, name = $2, position = $3
		where id = $4
		returning id, created_at, updated_at, name, board, position
	`)
	if err != nil {
		dao.log.Errorf("columns storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), column.Name, column.Position, column.ID).Scan(
		&column.ID,
		&column.CreatedAt,
		&column.UpdatedAt,
		&column.Name,
		&column.BoardID,
		&column.Position,
	); err != nil {
		if err == sql.ErrNoRows {
			err = sv.ErrRecordNotFound
		} else if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			switch pgErr.Constraint {
			case "columns_name_board_key":
				err = sv.ErrNameDuplicate
			case "columns_position_board_key":
				err = sv.ErrPositionDuplicate
			default:
				dao.log.Errorf("columns storage: integrity constraint violation: %v", err)
			}
		} else {
			dao.log.Errorf("columns storage: error while updating a row: %v", err)
		}

		return nil, err
	}

	return column, nil
}

// Delete will the column with the provided ID.
func (dao ColumnDAO) Delete(ID uint) error {
	res, err := dao.db.Exec(`delete from "columns" where id = $1`, ID)
	if err != nil {
		dao.log.Errorf("columns storage: error while deleting a column ID: %d: %v", ID, err)
		return err
	}

	rowsNum, err := res.RowsAffected()
	if err != nil {
		dao.log.Error(err)
		return err
	}

	if rowsNum != 1 {
		err = errors.Errorf("tried to delete %d rows, want 1", rowsNum)
		dao.log.Error(err)
		return err
	}

	return nil
}

// WithTx will return the ColumnDAO that will use the provided transaction
func (dao ColumnDAO) WithTx(tx *sql.Tx) sv.ColumnStorage {
	dao.db = tx
	return dao
}

// FindLeftColumn will find a column to the left of the one with the provided ID
func (dao ColumnDAO) CountColumnsByBoard(ID uint) (int, error) {
	var num int
	if err := dao.db.QueryRow(`select count(c1.id) from "columns" c1 where c1.board = $1`, ID).
		Scan(&num); err != nil {
		dao.log.Errorf("columns storage: error while counting columns by board: %v", err)
		return 0, err
	}

	return num, nil
}

func (dao ColumnDAO) FindColumnToTheLeft(ID uint) (uint, error) {
	var prev sql.NullInt64
	if err := dao.db.QueryRow(`
		select prev
		from (select id, lag(id) over (order by position) as prev from "columns") sub
		where id = $1`, ID).Scan(&prev); err != nil {
		dao.log.Errorf("columns storage: error while querying prev record: %v", err)
		return 0, err
	}

	if !prev.Valid || prev.Int64 <= 0 {
		err := errors.New("columns storage: invalid left column record")
		dao.log.Errorf("%v: %v", err, prev)
		return 0, err
	}

	return uint(prev.Int64), nil
}

func (dao ColumnDAO) FindColumnToTheRight(ID uint) (uint, error) {
	var next sql.NullInt64
	if err := dao.db.QueryRow(`
		select next
		from (select id, lead(id) over (order by position) as next from "columns") sub
		where id = $1`, ID).Scan(&next); err != nil {
		dao.log.Errorf("columns storage: error while querying prev record: %v", err)
		return 0, err
	}

	if !next.Valid || next.Int64 <= 0 {
		err := errors.New("columns storage: invalid right column record")
		dao.log.Errorf("%v: %v", err, next)
		return 0, err
	}

	return uint(next.Int64), nil
}
