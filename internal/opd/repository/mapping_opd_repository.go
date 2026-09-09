package repository

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"master-data-service/internal/opd/domain"
)

type MappingOpdRepository struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewMappingOpdRepository(
	db *sql.DB,
	logger *slog.Logger,
) *MappingOpdRepository {
	return &MappingOpdRepository{
		DB:     db,
		Logger: logger,
	}
}

func (repo *MappingOpdRepository) Create(
	ctx context.Context,
	mapping *domain.MappingOpd,
) error {
	const query = `
		INSERT INTO mapping_opd (
			sumber,
			kode_sumber,
			kode_master
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			created_at,
			updated_at
	`

	return repo.DB.QueryRowContext(
		ctx,
		query,
		mapping.Sumber,
		mapping.KodeSumber,
		mapping.KodeMaster,
	).Scan(
		&mapping.ID,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)
}

func (repo *MappingOpdRepository) FindByID(
	ctx context.Context,
	id int64,
) (*domain.MappingOpd, error) {
	const query = `
		SELECT
			id,
			sumber,
			kode_sumber,
			kode_master,
			created_at,
			updated_at
		FROM mapping_opd
		WHERE id = $1
	`

	var mapping domain.MappingOpd

	err := repo.DB.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&mapping.ID,
		&mapping.Sumber,
		&mapping.KodeSumber,
		&mapping.KodeMaster,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &mapping, nil
}

func (repo *MappingOpdRepository) FindBySumberAndKodeSumber(
	ctx context.Context,
	sumber string,
	kodeSumber string,
) (*domain.MappingOpd, error) {
	const query = `
		SELECT
			id,
			sumber,
			kode_sumber,
			kode_master,
			created_at,
			updated_at
		FROM mapping_opd
		WHERE sumber = $1
		  AND kode_sumber = $2
	`

	var mapping domain.MappingOpd

	err := repo.DB.QueryRowContext(
		ctx,
		query,
		sumber,
		kodeSumber,
	).Scan(
		&mapping.ID,
		&mapping.Sumber,
		&mapping.KodeSumber,
		&mapping.KodeMaster,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &mapping, nil
}

func (repo *MappingOpdRepository) FindBySumberAndKodeMaster(
	ctx context.Context,
	sumber string,
	kodeMaster string,
) (*domain.MappingOpd, error) {
	const query = `
		SELECT
			id,
			sumber,
			kode_sumber,
			kode_master,
			created_at,
			updated_at
		FROM mapping_opd
		WHERE sumber = $1
		  AND kode_master = $2
	`

	var mapping domain.MappingOpd

	err := repo.DB.QueryRowContext(
		ctx,
		query,
		sumber,
		kodeMaster,
	).Scan(
		&mapping.ID,
		&mapping.Sumber,
		&mapping.KodeSumber,
		&mapping.KodeMaster,
		&mapping.CreatedAt,
		&mapping.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &mapping, nil
}

func (repo *MappingOpdRepository) FindAll(
	ctx context.Context,
) ([]domain.MappingOpd, error) {
	const query = `
		SELECT
			id,
			sumber,
			kode_sumber,
			kode_master,
			created_at,
			updated_at
		FROM mapping_opd
		ORDER BY id
	`

	rows, err := repo.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mappings := make([]domain.MappingOpd, 0)

	for rows.Next() {
		var mapping domain.MappingOpd

		err := rows.Scan(
			&mapping.ID,
			&mapping.Sumber,
			&mapping.KodeSumber,
			&mapping.KodeMaster,
			&mapping.CreatedAt,
			&mapping.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		mappings = append(mappings, mapping)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mappings, nil
}

func (repo *MappingOpdRepository) Update(
	ctx context.Context,
	mapping *domain.MappingOpd,
) error {
	const query = `
		UPDATE mapping_opd
		SET
			sumber = $1,
			kode_sumber = $2,
			kode_master = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING updated_at
	`

	return repo.DB.QueryRowContext(
		ctx,
		query,
		mapping.Sumber,
		mapping.KodeSumber,
		mapping.KodeMaster,
		mapping.ID,
	).Scan(&mapping.UpdatedAt)
}

func (repo *MappingOpdRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	const query = `
		DELETE FROM mapping_opd
		WHERE id = $1
	`

	_, err := repo.DB.ExecContext(ctx, query, id)

	return err
}
