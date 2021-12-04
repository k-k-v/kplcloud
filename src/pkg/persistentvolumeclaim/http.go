/**
 * @Time : 8/11/21 4:26 PM
 * @Author : solacowa@gmail.com
 * @File : http
 * @Software: GoLand
 */

package persistentvolumeclaim

import (
	"context"
	"encoding/json"
	valid "github.com/asaskevich/govalidator"
	"github.com/go-kit/kit/endpoint"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/kplcloud/kplcloud/src/encode"
)

func MakeHTTPHandler(s Service, dmw []endpoint.Middleware, opts []kithttp.ServerOption) http.Handler {
	ems := []endpoint.Middleware{}

	ems = append(ems, dmw...)

	eps := NewEndpoint(s, map[string][]endpoint.Middleware{
		"Create": ems,
		"Sync":   ems,
		"List":   ems,
	})

	r := mux.NewRouter()

	r.Handle("/{cluster}/create/{namespace}", kithttp.NewServer(
		eps.CreateEndpoint,
		decodeCreateRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)
	r.Handle("/{cluster}/sync/{namespace}", kithttp.NewServer(
		eps.SyncEndpoint,
		kithttp.NopRequestDecoder,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/list/{namespace}", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/list/{storage}/storage", kithttp.NewServer(
		eps.ListEndpoint,
		decodeListRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/info/{namespace}/get/{name}", kithttp.NewServer(
		eps.GetEndpoint,
		decodeInfoRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodGet)
	r.Handle("/{cluster}/delete/{namespace}/name/{name}", kithttp.NewServer(
		eps.DeleteEndpoint,
		decodeInfoRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodDelete)

	return r
}

func decodeInfoRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	vars := mux.Vars(r)
	name, ok := vars["name"]
	if !ok {
		return nil, encode.InvalidParams.Error()
	}
	req.name = name
	return req, nil
}

func decodeListRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req listRequest
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if page < 0 {
		page = 1
	}
	if pageSize < 0 {
		page = 10
	}
	req.page = page
	req.pageSize = pageSize
	req.storage = r.URL.Query().Get("storage")
	if strings.EqualFold(req.storage, "") {
		vars := mux.Vars(r)
		req.storage, _ = vars["storage"]
	}

	return req, nil
}

func decodeCreateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, encode.InvalidParams.Wrap(err)
	}
	validResult, err := valid.ValidateStruct(req)
	if err != nil {
		return nil, encode.InvalidParams.Wrap(err)
	}
	if !validResult {
		return nil, encode.InvalidParams.Wrap(errors.New("valid false"))
	}
	return req, nil
}
