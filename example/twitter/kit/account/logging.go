package main

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"github.com/cihangir/gene/example/twitter/models"
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
