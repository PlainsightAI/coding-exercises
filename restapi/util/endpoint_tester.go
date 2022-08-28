package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
)

type EndpointTester struct {
	ts     *httptest.Server
	client http.Client
}

func NewEndpointTester(r *chi.Mux) *EndpointTester {
	return &EndpointTester{
		ts:     httptest.NewServer(r),
		client: http.Client{},
	}
}

type FormFile struct {
	Name string
	File []byte
}

type restError struct {
	StatusCode int
	Message    string
}

func (r *restError) Error() string {
	if r.StatusCode == 0 {
		return r.Message
	}
	return fmt.Sprintf("%v: %s", r.StatusCode, r.Message)
}

func (t *EndpointTester) SendAsFormData(method string, path string, files map[string]FormFile, fields map[string]string, response interface{}) *restError {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	for key, r := range files {
		fw, err := w.CreateFormFile(key, r.Name)
		if err != nil {
			return &restError{Message: err.Error()}
		}
		if _, err := fw.Write(r.File); err != nil {
			return &restError{Message: err.Error()}
		}
	}

	for key, r := range fields {
		fw, err := w.CreateFormField(key)
		if err != nil {
			return &restError{Message: err.Error()}
		}
		if _, err := fw.Write([]byte(r)); err != nil {
			return &restError{Message: err.Error()}
		}
	}

	if err := w.Close(); err != nil {
		return &restError{Message: err.Error()}
	}

	req, err := http.NewRequest(method, t.ts.URL+path, &b)
	if err != nil {
		return &restError{Message: err.Error()}
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return t.sendRest(req, response)
}

func (t *EndpointTester) Send(method string, path string, body map[string]interface{}, response interface{}) *restError {
	var bodyReader io.Reader

	if body != nil {
		encodedBody, err := json.Marshal(body)
		if err != nil {
			return &restError{Message: err.Error()}
		}

		bodyReader = bytes.NewReader(encodedBody)
	}

	req, err := http.NewRequest(method, t.ts.URL+path, bodyReader)
	if err != nil {
		return &restError{Message: err.Error()}
	}

	req.Header.Set("content-type", "application/json")

	return t.sendRest(req, response)
}

func (t *EndpointTester) sendRest(req *http.Request, response interface{}) *restError {
	resp, err := t.client.Do(req)
	if err != nil {
		return &restError{Message: err.Error()}
	}

	if resp.ContentLength == 0 || response == nil {
		if resp.StatusCode >= 400 {
			return &restError{StatusCode: resp.StatusCode, Message: resp.Status}
		}
		return nil
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &restError{Message: err.Error()}
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Message string `json:"message"`
		}
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return &restError{StatusCode: resp.StatusCode, Message: err.Error()}
		}

		return &restError{StatusCode: resp.StatusCode, Message: errResp.Message}
	}

	if err := json.Unmarshal(respBody, response); err != nil {
		return &restError{Message: err.Error()}
	}

	return nil
}
