package server

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenErrorJSONWithNil(t *testing.T) {
	res := genErrorJSON(0, nil)
	assert.Equal(t, `{"code":0,"error":""}`, string(res))
}

func TestSendError(t *testing.T) {
	res := httptest.NewRecorder()
	sendError(res, 1, nil)

	assert.Equal(t, 1, res.Code)
	assert.Equal(t, `{"code":1,"error":""}`, res.Body.String())
}
