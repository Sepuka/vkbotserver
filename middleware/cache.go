package middleware

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sepuka/vkbotserver/config"
	"github.com/sepuka/vkbotserver/domain"
	"github.com/sepuka/vkbotserver/message"
	"net/http"
)

type (
	httpWriter struct {
		w    http.ResponseWriter
		body []byte
	}
)

func NewHttpWriter(w http.ResponseWriter) *httpWriter {
	return &httpWriter{
		w: w,
	}
}

func (w *httpWriter) WriteHeader(statusCode int) {
	w.w.WriteHeader(statusCode)
}

func (w *httpWriter) Header() http.Header {
	return w.w.Header()
}

func (w *httpWriter) Write(data []byte) (int, error) {
	w.body = data

	return w.w.Write(data)
}

func Cache(cfg config.Config) func(handlerFunc HandlerFunc) HandlerFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(exec message.Executor, req *domain.Request, w http.ResponseWriter) error {
			const (
				cacheKeyTmpl = `%d_%s`
			)

			var (
				writer   *httpWriter
				client   *redis.Client
				ctx      = context.Background()
				cacheKey = fmt.Sprintf(cacheKeyTmpl, req.Object.Message.FromId, req.Object.Message.Text)
				cache    *redis.StringCmd
				value    []byte
				err      error
			)

			if cfg.Cache.Enabled {
				client = redis.NewClient(&redis.Options{})
				if client != nil {
					cache = client.Get(ctx, cacheKey)
					if cache != nil {
						if value, err = cache.Bytes(); err != nil {
							_, err = w.Write(value)

							return err
						}
					}
				}
			}

			writer = NewHttpWriter(w)
			err = next(exec, req, writer)

			if cfg.Cache.Enabled && client != nil {
				_ = client.Set(ctx, cacheKey, writer.body, cfg.Cache.Ttl)
			}

			return err
		}
	}
}
