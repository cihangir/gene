package kit

import (
	"encoding/json"

	"testing"

	"github.com/cihangir/gene/generators/common"
	"github.com/cihangir/gene/testdata"
	"github.com/cihangir/schema"
)

func TestLogging(t *testing.T) {
	s := &schema.Schema{}
	err := json.Unmarshal([]byte(testdata.TestDataFull), s)

	s = s.Resolve(s)

	sts, err := GenerateLogging(common.NewContext(), s)
	common.TestEquals(t, nil, err)
	common.TestEquals(t, loggingExpecteds[0], string(sts[0].Content))
}

var loggingExpecteds = []string{`package main

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	AccountService
}

func (mw loggingMiddleware) Create(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "AccountServiceCreate",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.AccountService.Create(ctx, req)
	return
}

func (mw loggingMiddleware) Delete(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "AccountServiceDelete",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.AccountService.Delete(ctx, req)
	return
}

func (mw loggingMiddleware) One(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "AccountServiceOne",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.AccountService.One(ctx, req)
	return
}

func (mw loggingMiddleware) Some(ctx context.Context, req *models.Account) (res *[]*models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "AccountServiceSome",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.AccountService.Some(ctx, req)
	return
}

func (mw loggingMiddleware) Update(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "AccountServiceUpdate",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.AccountService.Update(ctx, req)
	return
}
`}
