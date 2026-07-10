package web

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Service struct {
	repo       JobRepository
	dataFolder string
}

func NewService(repo JobRepository, dataFolder string) *Service {
	return &Service{
		repo:       repo,
		dataFolder: dataFolder,
	}
}

func (s *Service) Create(ctx context.Context, job *Job) error {
	return s.repo.Create(ctx, job)
}

func (s *Service) All(ctx context.Context) ([]Job, error) {
	return s.repo.Select(ctx, SelectParams{})
}

func (s *Service) Get(ctx context.Context, id string) (Job, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return s.repo.Delete(ctx, id)
	}

	csvName := CSVFileName(job.Name)
	datapath := filepath.Join(s.dataFolder, csvName)

	// Only delete the CSV file if no other jobs share this same sanitized name
	allJobs, err := s.repo.Select(ctx, SelectParams{})
	if err == nil {
		hasDuplicates := false
		for _, j := range allJobs {
			if j.ID != id && CSVFileName(j.Name) == csvName {
				hasDuplicates = true
				break
			}
		}
		if !hasDuplicates {
			if _, err := os.Stat(datapath); err == nil {
				_ = os.Remove(datapath)
			}
		}
	} else {
		// Fallback to simple removal if select fails
		if _, err := os.Stat(datapath); err == nil {
			_ = os.Remove(datapath)
		}
	}

	return s.repo.Delete(ctx, id)
}

func (s *Service) Update(ctx context.Context, job *Job) error {
	return s.repo.Update(ctx, job)
}

func (s *Service) SelectPending(ctx context.Context) ([]Job, error) {
	return s.repo.Select(ctx, SelectParams{Status: StatusPending, Limit: 1})
}

func (s *Service) GetCSV(ctx context.Context, id string) (string, error) {
	job, err := s.repo.Get(ctx, id)
	if err != nil {
		return "", err
	}

	csvName := CSVFileName(job.Name)
	datapath := filepath.Join(s.dataFolder, csvName)

	if _, err := os.Stat(datapath); os.IsNotExist(err) {
		return "", fmt.Errorf("csv file not found for job %s (name: %s)", id, job.Name)
	}

	return datapath, nil
}
