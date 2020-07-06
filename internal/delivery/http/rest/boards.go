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

// BoardHandler provides a Rest API http handlers for work with boards
type BoardHandler struct {
	service BoardService
	log     log.Logger
	router  routeAware
	resp    *responder
}

// NewBoardHandler is BoardHandler constructor
func NewBoardHandler(service BoardService, logger log.Logger, router routeAware) *BoardHandler {
	return &BoardHandler{
		service: service,
		log:     logger,
		router:  router,
		resp:    &responder{log: logger},
	}
}

// Create will call creation of the provided resource
func (h BoardHandler) Create(w http.ResponseWriter, r *http.Request) {
	var board models.Board
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &board); err != nil {
		h.log.Debugf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	newBoard, err := h.service.Create(&board)
	switch {
	case err == nil:
		url, err := h.router.GetURL("get_board", "id", strconv.Itoa(int(newBoard.ID)))
		if err != nil {
			h.log.Errorf("unable to build URL: %v", err)
		}
		w.Header().Set("Location", url.Path)
		h.resp.respondJSON(w, http.StatusCreated, newBoard)
	case errors.Is(err, services.ErrRecordAlreadyExist):
		h.log.Debugf("constraints error: %v", err)
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

// GetOneById will respond with the requested resource or an error
func (h BoardHandler) GetOneById(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Errorf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	board, err := h.service.FindOneById(ID)
	if err != nil {
		if err == services.ErrRecordNotFound {
			h.resp.respondError(w, http.StatusNotFound, "resource was not found")
			return
		}
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	w.Header().Set("Last-Modified", board.UpdatedAt.Format(http.TimeFormat))
	h.resp.respondJSON(w, http.StatusOK, board)
}

// Get will respond with the requested resources or an error
func (h BoardHandler) Get(w http.ResponseWriter, r *http.Request) {
	boards, err := h.service.Find()
	if err != nil {
		h.log.Errorf("error while getting records: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	h.resp.respondJSON(w, http.StatusOK, boards)
}

// Update will trigger update of the provided resource
func (h BoardHandler) Update(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Errorf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	var board models.Board
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &board); err != nil {
		h.log.Debugf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	board.ID = ID
	updatedBoard, err := h.service.Update(&board)
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
			h.resp.respondError(w, http.StatusInternalServerError, "internal error")
		}
	}
}

// Delete will trigger deletion of the provided resource
func (h BoardHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Errorf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	if err = h.service.Delete(ID); err != nil {
		h.log.Errorf("error while deleting a record: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
	}

	h.resp.respond(w, http.StatusNoContent, "")
}
