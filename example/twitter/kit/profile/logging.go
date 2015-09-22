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
	ProfileService
}

func (mw loggingMiddleware) Create(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "ProfileServiceCreate",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.ProfileService.Create(ctx, req)
	return
}

func (mw loggingMiddleware) Delete(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "ProfileServiceDelete",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.ProfileService.Delete(ctx, req)
	return
}

func (mw loggingMiddleware) One(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "ProfileServiceOne",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.ProfileService.One(ctx, req)
	return
}

func (mw loggingMiddleware) Update(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "ProfileServiceUpdate",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.ProfileService.Update(ctx, req)
	return
}
