// todo: refactor and clean-up!
// todo: consider using columns storage for operations with columns
// todo: log errors
// todo: return wrapped predefined errors for different actions
// todo: consider adding more diverse errors on integrity_constraint_violation error
// todo: improve work with column position change
// todo: review transaction isolation levels
package postgres

import (
	"database/sql"
	"time"

	"github.com/dnozdrin/detask/internal/app"
	"github.com/dnozdrin/detask/internal/domain/services"

	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// BoardDAO is a data access object for boards
type BoardDAO struct {
	db  *sql.DB
	log app.Logger
}

// NewBoardDAO represents a BoardDAO constructor
func NewBoardDAO(db *sql.DB, log app.Logger) *BoardDAO {
	return &BoardDAO{
		db:  db,
		log: log,
	}
}

// SaveWithDefaultColumn will store the provided board into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (b *BoardDAO) SaveWithDefaultColumn(board *models.Board) (*models.Board, error) {
	empty := &models.Board{}
	if board.ID > 0 {
		return empty, errors.Errorf("can not create a new record with an existing given ID: %d", board.ID)
	}

	tx, err := b.db.Begin()
	if err != nil {
		return nil, err
	}

	defer deferred(b.log, tx.Rollback)
	{
		stmt, err := tx.Prepare(`
		insert into boards (name, description)
		values ($1, $2)
		returning id, created_at, updated_at, name, description;`,
		)
		if err != nil {
			return empty, err
		}

		defer deferred(b.log, stmt.Close)
		if err = stmt.QueryRow(board.Name, board.Description).Scan(
			&board.ID,
			&board.CreatedAt,
			&board.UpdatedAt,
			&board.Name,
			&board.Description,
		); err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
				err = services.ErrRecordAlreadyExist
			}

			return empty, err
		}
	}

	{
		if res, err := tx.Exec(
			`insert into columns (name, board, position) values ('Default', $1, $2);`,
			board.ID,
			services.DefaultColPos,
		); err != nil {
			if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
				err = services.ErrRecordAlreadyExist

				return empty, err
			}
			if res != nil {
				if num, err := res.RowsAffected(); err != nil || num != 1 {
					b.log.Error("default column insertion error: err: %v, affected records: %d", err, num)
					return empty, err
				}
			}

			return empty, err
		}
	}

	if err = tx.Commit(); err != nil {
		return empty, err
	}

	return board, nil
}

// FindOneById will return a pointer to a board with the provided ID or
// a pointer to an empty board and an error
func (b *BoardDAO) FindOneById(ID uint) (*models.Board, error) {
	board := &models.Board{}
	err := b.db.QueryRow(`
		select id, created_at, updated_at, name, description
		from boards
		where id = $1
		order by name
		`, ID).
		Scan(
			&board.ID,
			&board.CreatedAt,
			&board.UpdatedAt,
			&board.Name,
			&board.Description,
		)
	if err == sql.ErrNoRows {
		err = services.ErrRecordNotFound
	}
	return board, err
}

// Find will return all found boards or an error
func (b *BoardDAO) Find() ([]*models.Board, error) {
	boards := make([]*models.Board, 0)

	rows, err := b.db.Query(`
		select id, created_at, updated_at, name, description
		from boards`)
	if err != nil {
		return nil, err
	}
	defer deferred(b.log, rows.Close)

	for rows.Next() {
		board := &models.Board{}
		if err := rows.Scan(
			&board.ID,
			&board.CreatedAt,
			&board.UpdatedAt,
			&board.Name,
			&board.Description,
		); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}

// Update will update the name and description of the persistent representation
// of the board
func (b *BoardDAO) Update(board *models.Board) (*models.Board, error) {
	stmt, err := b.db.Prepare(`
		update boards
		set updated_at = $1, name = $2, description = $3
		where id = $4
		returning id, created_at, updated_at, name, description
	`)
	if err != nil {
		return nil, err
	}
	defer deferred(b.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), board.Name, board.Description, board.ID).Scan(
		&board.ID,
		&board.CreatedAt,
		&board.UpdatedAt,
		&board.Name,
		&board.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		}
		return nil, err
	}

	return board, nil
}

// Delete will delete the record in the database
func (b *BoardDAO) Delete(ID uint) error {
	_, err := b.db.Exec(`delete from boards where id = $1`, ID)
	if err != nil {
		b.log.Warn(err)
		return err
	}

	return nil
}
