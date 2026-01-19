package course

import (
	"log"
	"time"

	"github.com/NicoJCastro/gocourse_domain/domain"
)

type (
	Filters struct {
		Name string
	}

	Service interface {
		Create(name, startDate, endDate string) (*domain.Course, error)
		Get(id string) (*domain.Course, error)
		GetAll(filters Filters, offset, limit int) ([]domain.Course, error)
		Delete(id string) error
		Update(id string, name *string, startDate *string, endDate *string) error
		Count(filters Filters) (int64, error)
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

func (s *service) Create(name, startDate, endDate string) (*domain.Course, error) {
	s.log.Println("---- Creating course ----")

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.log.Println("Error parsing start date:", err)
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.log.Println("Error parsing end date:", err)
		return nil, err
	}

	course := &domain.Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}

	if err := s.repo.Create(course); err != nil {
		s.log.Printf("Error creating course: %v\n", err)
		return nil, err
	}

	return course, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]domain.Course, error) {
	s.log.Println("---- Getting all courses ----")
	courses, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		s.log.Printf("Error getting courses: %v\n", err)
		return nil, err
	}
	return courses, nil
}

func (s service) Get(id string) (*domain.Course, error) {
	course, err := s.repo.Get(id)
	if err != nil {
		s.log.Printf("Error getting course: %v\n", err)
		return nil, err
	}
	return course, nil
}

func (s service) Delete(id string) error {
	s.log.Println("---- Deleting course ----")
	return s.repo.Delete(id)
}

func (s service) Update(id string, name *string, startDate *string, endDate *string) error {
	s.log.Println("---- Updating course ----")

	var startDateParsed, endDateParsed *time.Time

	if startDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			s.log.Printf("Error parsing start date: %v\n", err)
			return err
		}
		startDateParsed = &parsedDate
	}

	if endDate != nil {
		parsedDate, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			s.log.Printf("Error parsing end date: %v\n", err)
			return err
		}
		endDateParsed = &parsedDate
	}

	return s.repo.Update(id, name, startDateParsed, endDateParsed)
}

func (s service) Count(filters Filters) (int64, error) {
	return s.repo.Count(filters)
}
