package price

import (
	"context"
	"testing"
	"time"

	"github.com/evaafi/evaa-go-sdk/config"
)

func TestSingleEndpointProvider_GetRawData(t *testing.T) {
	service := newSingleEndpointProvider(nil)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-time.After(3 * time.Second)
		cancel()
	}()

	err := service.Update(ctx, "https://evaa.space/api/prices", time.Second)
	if err != nil {
		t.Fatalf("failed to GetRawData, err: %s", err)
	}

	rawData, err := service.GetRawData(context.Background(), "", "0xd3a8c0b9fd44fd25a49289c631e3ac45689281f2f8cf0744400b4c65bed38e5d")
	if err != nil {
		t.Fatalf("failed to GetRawData, err: %s", err)
	}
	if !rawData.verify(config.GetMainMainnetConfig().Assets) {
		t.Errorf("verify want true, got false")
	}
	for k, v := range rawData.Prices() {
		t.Logf("%s: %10s", k, v.String())
	}
}
