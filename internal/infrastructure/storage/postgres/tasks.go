package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dnozdrin/detask/internal/app/log"
	"time"

	"github.com/dnozdrin/detask/internal/domain/services"
	"github.com/lib/pq"

	"github.com/dnozdrin/detask/internal/domain/models"
)

// TaskDAO is a data access object for boards
type TaskDAO struct {
	db  *sql.DB
	log log.Logger
}

// NewTaskDAO represents a TaskDAO constructor
func NewTaskDAO(db *sql.DB, log log.Logger) *TaskDAO {
	return &TaskDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided task into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (t TaskDAO) Save(task *models.Task) (*models.Task, error) {
	empty := &models.Task{}
	if task.ID > 0 {
		return empty, services.ErrRecordAlreadyExist
	}

	stmt, err := t.db.Prepare(`
		insert into tasks (name, description, "column", position)
		values ($1, $2, $3, $4)
		returning id, created_at, updated_at, name, description, "column", position;`,
	)
	if err != nil {
		t.log.Error(err)
		return empty, err
	}
	defer deferred(t.log, stmt.Close)
	if err = stmt.QueryRow(task.Name, task.Description, task.ColumnID, task.Position).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Name,
		&task.Description,
		&task.ColumnID,
		&task.Position,
	); err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			switch pgErr.Constraint {
			case "tasks_column_fkey":
				err = services.ErrColumnRelation
			case "tasks_position_column_key":
				err = services.ErrPositionDuplicate
			default:
				t.log.Error(err)
			}
		} else {
			t.log.Error(err)
		}

		return empty, err
	}

	return task, nil
}

// FindOneById will return a pointer to a task with the provided ID or
// a pointer to an empty task and an error
func (t TaskDAO) FindOneById(ID uint) (*models.Task, error) {
	task := &models.Task{}
	err := t.db.QueryRow(`
		select id, created_at, updated_at, name, description, "column", position
		from tasks
		where id = $1
		`, ID).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.ColumnID, &task.Position)
	if err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		}
		return nil, err
	}

	return task, err
}

// Find will return all found tasks that meet the provided demand or an error
func (t TaskDAO) Find(demand services.TaskDemand) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	const querySelect = "t.id, t.created_at, t.updated_at, t.name, t.description, t.column, t.position"
	var join, where string

	where = "1=1"
	if boardID, ok := demand["board"]; ok {
		join = "join \"columns\" c on t.column = c.id"
		where = where + fmt.Sprintf(" and c.board = %d", boardID)
	}
	if columnID, ok := demand["column"]; ok {
		where = where + fmt.Sprintf(" and t.column = %d", columnID)
	}

	rows, err := t.db.Query(fmt.Sprintf(`select %s from tasks t %s where %s order by position;`, querySelect, join, where))
	if err != nil {
		return nil, err
	}

	defer deferred(t.log, rows.Close)

	for rows.Next() {
		task := &models.Task{}
		if err := rows.Scan(
			&task.ID,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Name,
			&task.Description,
			&task.ColumnID,
			&task.Position,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Update will update text of the persistent representation of the task
func (t TaskDAO) Update(task *models.Task) (*models.Task, error) {
	stmt, err := t.db.Prepare(`
		update tasks
		set updated_at = $1, name = $2, description = $3, position = $4, "column" = $5
		where id = $6
		returning id, created_at, updated_at, name, description, "column", position
	`)
	if err != nil {
		return nil, err
	}
	defer deferred(t.log, stmt.Close)
	if err = stmt.QueryRow(time.Now(), task.Name, task.Description, task.Position, task.ColumnID, task.ID).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Name,
		&task.Description,
		&task.ColumnID,
		&task.Position,
	); err != nil {
		if err == sql.ErrNoRows {
			err = services.ErrRecordNotFound
		} else if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			switch pgErr.Constraint {
			case "tasks_column_fkey":
				err = services.ErrColumnRelation
			case "tasks_position_column_key":
				err = services.ErrPositionDuplicate
			default:
				t.log.Error(err)
			}
		} else {
			t.log.Error(err)
		}

		return nil, err
	}

	return task, nil
}

// Delete will delete the record in the database
func (t TaskDAO) Delete(ID uint) error {
	_, err := t.db.Exec("delete from tasks where id = $1", ID)

	return err
}
