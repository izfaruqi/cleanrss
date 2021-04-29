package http

import (
	"cleanrss/domain"
	utils_http "cleanrss/utils/http"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

type cleanerHTTPHandler struct {
	u domain.CleanerUsecase
}

func NewCleanerHTTPHandler(u domain.CleanerUsecase) http.Handler {
	handler := cleanerHTTPHandler{u: u}
	router := chi.NewRouter()
	router.Get("/", handler.getAll)
	router.Get("/{id}", handler.getById)
	router.Post("/", handler.insert)
	router.Post("/{id}", handler.update)
	router.Delete("/{id}", handler.delete)
	router.Get("/clean/{id}", handler.cleanPage)
	return router
}

func (h cleanerHTTPHandler) cleanPage(w http.ResponseWriter, r *http.Request) {
	entryId, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	cleaned, err := h.u.GetCleanedEntry(entryId)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(cleaned))
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	return
}

func (h cleanerHTTPHandler) getAll(w http.ResponseWriter, r *http.Request) {
	cleaners, err := h.u.GetAll()
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, cleaners)
	return
}

func (h cleanerHTTPHandler) getById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	cleaner, err := h.u.GetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils_http.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	utils_http.WriteJson(w, cleaner)
	return
}

func (h cleanerHTTPHandler) insert(w http.ResponseWriter, r *http.Request) {
	var cleaner domain.Cleaner
	err := json.NewDecoder(r.Body).Decode(&cleaner)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.u.Insert(&cleaner)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]int64{"id": cleaner.Id})
	return
}

func (h cleanerHTTPHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, errors.New("id is invalid"))
		return
	}
	var cleaner domain.Cleaner
	err = json.NewDecoder(r.Body).Decode(&cleaner)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	cleaner.Id = id
	err = h.u.Update(cleaner)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]bool{"success": true})
	return
}

func (h cleanerHTTPHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	err = h.u.Delete(id)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]bool{"success": true})
	return
}
