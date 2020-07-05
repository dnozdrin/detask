package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dnozdrin/detask/internal/app/log"
	"time"

	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// CommentsDAO is a data access object for comments
type CommentsDAO struct {
	db  db
	log log.Logger
}

// NewCommentsDAO represents a CommentsDAO constructor
func NewCommentsDAO(db db, log log.Logger) *CommentsDAO {
	return &CommentsDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided comment into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (dao CommentsDAO) Save(comment *models.Comment) (*models.Comment, error) {
	if comment == nil {
		dao.log.Error("comments storage: nil pointer given")
		return nil, errors.New("nil comment pointer given")
	}
	if comment.ID > 0 {
		dao.log.Warnf("comments storage: %v, ID: %d", services.ErrRecordAlreadyExist, comment.ID)
		return nil, services.ErrRecordAlreadyExist
	}

	stmt, err := dao.db.Prepare(`
		insert into comments (text, task)
		values ($1, $2)
		returning id, created_at, updated_at, text, task;`,
	)
	if err != nil {
		dao.log.Errorf("comments storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
	if err = stmt.QueryRow(comment.Text, comment.TaskID).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Text,
		&comment.TaskID,
	); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			switch pgErr.Constraint {
			case "comments_task_fkey":
				err = services.ErrTaskRelation
			default:
				dao.log.Errorf("comments storage: integrity constraint violation: %v", err)
			}
		} else {
			dao.log.Errorf("comments storage: error while querying a row: %v", err)
		}

		return nil, err
	}

	return comment, nil
}

// FindOneById will return a pointer to a comment with the provided ID or
// a pointer to an empty comment and an error
func (dao CommentsDAO) FindOneById(ID uint) (*models.Comment, error) {
	comment := &models.Comment{}
	err := dao.db.QueryRow(`
		select id, created_at, updated_at, text, task
		from comments
		where id = $1
		`, ID).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt, &comment.Text, &comment.TaskID)
	if err != nil {
		if err != sql.ErrNoRows {
			dao.log.Errorf("comments storage: error while querying a row: %v", err)
			return nil, err
		}
		return nil, services.ErrRecordNotFound
	}

	return comment, err
}

// Find will return all found comments that meet the provided demand or an error
func (dao CommentsDAO) Find(demand services.CommentDemand) ([]*models.Comment, error) {
	const querySelect = "id, created_at, updated_at, text, task"
	where := "1=1"
	if taskID, ok := demand["task"]; ok {
		where = where + fmt.Sprintf(" and t.task = %d", taskID)
	}

	rows, err := dao.db.Query(
		fmt.Sprintf(`select %s from comments t where %s order by created_at desc;`, querySelect, where),
	)
	if err != nil {
		dao.log.Errorf("comments storage: error while querying rows: %v", err)
		return nil, err
	}
	defer deferred(dao.log, rows.Close)

	comments := make([]*models.Comment, 0)
	for rows.Next() {
		comment := &models.Comment{}
		if err := rows.Scan(
			&comment.ID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.Text,
			&comment.TaskID,
		); err != nil {
			dao.log.Errorf("comments storage: error while querying next row: %v", err)
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		dao.log.Errorf("comments storage: error while querying rows: %v", err)
		return nil, err
	}

	return comments, nil
}

// Update will update text of the persistent representation of the comment
func (dao CommentsDAO) Update(comment *models.Comment) (*models.Comment, error) {
	if comment == nil {
		dao.log.Error("comments storage: nil pointer given")
		return nil, errors.New("nil comment pointer given")
	}
	stmt, err := dao.db.Prepare(`
		update comments
		set updated_at = $1, text = $2
		where id = $3
		returning id, created_at, updated_at, text, task
	`)
	if err != nil {
		dao.log.Errorf("comments storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), comment.Text, comment.ID).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Text,
		&comment.TaskID,
	); err != nil {
		if err != sql.ErrNoRows {
			dao.log.Errorf("comments storage: error while updating a row: %v", err)
			return nil, err
		}

		return nil, services.ErrRecordNotFound
	}

	return comment, nil
}

// Delete will delete the record in the database
func (dao CommentsDAO) Delete(ID uint) error {
	_, err := dao.db.Exec("delete from comments where id = $1", ID)
	if err != nil {
		dao.log.Errorf("comments storage: error while deleting a row: %v", err)
		return err
	}

	return err
}
