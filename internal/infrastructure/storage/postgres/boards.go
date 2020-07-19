package postgres

import (
	"database/sql"
	"github.com/dnozdrin/detask/internal/app/log"
	"github.com/pkg/errors"
	"time"

	"github.com/dnozdrin/detask/internal/domain/models"
	sv "github.com/dnozdrin/detask/internal/domain/services"
)

// BoardDAO is a data access object for boards
type BoardDAO struct {
	db  querier
	log log.Logger
}

// NewBoardDAO represents a BoardDAO constructor
func NewBoardDAO(db querier, log log.Logger) BoardDAO {
	return BoardDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided board into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (dao BoardDAO) Save(board *models.Board) (*models.Board, error) {
	if board == nil {
		dao.log.Error("boards storage: nil pointer given")
		return nil, errors.New("nil board pointer given")
	}
	if board.ID > 0 {
		dao.log.Warnf("boards storage: %v, ID: %d", sv.ErrRecordAlreadyExist, board.ID)
		return nil, sv.ErrRecordAlreadyExist
	}

	stmt, err := dao.db.Prepare(`
		insert into boards (name, description)
		values ($1, $2)
		returning id, created_at, updated_at, name, description;`,
	)
	if err != nil {
		dao.log.Errorf("boards storage: failed to prepare statement: %v", err)
		return nil, err
	}

	defer deferred(dao.log, stmt.Close)
	if err = stmt.QueryRow(board.Name, board.Description).Scan(
		&board.ID,
		&board.CreatedAt,
		&board.UpdatedAt,
		&board.Name,
		&board.Description,
	); err != nil {
		dao.log.Errorf("boards storage: error while querying a row: %v", err)
		return nil, err
	}

	return board, nil
}

// FindOneById will return a pointer to a board with the provided ID or
// a pointer to an empty board and an error
func (dao BoardDAO) FindOneById(ID uint) (*models.Board, error) {
	board := &models.Board{}
	if err := dao.db.QueryRow(`
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
		); err != nil {
		if err != sql.ErrNoRows {
			dao.log.Errorf("boards storage: error while querying a row: %v", err)
			return nil, err
		}

		return nil, sv.ErrRecordNotFound
	}

	return board, nil
}

// Find will return all found boards or an error
func (dao BoardDAO) Find() ([]*models.Board, error) {
	boards := make([]*models.Board, 0)

	rows, err := dao.db.Query(`select id, created_at, updated_at, name, description from boards`)
	if err != nil {
		dao.log.Errorf("boards storage: error while querying rows: %v", err)
		return nil, err
	}
	defer deferred(dao.log, rows.Close)

	for rows.Next() {
		board := &models.Board{}
		if err := rows.Scan(
			&board.ID,
			&board.CreatedAt,
			&board.UpdatedAt,
			&board.Name,
			&board.Description,
		); err != nil {
			dao.log.Errorf("boards storage: error while querying next row: %v", err)
			return nil, err
		}
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		dao.log.Errorf("boards storage: an error on rows query: %v", err)
		return nil, err
	}

	return boards, nil
}

// Update will update the name and description of the persistent representation
// of the board
func (dao BoardDAO) Update(board *models.Board) (*models.Board, error) {
	if board == nil {
		dao.log.Error("boards storage: nil pointer given")
		return nil, errors.New("nil board pointer given")
	}
	stmt, err := dao.db.Prepare(`
		update boards
		set updated_at = $1, name = $2, description = $3
		where id = $4
		returning id, created_at, updated_at, name, description
	`)
	if err != nil {
		dao.log.Errorf("boards storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), board.Name, board.Description, board.ID).Scan(
		&board.ID,
		&board.CreatedAt,
		&board.UpdatedAt,
		&board.Name,
		&board.Description,
	); err != nil {
		if err != sql.ErrNoRows {
			dao.log.Errorf("boards storage: error while updating a row: %v", err)
			return board, err
		}

		return nil, sv.ErrRecordNotFound
	}

	return board, nil
}

// Delete will delete the record in the database
func (dao BoardDAO) Delete(ID uint) error {
	_, err := dao.db.Exec(`delete from boards where id = $1`, ID)
	if err != nil {
		dao.log.Errorf("boards storage: error while deleting a row: %v", err)
		return err
	}

	return nil
}

// WithTx will return the BoardDAO that will use the provided transaction
func (dao BoardDAO) WithTx(tx *sql.Tx) sv.BoardStorage {
	dao.db = tx
	return dao
}
