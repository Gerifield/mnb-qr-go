package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusMethodNotAllowed, resp.Code)
}

func TestWrongInput(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("something"))
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestInvalidSize(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, string(genErrorJSON(http.StatusBadRequest, errInvalidSize)), resp.Body.String())
}

func TestInvalidKind(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"pngSize":5}`))
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, string(genErrorJSON(http.StatusBadRequest, errInvalidKind)), resp.Body.String())
}

func TestInvalidBIC(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"pngSize":5,"kind":"HCT"}`))
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, string(genErrorJSON(http.StatusBadRequest, errors.New("invalid BIC length"))), resp.Body.String())
}

func TestInvalidExpiration(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"pngSize":5,"kind":"HCT","bic":"abcdefgh","name":"Test User","iban":"HU00123456789012345678901234"}`))
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Equal(t, string(genErrorJSON(http.StatusBadRequest, errors.New("negative validity period"))), resp.Body.String())
}

func TestMinimalGenSuccess(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"pngSize":5,"kind":"HCT","bic":"abcdefgh","name":"Test User","iban":"HU00123456789012345678901234","expire":20}`))
	resp := httptest.NewRecorder()
	New().GenerateHandler(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "image/png", resp.Header().Get("Content-Type"))
	assert.True(t, resp.Body.Len() > 100)
}
