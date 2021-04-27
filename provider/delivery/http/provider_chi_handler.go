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

type providerHTTPChiHandler struct {
	U domain.ProviderUsecase
}

func NewProviderHTTPChiHandler(u domain.ProviderUsecase) http.Handler {
	handler := providerHTTPChiHandler{U: u}
	router := chi.NewRouter()
	router.Get("/", handler.getAll)
	router.Get("/{id}", handler.getById)
	router.Post("/", handler.insert)
	router.Post("/{id}", handler.update)
	router.Delete("/{id}", handler.delete)
	return router
}

func (h providerHTTPChiHandler) getAll(w http.ResponseWriter, r *http.Request) {
	providers, err := h.U.GetAll()
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, providers)
}

func (h providerHTTPChiHandler) getById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	provider, err := h.U.GetById(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils_http.WriteErrorResponse(w, http.StatusNotFound, err)
			return
		}
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	utils_http.WriteJson(w, provider)
}

func (h providerHTTPChiHandler) insert(w http.ResponseWriter, r *http.Request) {
	var provider domain.Provider
	err := json.NewDecoder(r.Body).Decode(&provider)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.U.Insert(&provider)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]int64{"id": provider.Id})
	return
}

func (h providerHTTPChiHandler) update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	var provider domain.Provider
	err = json.NewDecoder(r.Body).Decode(&provider)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	provider.Id = id
	err = h.U.Update(provider)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(200)
	utils_http.WriteJson(w, map[string]bool{"success": true})
	return
}

func (h providerHTTPChiHandler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	err = h.U.Delete(id)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]bool{"success": true})
	return
}
