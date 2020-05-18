package qr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	oneHourSeconds = 60 * 60
)

func TestDateToString(t *testing.T) {
	testTable := []struct {
		input          Date
		expectedOutput string
	}{
		{Date(time.Date(2020, 05, 18, 10, 11, 23, 0, time.UTC)), "20200518101123+0"},
		{Date(time.Date(2020, 05, 18, 10, 11, 23, 0, time.FixedZone("testZone1", 2*oneHourSeconds))), "20200518101123+2"},
		{Date(time.Date(2020, 05, 18, 10, 11, 23, 0, time.FixedZone("testZone2", -oneHourSeconds))), "20200518101123-1"},
		{Date(time.Date(2020, 05, 18, 10, 11, 23, 0, time.FixedZone("testZone3", oneHourSeconds*11))), "20200518101123+11"},
	}

	for _, tt := range testTable {
		assert.Equal(t, tt.expectedOutput, tt.input.String())
	}
}
