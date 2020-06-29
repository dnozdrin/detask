package rest

import (
	"encoding/json"
	 "github.com/dnozdrin/detask/internal/app/log"
	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ColumnHandler struct {
	service *services.ColumnService
	log     log.Logger
	router  routeAware
	resp    *responder
}

func NewColumnHandler(service *services.ColumnService, logger log.Logger, router routeAware) *ColumnHandler {
	return &ColumnHandler{
		service: service,
		log:     logger,
		router:  router,
		resp:    &responder{log: logger},
	}
}

func (h ColumnHandler) Create(w http.ResponseWriter, r *http.Request) {
	var column models.Column
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &column); err != nil {
		h.log.Warnf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	newColumn, err := h.service.Create(&column)
	switch {
	case err == nil:
		url, err := h.router.GetURL("get_column", "id", strconv.Itoa(int(newColumn.ID)))
		if err != nil {
			h.log.Errorf("unable to build URL: %v", err)
		}
		w.Header().Set("Location", url.Path)
		h.resp.respondJSON(w, http.StatusCreated, newColumn)
	case errors.Is(err, services.ErrRecordAlreadyExist):
		h.log.Errorf("given resource already exists", err)
		h.resp.respondError(w, http.StatusBadRequest, "given resource already exists")
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

func (h ColumnHandler) GetOneById(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Warnf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	column, err := h.service.FindOneById(ID)
	if err != nil {
		if err == services.ErrRecordNotFound {
			h.resp.respondError(w, http.StatusNotFound, "resource was not found")
			return
		}
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	w.Header().Set("Last-Modified", column.UpdatedAt.Format(http.TimeFormat))
	h.resp.respondJSON(w, http.StatusOK, column)
}

func (h ColumnHandler) Get(w http.ResponseWriter, r *http.Request) {
	boards, err := h.service.Find()
	if err != nil {
		h.log.Errorf("error while getting records: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	h.resp.respondJSON(w, http.StatusOK, boards)
}

func (h ColumnHandler) Update(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Warnf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid resource identifier")
		return
	}

	var column models.Column
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &column); err != nil {
		h.log.Warnf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	column.ID = ID
	updatedBoard, err := h.service.Update(&column)
	switch err {
	case nil:
		h.resp.respondJSON(w, http.StatusOK, updatedBoard)
	case services.ErrRecordNotFound:
		h.log.Debugf("resource was not found %d", ID)
		h.resp.respondError(w, http.StatusNotFound, "resource was not found")
	default:
		if _, ok := err.(*v.Errors); ok {
			h.log.Debugf("resource was not updated: %v", err)
			h.resp.respondJSON(w, http.StatusBadRequest, err)
		} else {
			h.log.Errorf("resource was not updated: %v", err)
			h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		}
	}
}

func (h ColumnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Warnf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid resource identifier")
		return
	}

	if err = h.service.Delete(ID); err != nil {
		h.log.Errorf("error while deleting a record: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	h.resp.respond(w, http.StatusNoContent, "")
}
