package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dnozdrin/detask/internal/app"
	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// CommentsDAO is a data access object for comments
type CommentsDAO struct {
	db  *sql.DB
	log app.Logger
}

// NewCommentsDAO represents a CommentsDAO constructor
func NewCommentsDAO(db *sql.DB, log app.Logger) *CommentsDAO {
	return &CommentsDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided comment into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (c CommentsDAO) Save(comment *models.Comment) (*models.Comment, error) {
	empty := &models.Comment{}
	if comment.ID > 0 {
		return empty, errors.Errorf("can not create a new record with an existing given ID")
	}

	stmt, err := c.db.Prepare(`
		insert into comments (text, task)
		values ($1, $2)
		returning id, created_at, updated_at, text, task;`,
	)
	if err != nil {
		c.log.Error(err)
		return empty, err
	}
	defer deferred(c.log, stmt.Close)
	if err = stmt.QueryRow(comment.Text, comment.TaskID).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Text,
		&comment.TaskID,
	); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			err = services.ErrRecordAlreadyExist
		} else {
			c.log.Error(err)
		}

		return empty, err
	}

	return comment, nil
}

// FindOneById will return a pointer to a comment with the provided ID or
// a pointer to an empty comment and an error
func (c CommentsDAO) FindOneById(ID uint) (*models.Comment, error) {
	comment := &models.Comment{}
	err := c.db.QueryRow(`
		select id, created_at, updated_at, text, task
		from comments
		where id = $1
		`, ID).
		Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt, &comment.Text, &comment.TaskID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		}
		return nil, err
	}

	return comment, err
}

// Find will return all found comments that meet the provided demand or an error
func (c CommentsDAO) Find(demand services.CommentDemand) ([]*models.Comment, error) {
	const querySelect = "id, created_at, updated_at, text, task"
	where := "1=1"
	if taskID, ok := demand["task"]; ok {
		where = where + fmt.Sprintf(" and t.column = %d", taskID)
	}

	rows, err := c.db.Query(
		fmt.Sprintf(`select %s from comments t where %s order by created_at desc;`, querySelect, where),
	)
	if err != nil {
		return nil, err
	}
	defer deferred(c.log, rows.Close)

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
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// Update will update text of the persistent representation of the comment
func (c CommentsDAO) Update(comment *models.Comment) (*models.Comment, error) {
	stmt, err := c.db.Prepare(`
		update comments
		set updated_at = $1, text = $2
		where id = $3
		returning id, created_at, updated_at, text, task
	`)
	if err != nil {
		return nil, err
	}
	defer deferred(c.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), comment.Text, comment.ID).Scan(
		&comment.ID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.Text,
		&comment.TaskID,
	); err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		}
		return nil, err
	}

	return comment, nil
}

// Delete will delete the record in the database
func (c CommentsDAO) Delete(ID uint) error {
	_, err := c.db.Exec("delete from comments where id = $1", ID)

	return err
}
