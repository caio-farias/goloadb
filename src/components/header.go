package components

import (
	"encoding/json"
	"fmt"
	"log"
	"server/src/utils"
)

type Method = string

const (
	GET_METHOD    Method = "GET"
	POST_METHOD   Method = "POST"
	PUT_METHOD    Method = "PUT"
	PATCH_METHOD  Method = "PATCH"
	DELETE_METHOD Method = "DELETE"
)

type Status = string

const (
	Status_OK  Status = "200"
	Status_BAD Status = "500"
)

type HeaderContent = map[string]string

type Header struct {
	Path            string
	Method          Method
	Protocol        string
	ProtocolVersion string
	Status          Status
	HeaderContent   HeaderContent
}

func (header *Header) StringAsJson() string {
	to_string, err := json.MarshalIndent(header, "", "	")
	if err != nil {
		log.Println(err)
	}
	return string(to_string)
}

func (header *Header) String() string {
	if header.Method != "" {
		header_str := fmt.Sprintf("%s %s %s%s",
			header.Method, header.Path, header.ProtocolVersion, utils.END_OF_HEADER_LINE)

		for key, val := range header.HeaderContent {
			header_str += fmt.Sprintf("%s: %s%s", key, val, utils.END_OF_HEADER_LINE)
		}

		return header_str + utils.END_OF_HEADER_LINE
	}

	header_str := fmt.Sprintf("%s %s %s",
		header.ProtocolVersion, Status_OK, utils.END_OF_HEADER_LINE)

	for key, val := range header.HeaderContent {
		header_str += fmt.Sprintf("%s: %s %s", key, val, utils.END_OF_HEADER_LINE)
	}

	return header_str + utils.END_OF_HEADER_LINE
}
