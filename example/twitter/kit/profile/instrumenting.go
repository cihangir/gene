package main

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/twitter/models"
	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	ProfileService
}

func (mw instrumentingMiddleware) Create(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "ProfileServiceCreate"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.ProfileService.Create(ctx, req)
	return
}

func (mw instrumentingMiddleware) Delete(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "ProfileServiceDelete"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.ProfileService.Delete(ctx, req)
	return
}

func (mw instrumentingMiddleware) One(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "ProfileServiceOne"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.ProfileService.One(ctx, req)
	return
}

func (mw instrumentingMiddleware) Update(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "ProfileServiceUpdate"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.ProfileService.Update(ctx, req)
	return
}
