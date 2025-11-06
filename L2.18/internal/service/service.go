package service

import (
	"context"
	"time"

	"task-manager/internal/models"
	"task-manager/internal/repo"
)

type Service struct {
	repo repo.EventRepo
}

func NewService(r repo.EventRepo) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateEvent(ctx context.Context, e *models.Event) error {
	return s.repo.CreateEvent(ctx, e)
}

func (s *Service) UpdateEvent(ctx context.Context, e *models.Event) error {
	return s.repo.UpdateEvent(ctx, e)
}

func (s *Service) DeleteEvent(ctx context.Context, id int64) error {
	return s.repo.DeleteEvent(ctx, id)
}

func (s *Service) GetEventsForDay(ctx context.Context, userID int64, date time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForDay(ctx, userID, date)
}

func (s *Service) GetEventsForRange(ctx context.Context, userID int64, start, end time.Time) ([]models.Event, error) {
	return s.repo.GetEventsForRange(ctx, userID, start, end)
}
