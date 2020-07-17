package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"errors"
	"net/url"

	"github.com/go-chi/chi"
	"github.com/mustafadubul/product/internal/domain"
	"github.com/mustafadubul/product/internal/repository"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger  *zerolog.Logger
	service Service
}

// mockgen -source=http.go  -package=mocks -destination=../../../mocks/mocks_http_service.go -mock_names Service=MockHTTPService
type Service interface {
	Create(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Get(ctx context.Context, id uint64) (*domain.Product, error)
	Search(ctx context.Context, q *domain.Query) ([]domain.Product, error)

	Update(ctx context.Context, p *domain.Product) (*domain.Product, error)
	Delete(ctx context.Context, id uint64) error
}

func NewHandler(l *zerolog.Logger, svc Service) *Handler {
	componentLogger := l.With().Str("component", "http-handler").Logger()
	return &Handler{
		logger:  &componentLogger,
		service: svc,
	}
}

var (
	CreateEndpoint = "/product"

	GetEndpoint    = "/product/{id}"
	DeleteEndpoint = "/product/{id}"
	UpdateEndpoint = "/product/{id}"
	SearchEndpoint = "/q"
)

func (h *Handler) Setup() http.Handler {
	r := chi.NewRouter()
	r.Get(GetEndpoint, h.Get)

	r.Post(CreateEndpoint, h.Create)

	r.Get(SearchEndpoint, h.Search)

	r.Delete(DeleteEndpoint, h.Delete)
	r.Put(UpdateEndpoint, h.Update)

	return r
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := h.logger.With().Str("handler", "Search").Logger()
	l.WithContext(ctx)

	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		l.Info().Interface("payload", q).Msg("invalid url parameters")
		_ = writeError(w, http.StatusBadRequest, err)
		return
	}

	query, err := validateSearchInput(q)
	if err != nil {
		l.Info().Interface("query", query).Msg("invalid request query")
		_ = writeError(w, http.StatusBadRequest, err)
		return
	}

	results, err := h.service.Search(ctx, query)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			l.Error().Err(err).Interface("payload", q).Msg("failed to making db request")
			_ = writeError(w, http.StatusInternalServerError, err)
			return
		}
	}
	_ = writeJSON(w, http.StatusOK, results)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := h.logger.With().Str("handler", "Create").Logger()
	l.WithContext(ctx)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Info().Err(err).Msg("failed to read body")
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var product domain.Product
	if err = json.Unmarshal(data, &product); err != nil {
		l.Info().Err(err).Msg("failed to unmarshal body")
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.service.Create(ctx, &product)
	if err != nil {
		l.Info().Err(err).Msg("failed to create product")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	_ = writeJSON(w, http.StatusCreated, p)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := h.logger.With().Str("handler", "Delete").Logger()
	l.WithContext(ctx)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		l.Info().Interface("id", chi.URLParam(r, "id")).Msg("id not valid")
		_ = writeError(w, http.StatusBadRequest, err)
		return
	}
	err = h.service.Delete(ctx, id)
	if err != nil {
		l.Error().Err(err).Interface("id", chi.URLParam(r, "id")).Msg("failed to delete product")
		_ = writeError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := h.logger.With().Str("handler", "Update").Logger()
	l.WithContext(ctx)

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		l.Info().Err(err).Msg("failed to read body")
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var product domain.Product
	if err = json.Unmarshal(data, &product); err != nil {
		l.Info().Err(err).Msg("failed to unmarshal body")
		writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.service.Update(ctx, &product)
	if err != nil {
		l.Info().Err(err).Msg("failed to update product")
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	_ = writeJSON(w, http.StatusAccepted, p)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	l := h.logger.With().Str("handler", "Search").Logger()
	l.WithContext(ctx)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		l.Info().Interface("id", chi.URLParam(r, "id")).Msg("id not valid")
		_ = writeError(w, http.StatusBadRequest, err)
		return
	}

	p, err := h.service.Get(ctx, id)
	if err != nil {
		l.Error().Err(err).Interface("id", chi.URLParam(r, "id")).Msg("failed to delete product")
		_ = writeError(w, http.StatusInternalServerError, err)
		return
	}
	_ = writeJSON(w, http.StatusOK, p)
}

func validateSearchInput(v url.Values) (*domain.Query, error) {
	if v.Get("lat") == "" {
		return nil, fmt.Errorf("missing lat")
	}
	lat, err := strconv.ParseFloat(v.Get("lat"), 64)
	if err != nil {
		return nil, fmt.Errorf("lat invalid value")
	}

	if v.Get("lng") == "" {
		return nil, fmt.Errorf("missing lng")
	}
	lng, err := strconv.ParseFloat(v.Get("lng"), 64)
	if err != nil {
		return nil, fmt.Errorf("lng invalid value")
	}
	if v.Get("radius") == "" {
		return nil, fmt.Errorf("missing radius")
	}
	radius, err := strconv.ParseFloat(v.Get("radius"), 64)
	if err != nil {
		return nil, fmt.Errorf("lng invalid value")
	}

	return &domain.Query{
		Lat:    lat,
		Lng:    lng,
		Radius: radius,
		Term:   v.Get("term")}, nil
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) error {
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.WriteHeader(status)
	_, err = w.Write(body)
	return err
}

func writeError(w http.ResponseWriter, status int, err error) error {
	return writeJSON(w, status, ErrorStatus{Error: err.Error()})
}

type ErrorStatus struct {
	Error string `json:"error"`
}
