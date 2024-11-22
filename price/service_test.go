package price

import (
	"context"
	"testing"
	"time"

	"github.com/evaafi/evaa-go-sdk/config"
)

func TestService_GetPrices(t *testing.T) {
	cfg := config.GetMainMainnetConfig()
	service := NewService(cfg, newProvider(nil))
	prices, err := service.GetPrices(context.Background(), Endpoint, "https://evaa.space")
	if err != nil {
		t.Fatalf("failed to get prices, err: %s", err)
	}
	if len(prices.list) < len(cfg.Assets) {
		t.Errorf("prices count want >%d, got %d", len(cfg.Assets), len(prices.list))
	}
	for k, v := range prices.list {
		t.Logf("%5s: %10s", cfg.Assets[k].Name, v.String())
	}
	if prices.minTimestamp == 0 {
		t.Errorf("minTimestamp is zero")
	}
	if prices.data == nil {
		t.Errorf("data is empty")
	}
}

func TestService_GetPrices_multiEndpoint(t *testing.T) {
	cfg := config.GetMainMainnetConfig()
	service := NewService(cfg, newProvider(nil))
	prices, err := service.GetPrices(context.Background(), "https://evaa.space", "http://localhost", Endpoint)
	if err != nil {
		t.Fatalf("failed to get prices, err: %s", err)
	}
	if len(prices.list) < len(cfg.Assets) {
		t.Errorf("prices count want >%d, got %d", len(cfg.Assets), len(prices.list))
	}
	for k, v := range prices.list {
		t.Logf("%5s: %10s", cfg.Assets[k].Name, v.String())
	}
	if prices.minTimestamp == 0 {
		t.Errorf("minTimestamp is zero")
	}
	if prices.data == nil {
		t.Errorf("data is empty")
	}
}

func TestService_GetPrices_singleEndpoint(t *testing.T) {
	cfg := config.GetMainMainnetConfig()
	endpointProvider := newSingleEndpointProvider(nil)
	go endpointProvider.Update(context.Background(), "https://evaa.space/api/prices", time.Second)
	service := NewService(cfg, endpointProvider)
	time.Sleep(3 * time.Second)
	prices, err := service.GetPrices(context.Background())
	if err != nil {
		t.Fatalf("failed to get prices, err: %s", err)
	}
	if len(prices.list) < len(cfg.Assets) {
		t.Errorf("prices count want >%d, got %d", len(cfg.Assets), len(prices.list))
	}
	for k, v := range prices.list {
		t.Logf("%5s: %10s", cfg.Assets[k].Name, v.String())
	}
	if prices.minTimestamp == 0 {
		t.Errorf("minTimestamp is zero")
	}
	if prices.data == nil {
		t.Errorf("data is empty")
	}
}
