package course

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/NicoJCastro/gocourse_domain/domain"
)

type (
	Filters struct {
		Name string
	}

	Service interface {
		Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error)
		Get(ctx context.Context, id string) (*domain.Course, error)
		GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error)
		Delete(ctx context.Context, id string) error
		Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error
		Count(ctx context.Context, filters Filters) (int64, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, name, startDate, endDate string) (*domain.Course, error) {
	s.log.Println("---- Creating course ----")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println("Error parsing start date:", err)
		return nil, fmt.Errorf("%w: %v", ErrInvalidStartDate, err)
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println("Error parsing end date:", err)
		return nil, fmt.Errorf("%w: %v", ErrInvalidEndDate, err)
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(ctx, course); err != nil {
		s.log.Printf("Error creating course: %v\n", err)
		return nil, fmt.Errorf("%w: %v", ErrFailedToCreateCourse, err)
	}

	return course, nil
}

func (s service) GetAll(ctx context.Context, filters Filters, offset, limit int) ([]domain.Course, error) {
	s.log.Println("---- Getting all courses ----")
	courses, err := s.repo.GetAll(ctx, filters, offset, limit)
	if err != nil {
		s.log.Printf("Error getting courses: %v\n", err)
		// No envolvemos ErrNotFound, lo propagamos directamente
		var notFoundErr *ErrNotFound
		if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrFailedToGetAllCourses, err)
	}
	return courses, nil
}

func (s service) Get(ctx context.Context, id string) (*domain.Course, error) {
	course, err := s.repo.Get(ctx, id)
	if err != nil {
		s.log.Printf("Error getting course: %v\n", err)
		// No envolvemos ErrNotFound, lo propagamos directamente
		var notFoundErr *ErrNotFound
		if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
			return nil, err
		}
		return nil, fmt.Errorf("%w: %v", ErrFailedToGetCourse, err)
	}
	return course, nil
}

func (s service) Delete(ctx context.Context, id string) error {
	s.log.Println("---- Deleting course ----")
	err := s.repo.Delete(ctx, id)
	if err != nil {
		// No envolvemos ErrNotFound, lo propagamos directamente
		var notFoundErr *ErrNotFound
		if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
			return err
		}
		return fmt.Errorf("%w: %v", ErrFailedToDeleteCourse, err)
	}
	return nil
}

func (s service) Update(ctx context.Context, id string, name *string, startDate *string, endDate *string) error {
	s.log.Println("---- Updating course ----")

	var startDateParsed, endDateParsed *time.Time

	if startDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Printf("Error parsing start date: %v\n", err)
			return fmt.Errorf("%w: %v", ErrInvalidStartDate, err)
		}
		startDateParsed = &parsedDate
	}

	if endDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Printf("Error parsing end date: %v\n", err)
			return fmt.Errorf("%w: %v", ErrInvalidEndDate, err)
		}
		endDateParsed = &parsedDate
	}

	err := s.repo.Update(ctx, id, name, startDateParsed, endDateParsed)
	if err != nil {
		// No envolvemos ErrNotFound, lo propagamos directamente
		var notFoundErr *ErrNotFound
		if errors.As(err, &notFoundErr) || errors.Is(err, ErrNotFoundBase) {
			return err
		}
		return fmt.Errorf("%w: %v", ErrFailedToUpdateCourse, err)
	}
	return nil
}

func (s service) Count(ctx context.Context, filters Filters) (int64, error) {
	count, err := s.repo.Count(ctx, filters)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrFailedToCountCourses, err)
	}
	return count, nil
}
