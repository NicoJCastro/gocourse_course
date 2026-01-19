package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/NicoJCastro/go_lib_response/response"
	"github.com/NicoJCastro/gocourse_course/internal/course"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoint) http.Handler {
	mux := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	// 🎯 POST /courses - Crear course
	mux.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse,
		encodeResponse,
		opts...,
	)).Methods("POST")

	// 🎯 GETALL /courses - Obtener todos los cursos (con paginación y filtros)
	mux.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourses,
		encodeResponse,
		opts...,
	)).Methods("GET")

	// 🎯 GET /courses/{id} - Obtener un curso por ID
	mux.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse,
		encodeResponse,
		opts...,
	)).Methods("GET")

	// 🎯 PATCH /courses/{id} - Actualizar curso
	mux.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	// 🎯 DELETE /courses/{id} - Eliminar curso
	mux.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse,
		encodeResponse,
		opts...,
	)).Methods("DELETE")
	return mux
}

// 🎯 Decoder para CREATE: decodifica el body JSON
func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest("invalid JSON format")
	}
	return req, nil
}

// 🎯 Decoder para GET: extrae el ID de la URL
func decodeGetCourse(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, response.BadRequest(course.ErrIDRequired.Error())
	}
	return course.GetReq{ID: id}, nil
}

// 🎯 Decoder para GET ALL: extrae query parameters (limit, page, filters)
func decodeGetAllCourses(_ context.Context, r *http.Request) (interface{}, error) {
	// Extraer query parameters
	query := r.URL.Query()

	// Convertir limit y page a int
	limit := 0
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	page := 0
	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	// Construir GetAllReq con los query parameters
	req := course.GetAllReq{
		Name:  query.Get("name"),
		Limit: limit,
		Page:  page,
	}

	return req, nil
}

// 🎯 Decoder para UPDATE: extrae ID de la URL y body JSON
func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, response.BadRequest(course.ErrIDRequired.Error())
	}

	var req course.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest("invalid JSON format")
	}

	// Asignar el ID extraído de la URL
	req.ID = id
	return req, nil
}

// 🎯 Decoder para DELETE: extrae el ID de la URL
func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || id == "" {
		return nil, response.BadRequest(course.ErrIDRequired.Error())
	}
	return course.DeleteReq{ID: id}, nil
}

// 🎯 Encoder para todas las respuestas exitosas
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	respObj, ok := resp.(response.Response)
	if !ok {
		// Si no es response.Response, es un error de programación
		// encodeError debería haberlo manejado
		respObj = response.InternalServerError("invalid response type")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(respObj.StatusCode())
	return json.NewEncoder(w).Encode(respObj)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 🔍 Intentamos convertir el error a response.Response
	resp, ok := err.(response.Response)
	if !ok {
		// ❌ Si no es response.Response, es un error estándar de Go
		// 💡 Lo convertimos a InternalServerError como fallback seguro
		resp = response.InternalServerError(err.Error())
	}

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
