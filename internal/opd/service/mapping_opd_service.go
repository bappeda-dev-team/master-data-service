package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"master-data-service/internal/opd/domain"
	"master-data-service/internal/opd/repository"
	opdWeb "master-data-service/internal/opd/web"
)

var (
	ErrMappingOpdNotFound = errors.New("mapping opd tidak ditemukan")
	ErrMappingOpdInvalid  = errors.New("mapping opd tidak valid")
	ErrMappingOpdExists   = errors.New("mapping opd sudah ada")
)

type MappingOpdService struct {
	Repo   *repository.MappingOpdRepository
	Logger *slog.Logger
}

func NewMappingOpdService(
	repo *repository.MappingOpdRepository,
	logger *slog.Logger,
) *MappingOpdService {
	return &MappingOpdService{
		Repo:   repo,
		Logger: logger,
	}
}

func (s *MappingOpdService) Create(
	ctx context.Context,
	mapping *domain.MappingOpd,
) (*opdWeb.MappingOpdResponse, error) {
	if err := validateMappingOpd(mapping); err != nil {
		return nil, err
	}

	mapping.Sumber = strings.ToUpper(strings.TrimSpace(mapping.Sumber))
	mapping.KodeSumber = strings.TrimSpace(mapping.KodeSumber)
	mapping.KodeMaster = strings.TrimSpace(mapping.KodeMaster)

	existing, err := s.Repo.FindBySumberAndKodeSumber(
		ctx,
		mapping.Sumber,
		mapping.KodeSumber,
	)
	if err == nil && existing != nil {
		return nil, ErrMappingOpdExists
	}

	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// IDK WHAT IS THIS
	// existing, err = s.Repo.FindBySumberAndKodeMaster(
	// 	ctx,
	// 	mapping.Sumber,
	// 	mapping.KodeMaster,
	// )
	// if err == nil && existing != nil {
	// 	return nil, ErrMappingOpdExists
	// }

	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	if err := s.Repo.Create(ctx, mapping); err != nil {
		return nil, err
	}

	response := toMappingOpdResponse(*mapping)

	return &response, nil
}

func (s *MappingOpdService) FindByID(
	ctx context.Context,
	id int64,
) (*opdWeb.MappingOpdResponse, error) {
	mapping, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrMappingOpdNotFound
		}

		return nil, err
	}

	response := toMappingOpdResponse(*mapping)

	return &response, nil
}

func (s *MappingOpdService) FindBySumberAndKodeSumber(
	ctx context.Context,
	sumber string,
	kodeSumber string,
) (*opdWeb.MappingOpdResponse, error) {
	sumber = strings.ToUpper(strings.TrimSpace(sumber))
	kodeSumber = strings.TrimSpace(kodeSumber)

	mapping, err := s.Repo.FindBySumberAndKodeSumber(
		ctx,
		sumber,
		kodeSumber,
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrMappingOpdNotFound
		}

		return nil, err
	}

	response := toMappingOpdResponse(*mapping)

	return &response, nil
}

func (s *MappingOpdService) FindByKodeMaster(
	ctx context.Context,
	kodeMaster string,
) (*opdWeb.MappingOpdResponse, error) {
	kodeMaster = strings.TrimSpace(kodeMaster)

	mapping, err := s.Repo.FindBySumberAndKodeMaster(
		ctx,
		kodeMaster,
	)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrMappingOpdNotFound
		}

		return nil, err
	}

	response := toMappingOpdResponse(*mapping)

	return &response, nil
}

func (s *MappingOpdService) FindAll(
	ctx context.Context,
) ([]opdWeb.MappingOpdResponse, error) {
	mappings, err := s.Repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]opdWeb.MappingOpdResponse, 0, len(mappings))

	for _, mapping := range mappings {
		responses = append(
			responses,
			toMappingOpdResponse(mapping),
		)
	}

	return responses, nil
}

func (s *MappingOpdService) Update(
	ctx context.Context,
	mapping *domain.MappingOpd,
) (*opdWeb.MappingOpdResponse, error) {
	if mapping == nil || mapping.ID <= 0 {
		return nil, ErrMappingOpdInvalid
	}

	if err := validateMappingOpd(mapping); err != nil {
		return nil, err
	}

	mapping.Sumber = strings.ToUpper(strings.TrimSpace(mapping.Sumber))
	mapping.KodeSumber = strings.TrimSpace(mapping.KodeSumber)
	mapping.KodeMaster = strings.TrimSpace(mapping.KodeMaster)

	existing, err := s.Repo.FindByID(ctx, mapping.ID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrMappingOpdNotFound
		}

		return nil, err
	}

	if existing.Sumber != mapping.Sumber ||
		existing.KodeSumber != mapping.KodeSumber {

		found, err := s.Repo.FindBySumberAndKodeSumber(
			ctx,
			mapping.Sumber,
			mapping.KodeSumber,
		)

		if err == nil && found != nil && found.ID != mapping.ID {
			return nil, ErrMappingOpdExists
		}

		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return nil, err
		}
	}

	// IDK WHAT IS THIS
	// if existing.Sumber != mapping.Sumber ||
	// 	existing.KodeMaster != mapping.KodeMaster {

	// 	found, err := s.Repo.FindBySumberAndKodeMaster(
	// 		ctx,
	// 		mapping.Sumber,
	// 		mapping.KodeMaster,
	// 	)

	// 	if err == nil && found != nil && found.ID != mapping.ID {
	// 		return nil, ErrMappingOpdExists
	// 	}

	// 	if err != nil && !errors.Is(err, repository.ErrNotFound) {
	// 		return nil, err
	// 	}
	// }

	if err := s.Repo.Update(ctx, mapping); err != nil {
		return nil, err
	}

	response := toMappingOpdResponse(*mapping)

	return &response, nil
}

func (s *MappingOpdService) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return ErrMappingOpdInvalid
	}

	_, err := s.Repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrMappingOpdNotFound
		}

		return err
	}

	return s.Repo.Delete(ctx, id)
}

func validateMappingOpd(mapping *domain.MappingOpd) error {
	if mapping == nil {
		return ErrMappingOpdInvalid
	}

	if strings.TrimSpace(mapping.Sumber) == "" {
		return fmt.Errorf("%w: sumber wajib diisi", ErrMappingOpdInvalid)
	}

	if strings.TrimSpace(mapping.KodeSumber) == "" {
		return fmt.Errorf("%w: kode sumber wajib diisi", ErrMappingOpdInvalid)
	}

	if strings.TrimSpace(mapping.KodeMaster) == "" {
		return fmt.Errorf("%w: kode master wajib diisi", ErrMappingOpdInvalid)
	}

	return nil
}

func toMappingOpdResponse(
	mapping domain.MappingOpd,
) opdWeb.MappingOpdResponse {
	return opdWeb.MappingOpdResponse{
		ID:         mapping.ID,
		Sumber:     mapping.Sumber,
		KodeSumber: mapping.KodeSumber,
		KodeMaster: mapping.KodeMaster,
		CreatedAt:  mapping.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  mapping.UpdatedAt.Format(time.RFC3339),
	}
}
