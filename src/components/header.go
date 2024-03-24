package components

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"server/src/utils"
	"strings"
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
	StatusOK  Status = "200"
	StatusBAD Status = "500"
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

func NewHeader(content string, response bool) *Header {
	lines := strings.Split(content, utils.END_OF_HEADER_LINE)
	var protocol, version, path, method, protocol_version, status, statusMesage string
	if response {
		utils.Unpack(strings.Fields(lines[0]), &protocol_version, &status, &statusMesage, &statusMesage)
	} else {
		utils.Unpack(strings.Fields(lines[0]), &method, &path, &protocol_version)
	}
	utils.Unpack(strings.Split(protocol_version, "/"), &protocol, &version)

	header := &Header{
		Method:          method,
		Path:            path,
		Protocol:        protocol,
		ProtocolVersion: protocol_version,
	}

	header_map := make(map[string]string)
	var key, val string
	const key_ex = utils.HostKey + ":"

	for _, line := range lines[1:] {
		if strings.Contains(line, key_ex) {
			utils.Unpack(strings.Split(line, key_ex), &key, &val)
			continue
		}
		utils.Unpack(strings.Split(line, ": "), &key, &val)
		header_map[key] = val
	}

	header.HeaderContent = header_map
	return header
}

func NewDefaultHeader(path string, method Method) *Header {
	host := os.Getenv("host")
	return &Header{
		Path:            path,
		ProtocolVersion: DefaultHttp,
		Method:          method,
		HeaderContent: map[string]string{
			utils.HostKey:           host,
			utils.UserAgenteKey:     utils.UserAgent,
			utils.AcceptKey:         utils.AcceptAll,
			utils.ContentTypeKey:    utils.ApplicationJson,
			utils.AcceptEncodingKey: utils.AcceptEncoding,
		},
	}
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
		header.ProtocolVersion, StatusOK, utils.END_OF_HEADER_LINE)

	for key, val := range header.HeaderContent {
		header_str += fmt.Sprintf("%s: %s %s", key, val, utils.END_OF_HEADER_LINE)
	}

	return header_str + utils.END_OF_HEADER_LINE
}
