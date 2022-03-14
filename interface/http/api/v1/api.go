package v1

import (
	"encoding/json"
	"github.com/Diez37/logins/infrastructure/repository"
	"github.com/diez37/go-packages/clients/db"
	"github.com/diez37/go-packages/log"
	"github.com/go-http-utils/headers"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/ldez/mimetype"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type API struct {
	repository repository.Repository
	tracer     trace.Tracer
	logger     log.Logger
	validator  *validator.Validate
}

func NewAPI(repository repository.Repository, tracer trace.Tracer, logger log.Logger, validator *validator.Validate) *API {
	return &API{repository: repository, tracer: tracer, logger: logger, validator: validator}
}

func (handler *API) Add(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "Add")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	login := Login{}
	if err := json.Unmarshal(body, &login); err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		handler.logger.Error(err)
		return
	}

	if err := handler.validator.Struct(login); err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		handler.logger.Error(err)
		return
	}

	loginForRepository := &repository.Login{
		Uuid:      login.Uuid,
		Login:     login.Login,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	}
	if login.Banned != nil {
		loginForRepository.Banned = *login.Banned
	}

	loginForRepository, err = handler.repository.Insert(ctx, loginForRepository)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	handler.logger.Infof("api:v1:add: login '%s', uuid '%s'", loginForRepository.Login, loginForRepository.Uuid.String())

	content, err := json.Marshal(&Login{
		Uuid:      loginForRepository.Uuid,
		Login:     loginForRepository.Login,
		Banned:    &loginForRepository.Banned,
		CreatedAt: loginForRepository.CreatedAt,
		UpdateAt:  loginForRepository.UpdateAt,
	})
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
	}
}

func (handler *API) UpdateByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "UpdateByUuid")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	login := Login{}
	if err := json.Unmarshal(body, &login); err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		handler.logger.Error(err)
		return
	}

	loginFromRepository, err := handler.repository.FindByUuid(ctx, ctx.Value(UuidFieldName).(uuid.UUID))
	if err != nil && err != db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	if err == db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	loginFromRepository.Login = login.Login
	if login.Banned != nil {
		loginFromRepository.Banned = *login.Banned
	}

	loginFromRepository, err = handler.repository.Update(ctx, loginFromRepository)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	handler.logger.Infof("api:v1:update: login '%s'", loginFromRepository.Uuid.String())

	content, err := json.Marshal(&Login{
		Uuid:      loginFromRepository.Uuid,
		Login:     loginFromRepository.Login,
		Banned:    &loginFromRepository.Banned,
		CreatedAt: loginFromRepository.CreatedAt,
		UpdateAt:  loginFromRepository.UpdateAt,
	})
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
	}
}

func (handler *API) FindByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "FindByUuid")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	login, err := handler.repository.FindByUuid(ctx, ctx.Value(UuidFieldName).(uuid.UUID))
	if err != nil && err != db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	if err == db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	content, err := json.Marshal(&Login{
		Uuid:      login.Uuid,
		Login:     login.Login,
		Banned:    &login.Banned,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	})
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
	}
}

func (handler *API) BanByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "BanByUuid")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	_, err := handler.repository.BanByUuid(ctx, ctx.Value(UuidFieldName).(uuid.UUID))
	if err != nil && err != db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	if err == db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	handler.logger.Infof("api:v1:ban: login '%s'", ctx.Value(UuidFieldName).(uuid.UUID).String())

	writer.WriteHeader(http.StatusOK)
}

func (handler *API) FindByLogin(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "FindByLogin")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	login, err := handler.repository.FindByLogin(ctx, ctx.Value(LoginFieldName).(string))
	if err != nil && err != db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	if err == db.RecordNotFoundError {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	content, err := json.Marshal(&Login{
		Uuid:      login.Uuid,
		Login:     login.Login,
		Banned:    &login.Banned,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	})
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
	}
}

func (handler *API) Page(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "Page")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	page := ctx.Value(PageFieldName).(uint)
	limit := ctx.Value(LimitFieldName).(uint)

	var totalCount int64
	var models []*repository.Login

	wg := &errgroup.Group{}

	wg.Go(func() error {
		count, err := handler.repository.Count(ctx)
		totalCount = count

		return err
	})

	wg.Go(func() error {
		logins, err := handler.repository.Page(ctx, page-1, limit)
		models = logins

		return err
	})

	if err := wg.Wait(); err != nil && err != io.EOF {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	logins := make([]*Login, len(models))
	for index, login := range models {
		logins[index] = &Login{
			Uuid:      login.Uuid,
			Login:     login.Login,
			Banned:    &login.Banned,
			CreatedAt: login.CreatedAt,
			UpdateAt:  login.UpdateAt,
		}
	}

	content, err := json.Marshal(&Page{
		Meta: &Meta{
			Count: totalCount,
			Page:  page,
			Limit: limit,
		},
		Records: logins,
	})
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	writer.Header().Set(headers.ContentType, mimetype.ApplicationJSON)
	writer.Header().Set(CountHeaderName, strconv.FormatInt(totalCount, 10))
	writer.Header().Set(PageHeaderName, strconv.FormatUint(uint64(page), 10))
	writer.Header().Set(LimitHeaderName, strconv.FormatUint(uint64(limit), 10))
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
	}
}

func (handler *API) Count(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "Count")
	defer span.End()

	span.SetAttributes(
		attribute.String("interface", "http"),
		attribute.String("handler", "api.v1"),
	)

	count, err := handler.repository.Count(ctx)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
		return
	}

	writer.Header().Set(headers.ContentType, mimetype.TextPlain)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write([]byte(strconv.FormatInt(count, 10))); err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		handler.logger.Error(err)
	}
}
