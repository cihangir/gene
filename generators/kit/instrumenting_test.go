package kit

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestFunctions(t *testing.T) {
	s := &schema.Schema{}
	err := json.Unmarshal([]byte(testdata.TestDataFull), s)

	s = s.Resolve(s)

	sts, err := GenerateInstrumenting(common.NewContext(), s)
	common.TestEquals(t, nil, err)
	common.TestEquals(t, expecteds[0], string(sts[0].Content))
}

var expecteds = []string{`package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.TimeHistogram
	AccountService
}

func (mw instrumentingMiddleware) Create(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "AccountServiceCreate"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.AccountService.Create(ctx, req)
	return
}

func (mw instrumentingMiddleware) Delete(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "AccountServiceDelete"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.AccountService.Delete(ctx, req)
	return
}

func (mw instrumentingMiddleware) One(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "AccountServiceOne"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.AccountService.One(ctx, req)
	return
}

func (mw instrumentingMiddleware) Some(ctx context.Context, req *models.Account) (res *[]*models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "AccountServiceSome"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.AccountService.Some(ctx, req)
	return
}

func (mw instrumentingMiddleware) Update(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		methodField := metrics.Field{Key: "method", Value: "AccountServiceUpdate"}
		errorField := metrics.Field{Key: "error", Value: fmt.Sprintf("%v", err)}
		mw.requestCount.With(methodField).With(errorField).Add(1)
		mw.requestLatency.With(methodField).With(errorField).Observe(time.Since(begin))
	}(time.Now())

	res, err = mw.AccountService.Update(ctx, req)
	return
}
`}
