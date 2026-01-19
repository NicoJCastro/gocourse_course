package course

import (
	"errors"
	"fmt"
)

var ErrInvalidRequestType = errors.New("invalid request type")
var ErrInvalidDefaultLimitConfiguration = errors.New("invalid default limit configuration")
var ErrIDRequired = errors.New("id is required")
var ErrAtLeastOneFieldRequired = errors.New("at least one field is required")
var ErrNameRequired = errors.New("name is required")
var ErrStartDateAndEndDateRequired = errors.New("start_date and end_date are required")
var ErrFailedToCreateCourse = errors.New("failed to create course")
var ErrFailedToGetCourse = errors.New("failed to get course")
var ErrFailedToGetAllCourses = errors.New("failed to get all courses")
var ErrFailedToUpdateCourse = errors.New("failed to update course")
var ErrFailedToDeleteCourse = errors.New("failed to delete course")
var ErrFailedToCountCourses = errors.New("failed to count courses")
var ErrInvalidStartDate = errors.New("invalid start date format")
var ErrInvalidEndDate = errors.New("invalid end date format")

// ErrNotFound es un error personalizado que incluye el ID del curso no encontrado
type ErrNotFound struct {
	CourseID string
}

// Error implementa la interfaz error
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("course with ID %s not found", e.CourseID)
}

// Unwrap permite usar errors.Is() con este error
func (e *ErrNotFound) Unwrap() error {
	return ErrNotFoundBase
}

// NewErrNotFound crea una nueva instancia de ErrNotFound
func NewErrNotFound(courseID string) *ErrNotFound {
	return &ErrNotFound{CourseID: courseID}
}

// ErrNotFoundBase es un error sentinela para comparaciones con errors.Is()
var ErrNotFoundBase = errors.New("course not found")
