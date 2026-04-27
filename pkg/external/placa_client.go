// Package external  get placa en api externa
package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/oldfarmer96/vehicle-control-go/internal/models"
)

type PlacaClient struct {
	baseURL string
	token   string
	client  *http.Client
}

func NewPlacaClient(baseURL, token string) *PlacaClient {
	return &PlacaClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *PlacaClient) GetPlacaData(ctx context.Context, placa string) (*models.ExternalPlacaResponse, error) {
	payload := map[string]string{"placa": placa}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error al marshalear payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("error al crear request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al hacer request: %w", err)
	}
	defer resp.Body.Close()

	var result models.ExternalPlacaResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta: %w", err)
	}

	return &result, nil
}

