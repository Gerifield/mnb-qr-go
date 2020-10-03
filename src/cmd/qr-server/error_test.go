package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenErrorJSONWithNil(t *testing.T) {
	res := genErrorJSON(0, nil)
	assert.Equal(t, `{"code":0,"error":""}`, string(res))
}
