package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	// This service
	pb "github.com/adamryman/ambition-rello/rello-service"
)

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.RegisterRelloServer
	_ = io.Copy
)

type HTTPBody struct{}

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(ctx context.Context, endpoints Endpoints, logger log.Logger) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
		httptransport.ServerBefore(dumpBodyToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	m := http.NewServeMux()

	m.Handle("/", httptransport.NewServer(
		ctx,
		endpoints.CheckListWebhookEndpoint,
		HttpDecodeLogger(DecodeHTTPCheckListWebhookZeroRequest, logger),
		EncodeHTTPGenericResponse,
		serverOptions...,
	))
	return m
}

func HttpDecodeLogger(next httptransport.DecodeRequestFunc, logger log.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		logger.Log("method", r.Method, "url", r.URL.String())
		rv, err := next(ctx, r)
		if err != nil {
			logger.Log("method", r.Method, "url", r.URL.String(), "Error", err)
		}
		return rv, err
	}
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	code := http.StatusInternalServerError
	msg := err.Error()

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorWrapper{Error: msg})
}

func errorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// Server Decode

// DecodeHTTPCheckListWebhookZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded checklistwebhook request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPCheckListWebhookZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req pb.ChecklistUpdate
	err := json.NewDecoder(r.Body).Decode(&req)
	// err = io.EOF if r.Body was empty
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "decoding body of http request")
	}

	pathParams, err := PathParams(r.URL.Path, "/")
	_ = pathParams
	if err != nil {
		fmt.Printf("Error while reading path params: %v\n", err)
		return nil, errors.Wrap(err, "couldn't unmarshal path parameters")
	}
	queryParams, err := QueryParams(r.URL.Query())
	_ = queryParams
	if err != nil {
		fmt.Printf("Error while reading query params: %v\n", err)
		return nil, errors.Wrapf(err, "Error while reading query params: %v", r.URL.Query())
	}

	return &req, err
}

// Client Decode

// DecodeHTTPCheckListWebhook is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded Empty response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPCheckListWebhookResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errorDecoder(r)
	}
	var resp pb.Empty
	err := json.NewDecoder(r.Body).Decode(&resp)
	return &resp, err
}

// Client Encode

// EncodeHTTPCheckListWebhookZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a checklistwebhook request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPCheckListWebhookZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	fmt.Printf("Encoding request %v\n", request)
	req := request.(*pb.ChecklistUpdate)
	_ = req

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()

	// Set the body parameters
	var buf bytes.Buffer
	toRet := request.(*pb.ChecklistUpdate)

	toRet.Model = req.Model

	toRet.Action = req.Action

	if err := json.NewEncoder(&buf).Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	fmt.Printf("URL: %v\n", r.URL)
	return nil
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// Helper functions

// PathParams takes a url and a gRPC-annotation style url template, and
// returns a map of the named parameters in the template and their values in
// the given url.
//
// PathParams does not support the entirety of the URL template syntax defined
// in third_party/googleapis/google/api/httprule.proto. Only a small subset of
// the functionality defined there is implemented here.
func PathParams(url string, urlTmpl string) (map[string]string, error) {
	rv := map[string]string{}
	pmp := BuildParamMap(urlTmpl)

	expectedLen := len(strings.Split(strings.TrimRight(urlTmpl, "/"), "/"))
	recievedLen := len(strings.Split(strings.TrimRight(url, "/"), "/"))
	if expectedLen != recievedLen {
		return nil, fmt.Errorf("Expected a path containing %d parts, provided path contains %d parts", expectedLen, recievedLen)
	}

	parts := strings.Split(url, "/")
	for k, v := range pmp {
		rv[k] = parts[v]
	}

	return rv, nil
}

// BuildParamMap takes a string representing a url template and returns a map
// indicating the location of each parameter within that url, where the
// location is the index as if in a slash-separated sequence of path
// components. For example, given the url template:
//
//     "/v1/{a}/{b}"
//
// The returned param map would look like:
//
//     map[string]int {
//         "a": 2,
//         "b": 3,
//     }
func BuildParamMap(urlTmpl string) map[string]int {
	rv := map[string]int{}

	parts := strings.Split(urlTmpl, "/")
	for idx, part := range parts {
		if strings.ContainsAny(part, "{}") {
			param := RemoveBraces(part)
			rv[param] = idx
		}
	}
	return rv
}

// RemoveBraces replace all curly braces in the provided string, opening and
// closing, with empty strings.
func RemoveBraces(val string) string {
	val = strings.Replace(val, "{", "", -1)
	val = strings.Replace(val, "}", "", -1)
	return val
}

// QueryParams takes query parameters in the form of url.Values, and returns a
// bare map of the string representation of each key to the string
// representation for each value. The representations of repeated query
// parameters is undefined.
func QueryParams(vals url.Values) (map[string]string, error) {

	rv := map[string]string{}
	for k, v := range vals {
		rv[k] = v[0]
	}
	return rv, nil
}

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	for k, _ := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	return ctx
}

func dumpBodyToContext(ctx context.Context, r *http.Request) context.Context {
	var save []byte
	var err error
	save, r.Body, err = drainBody(r.Body)
	if err != nil {
		return ctx
	}
	ctx = context.WithValue(ctx, HTTPBody{}, save)
	return ctx
}

func drainBody(b io.ReadCloser) ([]byte, io.ReadCloser, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err := b.Close(); err != nil {
		return nil, b, err
	}

	return buf.Bytes(), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
