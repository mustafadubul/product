package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mustafadubul/product/internal/domain"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	httpHandler "github.com/mustafadubul/product/internal/handler/http"
	"github.com/mustafadubul/product/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

type Handler struct {
	*httpHandler.Handler

	cancel   context.CancelFunc
	mockCtrl *gomock.Controller

	service *mocks.MockHTTPService
}

func NewTestHandler(t *testing.T) *Handler {
	mockCtrl := gomock.NewController(t)
	mockService := mocks.NewMockHTTPService(mockCtrl)

	_, cancel := context.WithCancel(context.Background())
	l := zerolog.Nop()
	h := &Handler{
		Handler:  httpHandler.NewHandler(&l, mockService),
		service:  mockService,
		mockCtrl: mockCtrl,
	}

	h.cancel = cancel
	return h
}

func (h *Handler) Finish() {
	h.mockCtrl.Finish()
	h.cancel()
}

type testRequest struct {
	method    string
	endpoint  string
	handler   http.HandlerFunc
	payload   interface{}
	urlParams map[string]string
}

func TestHandler_Search(t *testing.T) {
	h := NewTestHandler(t)

	products := []domain.Product{
		{ItemName: "camera"},
	}

	query := &domain.Query{
		Lat:    15,
		Lng:    10,
		Radius: 5,
		Term:   "camera",
	}

	h.service.EXPECT().Search(gomock.Any(), query).Return(products, nil)

	request := testRequest{
		method:    http.MethodGet,
		endpoint:  "/q?radius=5&lng=10&lat=15&term=camera",
		handler:   h.Search,
		payload:   query,
		urlParams: nil,
	}

	res := httpTestRequestRecord(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	expectBody, _ := json.Marshal(products)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, expectBody, data)
}

func TestHandler_Create(t *testing.T) {
	h := NewTestHandler(t)
	defer h.Finish()

	product := &domain.Product{
		Lat:      15,
		Lng:      10,
		ItemName: "camera",
	}

	updatedProduct := &domain.Product{
		ID:       99,
		Lat:      15,
		Lng:      10,
		ItemName: "camera",
	}

	h.service.EXPECT().Create(gomock.Any(), product).Return(updatedProduct, nil)

	request := testRequest{
		method:    http.MethodPost,
		endpoint:  httpHandler.CreateEndpoint,
		handler:   h.Create,
		payload:   product,
		urlParams: nil,
	}

	res := httpTestRequestRecord(request)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	expected := &domain.Product{
		ID:       99,
		Lat:      15,
		Lng:      10,
		ItemName: "camera",
	}
	expectBody, _ := json.Marshal(expected)
	data, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, expectBody, data)
}

func TestHandler_Get(t *testing.T) {
	h := NewTestHandler(t)
	defer h.Finish()

	product := &domain.Product{
		ID:       99,
		Lat:      15,
		Lng:      10,
		ItemName: "camera",
	}
	h.service.EXPECT().Get(gomock.Any(), uint64(99)).Return(product, nil)

	request := testRequest{
		method:    http.MethodGet,
		endpoint:  httpHandler.GetEndpoint,
		handler:   h.Get,
		payload:   nil,
		urlParams: map[string]string{"id": "99"},
	}

	res := httpTestRequestRecord(request)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHandler_Delete(t *testing.T) {
	h := NewTestHandler(t)
	defer h.Finish()

	request := testRequest{
		method:    http.MethodDelete,
		endpoint:  httpHandler.DeleteEndpoint,
		handler:   h.Delete,
		payload:   nil,
		urlParams: map[string]string{"id": "99"},
	}

	h.service.EXPECT().Delete(gomock.Any(), uint64(99)).Return(nil)
	res := httpTestRequestRecord(request)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestHandler_Update(t *testing.T) {
	h := NewTestHandler(t)
	defer h.Finish()

	product := &domain.Product{
		ID:       99,
		Lat:      15,
		Lng:      10,
		ItemName: "camera",
	}

	h.service.EXPECT().Update(gomock.Any(), product).Return(product, nil)

	request := testRequest{
		method:    http.MethodPut,
		endpoint:  httpHandler.UpdateEndpoint,
		handler:   h.Update,
		payload:   product,
		urlParams: map[string]string{"id": "99"},
	}

	res := httpTestRequestRecord(request)
	assert.Equal(t, http.StatusAccepted, res.StatusCode)
}

func httpTestRequestRecord(t testRequest) *http.Response {
	var buf *bytes.Buffer
	data, _ := json.Marshal(&t.payload)
	buf = bytes.NewBuffer(data)

	// replace url parameters
	rctx := chi.NewRouteContext()
	var e string
	for k, v := range t.urlParams {
		e = strings.Replace(t.endpoint, "{"+k+"}", v, -1)
		rctx.URLParams.Add(k, v)
	}
	if e == "" {
		e = t.endpoint
	}

	req := httptest.NewRequest(t.method, e, buf)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rec := httptest.NewRecorder()

	t.handler(rec, req)

	return rec.Result()
}
