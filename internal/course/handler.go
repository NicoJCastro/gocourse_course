package course

import (
	"context"
	"errors"
	"strconv"

	"github.com/NicoJCastro/go_lib_response/response"
	"github.com/NicoJCastro/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoint struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		Name      string `json:"name"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	GetAllReq struct {
		Name  string `json:"name"`
		Limit int    `json:"limit"`
		Page  int    `json:"page"`
	}

	GetReq struct {
		ID string `json:"id"`
	}

	DeleteReq struct {
		ID string `json:"id"`
	}

	UpdateReq struct {
		ID        string  `json:"id"`
		Name      *string `json:"name"`
		StartDate *string `json:"start_date"`
		EndDate   *string `json:"end_date"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Error  string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

const (
	ErrMsgInvalidRequestType = "invalid request type"
)

func MakeEndpoint(s Service, config Config) Endpoint {
	return Endpoint{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateReq)
		if !ok {
			return nil, response.BadRequest(ErrMsgInvalidRequestType)
		}
		if req.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}
		if req.StartDate == "" || req.EndDate == "" {
			return nil, response.BadRequest(ErrStartDateAndEndDateRequired.Error())
		}
		course, err := s.Create(req.Name, req.StartDate, req.EndDate)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}
		return response.Created("Course created successfully", course, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetReq)
		if !ok {
			return nil, response.BadRequest(ErrMsgInvalidRequestType)
		}
		if req.ID == "" {
			return nil, response.BadRequest(ErrIDRequired.Error())
		}
		course, err := s.Get(req.ID)
		if err != nil {
			var notFoundErr *ErrNotFound
			if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("Course retrieved successfully", course, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(GetAllReq)
		if !ok {
			return nil, response.BadRequest(ErrMsgInvalidRequestType)
		}

		filters := Filters{
			Name: req.Name,
		}

		// Extraemos limit y page directamente del struct GetAllReq
		// Si los valores son 0 (no proporcionados), usaremos valores por defecto
		limit := req.Limit
		page := req.Page

		// 🔧 Validación: si limit es 0, usamos el valor por defecto de la configuración
		if limit <= 0 {
			defaultLimit, err := strconv.Atoi(config.LimPageDef)
			if err != nil {
				return nil, response.InternalServerError("invalid default limit configuration")
			}
			limit = defaultLimit
		}

		// 🔧 Validación: si page es 0 o negativo, establecemos página 1
		if page <= 0 {
			page = 1
		}

		count, err := s.Count(filters)
		if err != nil {
			return nil, response.InternalServerError("error counting courses: " + err.Error())
		}

		metaData, err := meta.New(page, limit, int(count), config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError("error generating metadata: " + err.Error())
		}

		courses, err := s.GetAll(filters, metaData.Offset(), metaData.Limit())
		if err != nil {
			return nil, response.InternalServerError("error retrieving courses: " + err.Error())
		}

		return response.OK("Courses retrieved successfully", courses, metaData), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		reqUpdate, ok := request.(UpdateReq)
		if !ok {
			return nil, response.BadRequest(ErrMsgInvalidRequestType)
		}
		if reqUpdate.ID == "" {
			return nil, response.BadRequest(ErrIDRequired.Error())
		}
		if reqUpdate.Name == nil && reqUpdate.StartDate == nil && reqUpdate.EndDate == nil {
			return nil, response.BadRequest(ErrAtLeastOneFieldRequired.Error())
		}

		// 🔧 Validación: si se proporciona un campo, no puede estar vacío
		if reqUpdate.Name != nil && *reqUpdate.Name == "" {
			return nil, response.BadRequest(ErrNameRequired.Error())
		}

		// 🔧 Validación: si se proporciona StartDate, debe tener un formato válido y no estar vacío
		if reqUpdate.StartDate != nil && *reqUpdate.StartDate == "" {
			return nil, response.BadRequest(ErrStartDateAndEndDateRequired.Error())
		}

		// 🔧 Validación: si se proporciona EndDate, debe tener un formato válido y no estar vacío
		if reqUpdate.EndDate != nil && *reqUpdate.EndDate == "" {
			return nil, response.BadRequest(ErrStartDateAndEndDateRequired.Error())
		}

		err := s.Update(reqUpdate.ID, reqUpdate.Name, reqUpdate.StartDate, reqUpdate.EndDate)
		if err != nil {
			var notFoundErr *ErrNotFound
			if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("Course updated successfully", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(DeleteReq)
		if !ok {
			return nil, response.BadRequest(ErrMsgInvalidRequestType)
		}
		if req.ID == "" {
			return nil, response.BadRequest(ErrIDRequired.Error())
		}
		err := s.Delete(req.ID)
		if err != nil {
			var notFoundErr *ErrNotFound
			if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}
		return response.OK("Course deleted successfully", nil, nil), nil
	}
}
