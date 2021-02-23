package smartYieldState

import (
	"database/sql"
	"strconv"
	"sync"

	"github.com/alethio/web3-go/ethrpc"
	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
)

type Config struct {
	ComptrollerAddress string
}

type Storable struct {
	config Config
	raw    *types.RawData

	abis map[string]abi.ABI
	eth  *ethrpc.ETH

	Preprocessed struct {
		BlockTimestamp int64
		BlockNumber    int64
	}
}

func New(config Config, raw *types.RawData, abis map[string]abi.ABI, eth *ethrpc.ETH) (*Storable, error) {
	var s = &Storable{
		config: config,
		raw:    raw,
		abis:   abis,
		eth:    eth,
	}

	var err error
	s.Preprocessed.BlockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "unable to process block number")
	}

	s.Preprocessed.BlockTimestamp, err = strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not parse block timestamp")
	}

	return s, nil
}

func (s Storable) ToDB(tx *sql.Tx) error {
	var wg = &errgroup.Group{}

	var results = make(map[string]*State)
	var mu = &sync.Mutex{}

	for _, p := range state.Pools() {
		p := p

		results[p.SmartYieldAddress] = &State{
			PoolAddress: p.SmartYieldAddress,
		}

		s.getTotalLiquidity(wg, p, mu, results)
		s.getJuniorLiquidity(wg, p, mu, results)
		s.getPrice(wg, p, mu, results)

		if p.ProtocolId == "compound/v2" {
			s.getCompoundAPY(wg, p, mu, results)
		}
	}

	err := wg.Wait()
	if err != nil {
		return err
	}

	spew.Dump(results)

	return nil
}