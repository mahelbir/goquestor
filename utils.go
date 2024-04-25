package goquestor

import (
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"strings"
)

func EncodeBody(body url.Values) io.Reader {
	return strings.NewReader(body.Encode())
}

func JSONBody(body any) io.Reader {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(jsonData)
}
