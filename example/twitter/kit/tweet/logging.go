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
	TweetService
}

func (mw loggingMiddleware) Create(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "TweetServiceCreate",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.TweetService.Create(ctx, req)
	return
}

func (mw loggingMiddleware) Delete(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "TweetServiceDelete",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.TweetService.Delete(ctx, req)
	return
}

func (mw loggingMiddleware) One(ctx context.Context, req *models.Account) (res *models.Account, err error) {
	defer func(begin time.Time) {
		input, _ := json.Marshal(req)
		output, _ := json.Marshal(res)
		_ = mw.logger.Log(
			"method", "TweetServiceOne",
			"input", string(input),
			"output", string(output),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	res, err = mw.TweetService.One(ctx, req)
	return
}
