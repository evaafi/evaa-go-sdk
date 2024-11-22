package price

import (
	"context"
	"fmt"
	"math/big"
	"sort"

	"github.com/xssnick/tonutils-go/tvm/cell"
	"golang.org/x/sync/errgroup"

	"github.com/evaafi/evaa-go-sdk/config"
)

type Provider interface {
	GetRawData(ctx context.Context, baseURL, address string) (*RawData, error)
}

type Prices struct {
	list         map[string]*big.Int
	data         *cell.Cell
	minTimestamp int64
}

func (p *Prices) Get(asset string) *big.Int {
	return p.list[asset]
}

func (p *Prices) Data() *cell.Cell {
	return p.data
}

func (p *Prices) MinTimestamp() int64 {
	return p.minTimestamp
}

type Service struct {
	config        *config.Config
	provider      Provider
	proofSkeleton *cell.ProofSkeleton
}

func NewService(config *config.Config, provider Provider) *Service {
	if provider == nil {
		provider = newProvider(nil)
	}
	proofSkeleton := cell.CreateProofSkeleton()
	proofSkeleton.SetRecursive()
	return &Service{config: config, provider: provider, proofSkeleton: proofSkeleton}
}

type Data struct {
	*RawData
	oracleID uint64
}

func (s *Service) GetPrices(ctx context.Context, endpoint ...string) (*Prices, error) {
	if len(endpoint) == 0 {
		endpoint = append(endpoint, Endpoint)
	}
	ch := make(chan *Data)
	var wg errgroup.Group
	for _, oracle := range s.config.Oracles {
		oracle := oracle
		wg.Go(func() error {
			if len(endpoint) == 1 {
				data, err := s.provider.GetRawData(ctx, endpoint[0], oracle.Address)
				if err != nil {
					return fmt.Errorf("failed to get oracle %d data form %s, err: %w", oracle.ID, endpoint[0], err)
				}
				ch <- &Data{
					RawData:  data,
					oracleID: oracle.ID,
				}
				return nil
			}

			var wg1 errgroup.Group
			oracleCh := make(chan *RawData)
			for _, baseURL := range endpoint {
				baseURL := baseURL
				wg1.Go(func() error {
					data, err := s.provider.GetRawData(ctx, baseURL, oracle.Address)
					if err != nil {
						return fmt.Errorf("failed to get oracle %d data from %s, err: %w", oracle.ID, baseURL, err)
					}
					oracleCh <- data
					return nil
				})
			}

			var endpointErr error
			go func() {
				endpointErr = wg1.Wait()
				close(oracleCh)
			}()

			data, ok := <-oracleCh
			if !ok {
				return endpointErr
			}
			ch <- &Data{
				RawData:  data,
				oracleID: oracle.ID,
			}
			return nil
		})
	}

	var oracleErr error
	go func() {
		oracleErr = wg.Wait()
		close(ch)
	}()

	acceptedPrices := make([]*Data, 0, len(s.config.Oracles))
	for data := range ch {
		if !data.verify(s.config.Assets) {
			continue
		}
		acceptedPrices = append(acceptedPrices, data)
	}

	if len(acceptedPrices) < s.config.MinimalOracles {
		return nil, fmt.Errorf("prices is outdated, err: %w", oracleErr)
	}

	sort.Slice(acceptedPrices, func(i, j int) bool {
		return acceptedPrices[i].Timestamp > acceptedPrices[j].Timestamp
	})

	if len(acceptedPrices) != s.config.MinimalOracles {
		acceptedPrices = acceptedPrices[:s.config.MinimalOracles]
	}

	minTimestamp := acceptedPrices[0].Timestamp
	isOddMinOraclesCount := s.config.MinimalOracles%2 == 1
	medianIndex := s.config.MinimalOracles / 2

	medianPrices := make(map[string]*big.Int, len(s.config.Assets))
	for k, _ := range s.config.Assets {
		sort.SliceStable(acceptedPrices, func(i, j int) bool {
			return acceptedPrices[i].Prices()[k].Cmp(acceptedPrices[j].Prices()[k]) != 1
		})
		if isOddMinOraclesCount {
			medianPrices[k] = acceptedPrices[medianIndex].Prices()[k]
		} else {
			medianPrices[k] = new(big.Int).Div(new(big.Int).Add(acceptedPrices[medianIndex-1].Prices()[k], acceptedPrices[medianIndex].Prices()[k]), big.NewInt(2))
		}
	}

	var packedMedianData *cell.Cell
	for asset, median := range medianPrices {
		packedMedianData = cell.BeginCell().
			MustStoreBigUInt(s.config.Assets[asset].ID, 256).
			MustStoreBigCoins(median).
			MustStoreMaybeRef(packedMedianData).
			EndCell()
	}

	sort.Slice(acceptedPrices, func(i, j int) bool {
		return acceptedPrices[i].oracleID > acceptedPrices[j].oracleID
	})

	var packedOracleData *cell.Cell
	for _, price := range acceptedPrices {
		prf, err := cell.BeginCell().
			MustStoreUInt(uint64(price.Timestamp), 32).
			MustStoreMaybeRef(price.PricesDict.AsCell()).
			EndCell().CreateProof(s.proofSkeleton)
		if err != nil {
			return nil, fmt.Errorf("createProof err: %s", err)
		}

		packedOracleData = cell.BeginCell().
			MustStoreBuilder(cell.BeginCell().
				MustStoreUInt(price.oracleID, 32).
				MustStoreRef(prf).
				MustStoreSlice(price.Signature, 8*uint(len(price.Signature))),
			).MustStoreMaybeRef(packedOracleData).EndCell()
	}

	return &Prices{
		list:         medianPrices,
		data:         cell.BeginCell().MustStoreRef(packedMedianData).MustStoreRef(packedOracleData).EndCell(),
		minTimestamp: minTimestamp,
	}, nil
}
