package rest

import (
	"encoding/json"
	"github.com/dnozdrin/detask/internal/app/log"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dnozdrin/detask/internal/domain/models"
	"github.com/dnozdrin/detask/internal/domain/services"
	v "github.com/dnozdrin/detask/internal/domain/validation"
	"github.com/pkg/errors"
)

type CommentHandler struct {
	service CommentService
	log     log.Logger
	router  routeAware
	resp    *responder
}

func NewCommentHandler(service CommentService, logger log.Logger, router routeAware) *CommentHandler {
	return &CommentHandler{
		service: service,
		log:     logger,
		router:  router,
		resp:    &responder{log: logger},
	}
}

func (h CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &comment); err != nil {
		h.log.Warnf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	newComment, err := h.service.Create(&comment)
	switch {
	case err == nil:
		url, err := h.router.GetURL("get_comment", "id", strconv.Itoa(int(newComment.ID)))
		if err != nil {
			h.log.Errorf("unable to build URL: %v", err)
		}
		w.Header().Set("Location", url.Path)
		h.resp.respondJSON(w, http.StatusCreated, newComment)
	case errors.Is(err, services.ErrTaskRelation):
		h.log.Warnf("constraints error: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, services.ErrRecordAlreadyExist):
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

func (h CommentHandler) GetOneById(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Warnf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	comment, err := h.service.FindOneById(ID)
	if err != nil {
		if err == services.ErrRecordNotFound {
			h.resp.respondError(w, http.StatusNotFound, "resource was not found")
			return
		}
		h.resp.respondError(w, http.StatusInternalServerError, "invalid resource identifier")
		return
	}

	w.Header().Set("Last-Modified", comment.UpdatedAt.Format(http.TimeFormat))
	h.resp.respondJSON(w, http.StatusOK, comment)
}

func (h CommentHandler) Get(w http.ResponseWriter, r *http.Request) {
	demand := make(services.CommentDemand)
	err := parseFilter(r, demand)
	if err != nil {
		h.log.Debug(err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid filter params")
		return
	}

	boards, err := h.service.Find(demand)
	if err != nil {
		h.log.Errorf("error while getting records: %v", err)
		h.resp.respondError(w, http.StatusInternalServerError, errInternalServer)
		return
	}

	h.resp.respondJSON(w, http.StatusOK, boards)
}

func (h CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	ID, err := h.router.GetIDVar(r)
	if err != nil {
		h.log.Warnf("error on parsing resource identifier: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "invalid resource identifier")
		return
	}

	var comment models.Comment
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error on request body read: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, "error on request body read")
		return
	}
	if err := json.Unmarshal(reqBody, &comment); err != nil {
		h.log.Warnf("error on request body parsing: %v", err)
		h.resp.respondError(w, http.StatusBadRequest, errInvalidJSON)
		return
	}

	comment.ID = ID
	updatedBoard, err := h.service.Update(&comment)
	switch  {
	case err == nil:
		h.resp.respondJSON(w, http.StatusOK, updatedBoard)
	case errors.Is(err, services.ErrRecordNotFound):
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

func (h CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
