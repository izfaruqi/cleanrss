package http

import (
	"cleanrss/domain"
	utils_http "cleanrss/utils/http"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/url"
	"strconv"
)

type entryHTTPChiHandler struct {
	u domain.EntryUsecase
}

func NewEntryHTTPChiHandler(u domain.EntryUsecase) http.Handler {
	handler := entryHTTPChiHandler{u: u}
	router := chi.NewRouter()
	router.Get("/refresh", handler.refreshFromAllProviders)
	router.Get("/refresh/provider/:id", handler.refreshFromProvider)
	router.Get("/query", handler.getByQuery)
	return router
}

func (h entryHTTPChiHandler) refreshFromAllProviders(w http.ResponseWriter, r *http.Request) {
	err := h.u.TriggerRefreshAll()
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]bool{"success": true})
	return
}

func (h entryHTTPChiHandler) refreshFromProvider(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	err = h.u.TriggerRefresh(id)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, map[string]bool{"success": true})
	return
}

func (h entryHTTPChiHandler) getByQuery(w http.ResponseWriter, r *http.Request) {
	var entries []domain.Entry
	queryParams, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}
	query := queryParams.Get("q")
	dateFrom, err := strconv.ParseInt(utils_http.GetParam(queryParams, "date_from", "-1"), 10, 64)
	dateUntil, err := strconv.ParseInt(utils_http.GetParam(queryParams, "date_until", "-1"), 10, 64)
	providerId, err := strconv.ParseInt(utils_http.GetParam(queryParams, "provider_id", "-1"), 10, 64)
	limit, err := strconv.ParseInt(utils_http.GetParam(queryParams, "limit", "40"), 10, 64)
	offset, err := strconv.ParseInt(utils_http.GetParam(queryParams, "offset", "0"), 10, 64)
	includeAll, err := strconv.ParseBool(utils_http.GetParam(queryParams, "include_all", "false"))
	allowRefresh, err := strconv.ParseBool(utils_http.GetParam(queryParams, "allow_refresh", "true"))

	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusBadRequest, errors.New("id is invalid"))
		return
	}
	entries, err = h.u.GetByQuery(query, dateFrom, dateUntil, providerId, limit, offset, includeAll)
	if (entries == nil || len(entries) == 0) && providerId != -1 && allowRefresh && offset == 0 {
		err = h.u.TriggerRefresh(providerId)
		if err != nil {
			utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		entries, err = h.u.GetByQuery(query, dateFrom, dateUntil, providerId, limit, offset, includeAll)
	}
	if err != nil {
		utils_http.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils_http.WriteJson(w, entries)
	return
}
