package rest

import (
	"encoding/json"
	"github.com/dnozdrin/detask/internal/app/log"
	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/pkg/errors"
)

type TaskHandler struct {
	service TaskService
	log     log.Logger
	router  routeAware
	resp    *responder
}

func NewTaskHandler(service TaskService, logger log.Logger, router routeAware) *TaskHandler {
	return &TaskHandler{
		service: service,
		log:     logger,
		router:  router,
		resp:    &responder{log: logger},
	}
}

func (h TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &task); err != nil {
		h.log.Warnf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	newTask, err := h.service.Create(&task)
	switch {
	case err == nil:
		url, err := h.router.GetURL("get_task", "id", strconv.Itoa(int(newTask.ID)))
		if err != nil {
			h.log.Errorf("unable to build URL: %v", err)
		}
		w.Header().Set("Location", url.Path)
		h.resp.respondJSON(w, http.StatusCreated, newTask)
	case errors.Is(err, services.ErrColumnRelation):
		h.log.Warnf("constraints error: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrRecordAlreadyExist),
		errors.Is(err, services.ErrPositionDuplicate):
		h.log.Warnf("constraints error: %v", err)
		h.resp.respondError(w, http.StatusConflict, err.Error())
	default:
		if _, ok := err.(*v.Errors); ok {
			h.log.Debug("resource was not created", err)
			h.resp.respondJSON(w, http.StatusBadRequest, err)
		} else {
			h.log.Errorf("resource was not created", err)
			h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		}
	}
}

func (h TaskHandler) GetOneById(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Errorf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	task, err := h.service.FindOneById(ID)
	if err != nil {
		if err == services.ErrRecordNotFound {
			h.resp.respondError(w, http.StatusNotFound, "resource was not found")
			return
		}
		h.log.Errorf("resource was not found: %v", err)
	}

	w.Header().Set("Last-Modified", task.UpdatedAt.Format(http.TimeFormat))
	h.resp.respondJSON(w, http.StatusOK, task)
}

func (h TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	demand := make(services.TaskDemand)
	err := parseFilter(r, &demand)
	if err != nil {
		h.log.Debug(err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid filter params")
		return
	}

	tasks, err := h.service.Find(demand)
	if err != nil {
		h.log.Errorf("error while getting records: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	h.resp.respondJSON(w, http.StatusOK, tasks)
}

func (h TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Errorf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid resource identifier")
		return
	}

	var task models.Task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &task); err != nil {
		h.log.Warnf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	task.ID = ID
	updatedTask, err := h.service.Update(&task)
	switch {
	case err == nil:
		h.resp.respondJSON(w, http.StatusOK, updatedTask)
	case errors.Is(err, services.ErrRecordNotFound):
		h.log.Debugf("resource was not found %d", ID)
		h.resp.respondError(w, http.StatusNotFound, "resource was not found")
	case errors.Is(err, services.ErrColumnRelation):
		h.log.Warnf("constraints error: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrPositionDuplicate):
		h.log.Warnf("constraints error: %v", err)
		h.resp.respondError(w, http.StatusConflict, err.Error())
	default:
		if _, ok := err.(*v.Errors); ok {
			h.log.Debugf("resource was not updated: %v", err)
			h.resp.respondJSON(w, http.StatusBadRequest, err)
		} else {
			h.log.Errorf("resource was not updated: %v", err)
			h.resp.respond(w, http.StatusNoContent, "")
		}
	}
}

func (h TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Errorf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid resource identifier")
		return
	}

	if err = h.service.Delete(ID); err != nil {
		h.log.Errorf("error while deleting a record: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
	}

	h.resp.respond(w, http.StatusNoContent, "")
}
