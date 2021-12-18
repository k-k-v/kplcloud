/**
 * @Time : 2021/12/6 10:38 AM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package terminal

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/kplcloud/kplcloud/src/encode"
	"github.com/kplcloud/kplcloud/src/middleware"
	"time"
)

type (
	tokenRequest struct {
		PodName     string `json:"podName"`
		Container   string `json:"container"`
		ServiceName string `json:"serviceName"`
	}
	tokenResult struct {
		Namespace   string    `json:"namespace"`
		PodName     string    `json:"podName"`
		Container   string    `json:"container"`
		ErrMsg      string    `json:"errMsg"`
		SessionId   string    `json:"sessionId"`
		Token       string    `json:"token"`
		Cluster     string    `json:"cluster"`
		Containers  []string  `json:"containers"`
		Phase       string    `json:"phase"`
		HostIp      string    `json:"hostIp"`
		PodIp       string    `json:"podIp"`
		StartTime   time.Time `json:"startTime"`
		ServiceName string    `json:"serviceName"`
		Pods        []string  `json:"pods"`
	}
)

type Endpoints struct {
	TokenEndpoint endpoint.Endpoint
}

func NewEndpoint(s Service, dmw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		TokenEndpoint: makeTokenEndpoint(s),
	}

	for _, m := range dmw["Token"] {
		eps.TokenEndpoint = m(eps.TokenEndpoint)
	}
	return eps
}

func makeTokenEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		clusterId, ok := ctx.Value(middleware.ContextKeyClusterId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		userId, ok := ctx.Value(middleware.ContextUserId).(int64)
		if !ok {
			return nil, encode.ErrClusterNotfound.Error()
		}
		ns, _ := ctx.Value(middleware.ContextKeyNamespaceName).(string)
		if !ok {
			return nil, encode.ErrNamespaceNotfound.Error()
		}
		req := request.(tokenRequest)
		res, err := s.Token(ctx, userId, clusterId, ns, req.ServiceName, req.PodName)
		return encode.Response{
			Data:  res,
			Error: err,
		}, err
	}
}
