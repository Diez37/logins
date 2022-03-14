package v1

import (
	"context"
	"github.com/diez37/go-packages/log"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"net/http"
)

type Uuid struct {
	logger log.Logger
}

func NewUuid(logger log.Logger) *Uuid {
	return &Uuid{logger: logger}
}

func (middleware *Uuid) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		uuid, err := uuid.Parse(chi.URLParam(request, UuidFieldName))

		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			middleware.logger.Error(err)
			return
		}

		ctx := context.WithValue(request.Context(), UuidFieldName, uuid)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func LoginField(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		login := chi.URLParam(request, LoginFieldName)

		if login == "" {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), LoginFieldName, login)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func PageField(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		generalPageNumber := "1"

		if page := request.Header.Get(PageHeaderName); page != "" {
			generalPageNumber = page
		}

		if page := request.URL.Query().Get(PageFieldName); page != "" {
			generalPageNumber = page
		}

		u, err := cast.ToUintE(generalPageNumber)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if u < 1 {
			u = PageDefault
		}
		ctx := context.WithValue(request.Context(), PageFieldName, u)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func LimitField(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		generalLimitNumber := "0"

		if limit := request.Header.Get(LimitHeaderName); limit != "" {
			generalLimitNumber = limit
		}

		if limit := request.URL.Query().Get(LimitFieldName); limit != "" {
			generalLimitNumber = limit
		}

		u, err := cast.ToUintE(generalLimitNumber)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if u == 0 {
			u = LimitDefault
		}
		ctx := context.WithValue(request.Context(), LimitFieldName, u)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
