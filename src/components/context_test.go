package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContext(tst *testing.T) {
	context := &Context{
		Header: &Header{
			Method:          "GET",
			Path:            "",
			Protocol:        "HTTP/1.1",
			ProtocolVersion: "HTTP/1.1",
			HeaderContent:   map[string]string{},
		},
		Body: " something",
	}

	str_input := "GET HTTP/1.1 \r\n\r\n something"

	assert.Equal(tst, context, NewContext(str_input))
}
