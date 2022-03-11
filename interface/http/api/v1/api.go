package v1

import (
	"encoding/json"
	"errors"
	"github.com/Diez37/logins/infrastructure/repository"
	"github.com/diez37/go-packages/clients/db"
	"github.com/diez37/go-packages/log"
	"github.com/diez37/go-packages/server/http/helpers"
	"github.com/go-http-utils/headers"
	"github.com/google/uuid"
	"github.com/ldez/mimetype"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strconv"
)

type API struct {
	repository  repository.Repository
	tracer      trace.Tracer
	errorHelper *helpers.Error
	logger      log.Logger
}

func NewAPI(repository repository.Repository, tracer trace.Tracer, errorHelper *helpers.Error, logger log.Logger) *API {
	return &API{repository: repository, tracer: tracer, errorHelper: errorHelper, logger: logger}
}

func (handler *API) FindByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.FindByUuid")
	defer span.End()

	value := ctx.Value(UuidFieldName)

	s, err := cast.ToStringE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	uuid, err := uuid.Parse(s)
	if err != nil {
		handler.errorHelper.Error(http.StatusBadRequest, err, writer)
		return
	}

	login, err := handler.repository.FindByUuid(ctx, uuid)
	if err != nil && err != db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	if err == db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)), writer)
		return
	}

	content, err := json.Marshal(&Login{
		Uuid:      login.Uuid,
		Login:     login.Login,
		Ban:       login.Ban,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) FindByLogin(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.FindByLogin")
	defer span.End()

	value := ctx.Value(LoginFieldName)

	s, err := cast.ToStringE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	login, err := handler.repository.FindByLogin(ctx, s)
	if err != nil && err != db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	if err == db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)), writer)
		return
	}

	content, err := json.Marshal(&Login{
		Uuid:      login.Uuid,
		Login:     login.Login,
		Ban:       login.Ban,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) Add(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.Add")
	defer span.End()

	value := ctx.Value(LoginFieldName)

	s, err := cast.ToStringE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	login, err := handler.repository.Save(ctx, &repository.Login{Login: s})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	content, err := json.Marshal(&Login{
		Uuid:      login.Uuid,
		Login:     login.Login,
		Ban:       login.Ban,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) BanByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.BanByUuid")
	defer span.End()

	value := ctx.Value(UuidFieldName)

	s, err := cast.ToStringE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	uuid, err := uuid.Parse(s)
	if err != nil {
		handler.errorHelper.Error(http.StatusBadRequest, err, writer)
		return
	}

	_, err = handler.repository.BanByUuid(ctx, uuid)
	if err != nil && err != db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	if err == db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)), writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *API) BanByLogin(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.BanByUuid")
	defer span.End()

	value := ctx.Value(LoginFieldName)

	login, err := cast.ToStringE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	_, err = handler.repository.BanByLogin(ctx, login)
	if err != nil && err != db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	if err == db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)), writer)
		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (handler *API) Page(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.Page")
	defer span.End()

	var page, limit uint

	value := ctx.Value(PageFieldName)
	u, err := cast.ToUintE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusBadRequest, err, writer)
		return
	}
	page = u - 1

	value = ctx.Value(LimitFieldName)
	u, err = cast.ToUintE(value)
	if err != nil {
		handler.errorHelper.Error(http.StatusBadRequest, err, writer)
		return
	}
	limit = u

	var totalCount int64
	var models []*repository.Login

	wg := &errgroup.Group{}

	wg.Go(func() error {
		count, err := handler.repository.Count(ctx)
		totalCount = count

		return err
	})

	wg.Go(func() error {
		logins, err := handler.repository.Page(ctx, page, limit)
		models = logins

		return err
	})

	if err := wg.Wait(); err != nil && err != io.EOF {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	logins := make([]*Login, len(models))
	for index, login := range models {
		logins[index] = &Login{
			Uuid:      login.Uuid,
			Login:     login.Login,
			Ban:       login.Ban,
			CreatedAt: login.CreatedAt,
			UpdateAt:  login.UpdateAt,
		}
	}

	content, err := json.Marshal(&Page{
		Meta: &Meta{
			Count: totalCount,
			Page:  page + 1,
			Limit: limit,
		},
		Records: logins,
	})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.Header().Set(headers.ContentType, mimetype.ApplicationJSON)
	writer.Header().Set(CountHeaderName, strconv.FormatInt(totalCount, 10))
	writer.Header().Set(PageHeaderName, strconv.FormatUint(uint64(page), 10))
	writer.Header().Set(LimitHeaderName, strconv.FormatUint(uint64(limit), 10))
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}
