package octrace

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestDdTraceIDtoOtel(t *testing.T) {
	a := assert.New(t)
	traceID := "4014821548718594605"
	otelTraceID, err := ddTraceIDtoOtel(context.Background(), traceID)
	a.NoError(err)
	a.Equal("37b782ed57f49a2d0000000000000000", otelTraceID.String())
}

func TestDdSpanIDtoOtel(t *testing.T) {
	a := assert.New(t)
	parentSpanID := "1424722419751834193"
	otelSpanID := ddSpanIDtoOtel(parentSpanID)

	spew.Dump(otelSpanID)
	a.Equal("13c5a17e04d7d251", otelSpanID.String())
}
