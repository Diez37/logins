package v1

import (
	"encoding/json"
	"errors"
	"github.com/Diez37/logins/domain"
	"github.com/Diez37/logins/infrastructure/repository"
	"github.com/diez37/go-packages/clients/db"
	"github.com/diez37/go-packages/log"
	"github.com/diez37/go-packages/server/http/helpers"
	"github.com/go-http-utils/headers"
	"github.com/google/uuid"
	"github.com/ldez/mimetype"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/errgroup"
	"io"
	"io/ioutil"
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

func (handler *API) Add(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.Add")
	defer span.End()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	login := domain.Login{}

	if err := json.Unmarshal(body, &login); err != nil {
		handler.errorHelper.Error(http.StatusBadRequest, err, writer)
		return
	}

	loginFromRepository, err := handler.repository.Insert(ctx, &repository.Login{
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

	content, err := json.Marshal(&Login{
		Uuid:      loginFromRepository.Uuid,
		Login:     loginFromRepository.Login,
		Ban:       loginFromRepository.Ban,
		CreatedAt: loginFromRepository.CreatedAt,
		UpdateAt:  loginFromRepository.UpdateAt,
	})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) UpdateByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.UpdateByUuid")
	defer span.End()

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	login := domain.Login{}

	if err := json.Unmarshal(body, &login); err != nil {
		handler.errorHelper.Error(http.StatusBadRequest, err, writer)
		return
	}

	loginFromRepository, err := handler.repository.Update(ctx, &repository.Login{
		Uuid:      ctx.Value(UuidFieldName).(uuid.UUID),
		Login:     login.Login,
		Ban:       login.Ban,
		CreatedAt: login.CreatedAt,
		UpdateAt:  login.UpdateAt,
	})
	if err != nil && err != db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	if err == db.RecordNotFoundError {
		handler.errorHelper.Error(http.StatusNotFound, errors.New(http.StatusText(http.StatusNotFound)), writer)
		return
	}

	content, err := json.Marshal(&Login{
		Uuid:      loginFromRepository.Uuid,
		Login:     loginFromRepository.Login,
		Ban:       loginFromRepository.Ban,
		CreatedAt: loginFromRepository.CreatedAt,
		UpdateAt:  loginFromRepository.UpdateAt,
	})
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) FindByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.FindByUuid")
	defer span.End()

	login, err := handler.repository.FindByUuid(ctx, ctx.Value(UuidFieldName).(uuid.UUID))
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

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) BanByUuid(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.BanByUuid")
	defer span.End()

	_, err := handler.repository.BanByUuid(ctx, ctx.Value(UuidFieldName).(uuid.UUID))
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

func (handler *API) FindByLogin(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.FindByLogin")
	defer span.End()

	login, err := handler.repository.FindByLogin(ctx, ctx.Value(LoginFieldName).(string))
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

	writer.Header().Add(headers.ContentType, mimetype.ApplicationJSON)
	writer.WriteHeader(http.StatusOK)

	if _, err := writer.Write(content); err != nil {
		handler.logger.Error(err)
	}
}

func (handler *API) Page(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.Page")
	defer span.End()

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
			Page:  page,
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

func (handler *API) Count(writer http.ResponseWriter, request *http.Request) {
	ctx, span := handler.tracer.Start(request.Context(), "api:v1.Count")
	defer span.End()

	count, err := handler.repository.Count(ctx)
	if err != nil {
		handler.errorHelper.Error(http.StatusInternalServerError, err, writer)
		return
	}

	writer.Header().Set(headers.ContentType, mimetype.TextPlain)
	writer.WriteHeader(http.StatusOK)
	if _, err := writer.Write([]byte(strconv.FormatInt(count, 10))); err != nil {
		handler.logger.Error(err)
	}
}
