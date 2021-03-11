package smartYieldRewards

import (
	"database/sql"
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/barnbridge/barnbridge-backend/state"
	"github.com/barnbridge/barnbridge-backend/types"
)

var log = logrus.WithField("module", "storable(smart yield rewards)")

type Storable struct {
	config      Config
	raw         *types.RawData
	syRewardABI abi.ABI

	processed struct {
		stakingActions []StakingAction
		claims         []ClaimEvent
		blockNumber    int64
		blockTimestamp int64
	}
}

func NewStorable(config Config, raw *types.RawData, syRewardABI abi.ABI) *Storable {
	return &Storable{
		config:      config,
		raw:         raw,
		syRewardABI: syRewardABI,
	}
}

func (s *Storable) ToDB(tx *sql.Tx) error {
	var rewardLogs []web3types.Log

	for _, data := range s.raw.Receipts {
		for _, log := range data.Logs {
			if state.RewardPoolByAddress(log.Address) != nil {
				rewardLogs = append(rewardLogs, log)
			}
		}
	}

	if len(rewardLogs) == 0 {
		log.WithField("handler", "smart yield rewards").Debug("No events found")
		return nil
	}

	err := s.decodeEvents(rewardLogs)
	if err != nil {
		return err
	}

	s.processed.blockNumber, err = strconv.ParseInt(s.raw.Block.Number, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	s.processed.blockTimestamp, err = strconv.ParseInt(s.raw.Block.Timestamp, 0, 64)
	if err != nil {
		return errors.Wrap(err, "could not get block number")
	}

	err = s.storeProcessed(tx)
	if err != nil {
		return err
	}

	return nil
}