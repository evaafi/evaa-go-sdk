package price

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var _ Provider = (*SingleEndpointProvider)(nil)

type SingleEndpointProvider struct {
	client *http.Client
	list   map[string]string
	mtx    sync.RWMutex
}

func newSingleEndpointProvider(client *http.Client) *SingleEndpointProvider {
	if client == nil {
		client = http.DefaultClient
	}
	return &SingleEndpointProvider{client: client, list: make(map[string]string)}
}

func (p *SingleEndpointProvider) update(ctx context.Context, url string) (map[string]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	outputsResp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer outputsResp.Body.Close()

	var list = map[string]string{}
	if err := json.NewDecoder(outputsResp.Body).Decode(&list); err != nil {
		return nil, fmt.Errorf("failed to decode response: %s %w", outputsResp.Status, err)
	}

	return list, nil
}

func (p *SingleEndpointProvider) Update(ctx context.Context, url string, interval time.Duration) error {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			res, err := p.update(ctx, url)
			if errors.Is(err, context.Canceled) {
				return nil
			}
			if err != nil {
				return err
			}

			p.mtx.Lock()
			p.list = res
			p.mtx.Unlock()
		}
	}
}

func (p *SingleEndpointProvider) GetRawData(ctx context.Context, _, address string) (*RawData, error) {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	data, ok := p.list[address]
	if !ok {
		return nil, fmt.Errorf("data not inisialized yet")
	}

	rawData, err := Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data, err: %w", err)
	}
	return rawData, nil
}
