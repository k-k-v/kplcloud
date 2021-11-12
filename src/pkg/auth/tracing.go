/**
 * @Time : 8/11/21 4:22 PM
 * @Author : solacowa@gmail.com
 * @File : tracing
 * @Software: GoLand
 */

package auth

import (
	"context"

	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// 链路追踪中间件
type tracing struct {
	next   Service
	tracer stdopentracing.Tracer
}

func (s *tracing) Register(ctx context.Context, username, email, password, mobile, remark string) (err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Register", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Auth",
	})
	defer func() {
		span.LogKV(
			"username", username,
			"mobile", mobile,
			"remark", remark,
			"email", email,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.Register(ctx, username, email, password, mobile, remark)
}

func (s *tracing) Login(ctx context.Context, username, password string) (rs string, sessionTimeout int64, err error) {
	span, ctx := stdopentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "Login", stdopentracing.Tag{
		Key:   string(ext.Component),
		Value: "package.Auth",
	})
	defer func() {
		span.LogKV(
			"username", username,
			"err", err,
		)
		span.Finish()
	}()
	return s.next.Login(ctx, username, password)
}

func NewTracing(otTracer stdopentracing.Tracer) Middleware {
	return func(next Service) Service {
		return &tracing{
			next:   next,
			tracer: otTracer,
		}
	}
}