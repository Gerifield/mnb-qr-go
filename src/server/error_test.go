package server

import (
	"errors"
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
	sendError(res, 500, errors.New("test err"))

	assert.Equal(t, 500, res.Code)
	assert.Equal(t, `{"code":500,"error":"test err"}`, res.Body.String())
}
