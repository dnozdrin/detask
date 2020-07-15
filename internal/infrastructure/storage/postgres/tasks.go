package postgres

import (
	"database/sql"
	"fmt"
	"github.com/dnozdrin/detask/internal/app/log"
	"github.com/pkg/errors"
	"time"

	sv "github.com/dnozdrin/detask/internal/domain/services"
	"github.com/lib/pq"

	"github.com/dnozdrin/detask/internal/domain/models"
)

// TaskDAO is a data access object for boards
type TaskDAO struct {
	db  querier
	log log.Logger
}

// NewTaskDAO represents a TaskDAO constructor
func NewTaskDAO(db querier, log log.Logger) TaskDAO {
	return TaskDAO{
		db:  db,
		log: log,
	}
}

// Save will store the provided task into the database and return
// a pointer to the saved entity. Returns nil and an error in case of error.
func (dao TaskDAO) Save(task *models.Task) (*models.Task, error) {
	if task == nil {
		dao.log.Error("tasks storage: nil pointer given")
		return nil, errors.New("nil tasks pointer given")
	}
	if task.ID > 0 {
		dao.log.Warnf("tasks storage: %v, ID: %d", sv.ErrRecordAlreadyExist, task.ID)
		return nil, sv.ErrRecordAlreadyExist
	}

	stmt, err := dao.db.Prepare(`
		insert into tasks (name, description, "column", position)
		values ($1, $2, $3, $4)
		returning id, created_at, updated_at, name, description, "column", position;`,
	)
	if err != nil {
		dao.log.Errorf("tasks storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
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
				err = sv.ErrColumnRelation
			case "tasks_position_column_key":
				err = sv.ErrPositionDuplicate
			default:
				dao.log.Errorf("tasks storage: integrity constraint violation: %v", err)
			}
		} else {
			dao.log.Errorf("tasks storage: error while querying a row: %v", err)
		}

		return nil, err
	}

	return task, nil
}

// FindOneById will return a pointer to a task with the provided ID or
// a pointer to an empty task and an error
func (dao TaskDAO) FindOneById(ID uint) (*models.Task, error) {
	task := &models.Task{}
	err := dao.db.QueryRow(`
		select id, created_at, updated_at, name, description, "column", position
		from tasks
		where id = $1
		`, ID).
		Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.Name, &task.Description, &task.ColumnID, &task.Position)
	if err != nil {
		if err != sql.ErrNoRows {
			dao.log.Errorf("tasks storage: error while querying a row: %v", err)
			return nil, err
		}

		return nil, sv.ErrRecordNotFound
	}

	return task, err
}

// Find will return all found tasks that meet the provided demand or an error
func (dao TaskDAO) Find(demand sv.TaskDemand) ([]*models.Task, error) {
	tasks := make([]*models.Task, 0)

	const querySelect = "t.id, t.created_at, t.updated_at, t.name, t.description, t.column, t.position"
	var join, where string

	where = "1=1"
	if boardID, ok := demand["board"]; ok {
		join = `join "columns" c on t.column = c.id`
		where = where + fmt.Sprintf(" and c.board = %d", boardID)
	}
	if columnID, ok := demand["column"]; ok {
		where = where + fmt.Sprintf(" and t.column = %d", columnID)
	}

	rows, err := dao.db.Query(fmt.Sprintf(`select %s from tasks t %s where %s order by position;`, querySelect, join, where))
	if err != nil {
		dao.log.Errorf("tasks storage: error while querying rows: %v", err)
		return nil, err
	}

	defer deferred(dao.log, rows.Close)

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
		dao.log.Errorf("tasks storage: rows query error: %v", err)
		return nil, err
	}

	return tasks, nil
}

// Update will update text of the persistent representation of the task
func (dao TaskDAO) Update(task *models.Task) (*models.Task, error) {
	if task == nil {
		dao.log.Error("tasks storage: nil pointer given")
		return nil, errors.New("nil tasks pointer given")
	}
	stmt, err := dao.db.Prepare(`
		update tasks
		set updated_at = $1, name = $2, description = $3, position = $4, "column" = $5
		where id = $6
		returning id, created_at, updated_at, name, description, "column", position
	`)
	if err != nil {
		dao.log.Errorf("tasks storage: failed to prepare statement: %v", err)
		return nil, err
	}
	defer deferred(dao.log, stmt.Close)
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
			err = sv.ErrRecordNotFound
		} else if pgErr, ok := err.(*pq.Error); ok && pgErr.Code.Class().Name() == "integrity_constraint_violation" {
			switch pgErr.Constraint {
			case "tasks_column_fkey":
				err = sv.ErrColumnRelation
			case "tasks_position_column_key":
				err = sv.ErrPositionDuplicate
			default:
				dao.log.Errorf("tasks storage: integrity constraint violation: %v", err)
			}
		} else {
			dao.log.Errorf("tasks storage: error while updating a row: %v", err)
		}

		return nil, err
	}

	return task, nil
}

// MoveToColumn will move all tasks from source column to target column
func (dao TaskDAO) MoveToColumn(sourceID, targetID uint) error {
	if _, err := dao.db.Exec(`update tasks set "column" = $1 where "column" = $2`, targetID, sourceID); err != nil {
		dao.log.Errorf(
			"tasks storage: error while moving tasks from column %d to column %d: %v",
			sourceID,
			targetID,
			err,
		)
		return err
	}

	return nil
}

// Delete will delete the record in the database
func (dao TaskDAO) Delete(ID uint) error {
	if _, err := dao.db.Exec("delete from tasks where id = $1", ID); err != nil {
		dao.log.Errorf("tasks storage: error while deleting a row: %v", err)
		return err
	}

	return nil
}

// WithTx will return the TaskDAO that will use the provided transaction
func (dao TaskDAO) WithTx(tx *sql.Tx) sv.TaskStorage {
	dao.db = tx
	return dao
}
