package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"master-data-service/internal/opd/web"
)

// SimpegConfig adalah konfigurasi external service Simpeg
type SimpegConfig struct {
	BaseURL string
	ApiPath string
}

// SimpegClient adalah HTTP client untuk berkomunikasi dengan Simpeg
type SimpegClient struct {
	Config SimpegConfig
	HTTP   *http.Client
	Logger *slog.Logger
}

func NewSimpegClient(
	cfg SimpegConfig,
	// logger taruh paling bawah
	logger *slog.Logger,
) *SimpegClient {
	return &SimpegClient{
		Config: cfg,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		Logger: logger,
	}
}

// MasterOpdSimpeg membuat GET request dummy ke Simpeg BASE_URL
// untuk mengambil data master OPD dari upstream Simpeg
func (c *SimpegClient) MasterOpdSimpeg(
	ctx context.Context,
) ([]web.SimpegOpdResponse, error) {
	const op = "simpeg-client-MasterOpdSimpeg"

	url := fmt.Sprintf("%s/%s/master_opd/1", c.Config.BaseURL, c.Config.ApiPath)

	c.Logger.Debug(
		"get master opd from simpeg",
		"op", op,
		"url", url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		c.Logger.Error(
			"failed building simpeg request",
			"op", op,
			"url", url,
			"error", err,
		)
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		c.Logger.Error(
			"simpeg request failed",
			"op", op,
			"url", url,
			"error", err,
		)
		return nil, fmt.Errorf("simpeg request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Logger.Error(
			"simpeg unexpected status code",
			"op", op,
			"url", url,
			"status_code", resp.StatusCode,
		)
		return nil, fmt.Errorf(
			"simpeg unexpected status code %d",
			resp.StatusCode,
		)
	}

	var response web.BaseSimpegResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.Logger.Error(
			"failed decoding simpeg response",
			"op", op,
			"url", url,
			"error", err,
		)

		return nil, fmt.Errorf(
			"%s: decode response: %w",
			op,
			err,
		)
	}

	opds := response.Data

	c.Logger.Debug(
		"successfully fetched master opd from simpeg",
		"op", op,
		"url", url,
		"count", len(opds),
	)

	return opds, nil
}
