package components

import (
	"server/src/utils"
	"strings"
)

type Context struct {
	Url           string
	Header        *Header
	Endpoint      string
	Body          string
	DiscoveryInfo string
	Timeout       int
}

func NewContext(content string) *Context {
	var header_content, body, rest string

	utils.Unpack(strings.Split(content, utils.END_HEADER_PATTERN), &header_content, &body, &rest)
	header := NewHeader(header_content, true)

	return &Context{
		Header: header,
		Body:   body + rest,
	}
}
