package price

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/evaafi/evaa-go-sdk/config"
	"github.com/xssnick/tonutils-go/tvm/cell"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Endpoint        = "https://api.stardust-mainnet.iotaledger.net"
	outputsEndpoint = "/api/indexer/v1/outputs/nft/"
	coreEndpoint    = "/api/core/v2/outputs/"
)

type provider struct {
	client *http.Client
}

func newProvider(client *http.Client) Provider {
	if client == nil {
		client = http.DefaultClient
	}
	return &provider{client: client}
}

func (p *provider) GetRawData(ctx context.Context, baseURL, address string) (*RawData, error) {
	if baseURL == "" {
		baseURL = Endpoint
	}
	output, err := p.getOutputID(ctx, baseURL, address)
	if err != nil {
		return nil, fmt.Errorf("failed to get outputID, err: %w", err)
	}
	feature, err := p.getFeature(ctx, baseURL, output)
	if err != nil {
		return nil, fmt.Errorf("failed to get feature, err: %w", err)
	}
	rawData, err := Parse(feature)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data, err: %w", err)
	}
	return rawData, nil
}

func (p *provider) getOutputID(ctx context.Context, baseURL string, address string) (outputID string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+outputsEndpoint+address, nil)
	if err != nil {
		return "", fmt.Errorf("error creating output request: %w", err)
	}
	outputsResp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting outputs: %w", err)
	}
	defer outputsResp.Body.Close()

	var outputIDResp struct {
		Items []string `json:"items"`
	}
	if err := json.NewDecoder(outputsResp.Body).Decode(&outputIDResp); err != nil {
		return "", fmt.Errorf("error decoding output response: %s %w", outputsResp.Status, err)
	}
	if len(outputIDResp.Items) == 0 {
		return "", fmt.Errorf("no items found for NFT ID %s", address)
	}
	return outputIDResp.Items[0], nil
}

func (p *provider) getFeature(ctx context.Context, baseURL, outputID string) (feature string, err error) {
	req, err := http.NewRequestWithContext(ctx, "GET", baseURL+coreEndpoint+outputID, nil)
	if err != nil {
		return "", fmt.Errorf("error creating core request: %w", err)
	}
	coreResp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error getting core data: %w", err)
	}
	defer coreResp.Body.Close()

	var resData struct {
		Output struct {
			Features []struct {
				Data string `json:"data"`
			} `json:"features"`
		} `json:"output"`
	}
	if err := json.NewDecoder(coreResp.Body).Decode(&resData); err != nil {
		return "", fmt.Errorf("error decoding core response: %s %w", coreResp.Status, err)
	}
	if len(resData.Output.Features) == 0 {
		return "", fmt.Errorf("no features found for OutputID %s", outputID)
	}

	return resData.Output.Features[0].Data, nil
}

func Parse(data string) (*RawData, error) {
	if len(data)%2 == 1 {
		return nil, errors.New("invalid price data")
	}
	var dataBuilder strings.Builder
	for i := 2; i < len(data); i += 2 {
		dataBuilder.WriteString("%" + data[i:i+2])
	}
	unescapedData, err := url.QueryUnescape(dataBuilder.String())
	if err != nil {
		return nil, fmt.Errorf("failed to unescape data, err: %w", err)
	}

	var jsonData struct {
		PackedPrices string `json:"packedPrices"`
		Signature    string `json:"signature"`
		PublicKey    string `json:"publicKey"`
		Timestamp    int64  `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(unescapedData), &jsonData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal unescapedData, err: %w", err)
	}

	pricesCell, err := hex.DecodeString(jsonData.PackedPrices)
	if err != nil {
		return nil, fmt.Errorf("failed to decode prices cell from hex, err: %w", err)
	}

	signature, err := hex.DecodeString(jsonData.Signature)
	if err != nil {
		return nil, fmt.Errorf("failed to decode signature from hex, err: %w", err)
	}

	pubKey, err := hex.DecodeString(jsonData.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode pubkey from hex, err: %w", err)
	}

	packedPrices, err := cell.FromBOC(pricesCell)
	if err != nil {
		return nil, fmt.Errorf("failed to decode prices from BOC, err: %w", err)
	}

	priceCell, err := packedPrices.BeginParse().LoadRefCell()
	if err != nil {
		return nil, fmt.Errorf("failed to load refCell from packedPrices, err: %w", err)
	}

	return &RawData{
		PricesDict: priceCell.AsDict(256),
		Signature:  signature,
		PubKey:     pubKey,
		Timestamp:  jsonData.Timestamp,
	}, nil
}

type RawData struct {
	PricesDict *cell.Dictionary
	Signature  []byte
	PubKey     []byte
	Timestamp  int64

	prices map[string]*big.Int
}

func (d *RawData) Prices() map[string]*big.Int {
	if d.prices != nil {
		return d.prices
	}

	d.prices = make(map[string]*big.Int)

	kvs, err := d.PricesDict.LoadAll()
	if err != nil {
		return nil
	}

	for _, kv := range kvs {
		k, err := kv.Key.LoadBigUInt(256)
		if err != nil {
			return nil
		}
		v, err := kv.Value.LoadVarUInt(16)
		if err != nil {
			return nil
		}

		d.prices[k.String()] = v
	}

	return d.prices
}

const ttlOracleData = 120 * time.Second

func (d *RawData) verify(assets map[string]*config.AssetConfig) bool {
	if time.Since(time.Unix(d.Timestamp, 0)) > ttlOracleData {
		return false
	}

	prices := d.Prices()
	if len(prices) < len(assets) {
		return false
	}

	for k, _ := range assets {
		price, ok := prices[k]
		if !ok || price == nil || price.Sign() != 1 {
			return false
		}
	}

	return true
}
