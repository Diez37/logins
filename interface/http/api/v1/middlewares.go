package v1

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cast"
	"net/http"
)

func UuidField(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid := chi.URLParam(r, UuidFieldName)

		s, err := cast.ToStringE(uuid)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if s == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), UuidFieldName, s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoginField(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		login := chi.URLParam(r, LoginFieldName)

		s, err := cast.ToStringE(login)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if s == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), LoginFieldName, s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PageFields(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// page
		page := r.URL.Query().Get(PageFieldName)
		if page == "" {
			page = "1"
		}

		i, err := cast.ToIntE(page)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if i < 1 {
			i = 1
		}
		ctx := context.WithValue(r.Context(), PageFieldName, i)

		// limit
		limit := r.URL.Query().Get(LimitFieldName)
		if limit == "" {
			limit = "0"
		}

		u, err := cast.ToUintE(limit)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if u == 0 {
			u = LimitDefault
		}
		ctx = context.WithValue(ctx, LimitFieldName, u)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
