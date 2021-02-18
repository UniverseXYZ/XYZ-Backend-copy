package smartYield

import (
	"database/sql"
	"encoding/hex"
	"math/big"

	web3types "github.com/alethio/web3-go/types"
	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/barnbridge/barnbridge-backend/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

type SeniorBondBuyTrade struct {
	*types.Event

	SYAddress    string
	BuyerAddress string
	SeniorBondID *big.Int
	UnderlyingIn *big.Int
	Gain         *big.Int
	ForDays      *big.Int
}

type SeniorBondRedeemTrade struct {
	*types.Event

	SYAddress    string
	OwnerAddress string
	SeniorBondID *big.Int
	Fee          *big.Int
}

func (s *Storable) decodeSeniorBondBuyEvent(log web3types.Log, event string) (*SeniorBondBuyTrade, error) {
	var t SeniorBondBuyTrade
	t.SYAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(&t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.BuyerAddress = utils.Topic2Address(log.Topics[1])

	n, ok := new(big.Int).SetString(utils.Trim0x(log.Topics[2]), 16)
	if !ok {
		return nil, errors.New("could not convert seniorBondId ")
	}

	t.SeniorBondID = n
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) decodeSeniorBondRedeemEvent(log web3types.Log, event string) (*SeniorBondRedeemTrade, error) {
	var t SeniorBondRedeemTrade
	t.SYAddress = utils.NormalizeAddress(log.Address)

	data, err := hex.DecodeString(utils.Trim0x(log.Data))
	if err != nil {
		return nil, errors.Wrap(err, "could not decode log data")
	}

	err = s.abis["smartyield"].UnpackIntoInterface(&t, event, data)
	if err != nil {
		return nil, errors.Wrap(err, "could not unpack log data")
	}

	t.OwnerAddress = utils.Topic2Address(log.Topics[1])

	n, ok := new(big.Int).SetString(utils.Trim0x(log.Topics[2]), 16)
	if !ok {
		return nil, errors.New("could not convert seniorBondId ")
	}

	t.SeniorBondID = n
	t.Event, err = new(types.Event).Build(log)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *Storable) storeSeniorBuyTrades(tx *sql.Tx) error {
	if len(s.processed.seniorActions.seniorBondBuys) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_senior_buy", "sy_address", "buyer_address", "senior_bond_id", "underlying_in", "gain", "for_days", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.seniorActions.seniorBondBuys {
		_, err = stmt.Exec(a.SYAddress, a.BuyerAddress, a.SeniorBondID.String(), a.UnderlyingIn.String(), a.Gain.String(), a.ForDays.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Storable) storeSeniorRedeemTrades(tx *sql.Tx) error {
	if len(s.processed.seniorActions.seniorBondRedeems) == 0 {
		return nil
	}

	stmt, err := tx.Prepare(pq.CopyIn("smart_yield_senior_redeem", "sy_address", "owner_address", "senior_bond_id", "fee", "tx_hash", "tx_index", "log_index", "block_timestamp", "included_in_block"))
	if err != nil {
		return err
	}

	for _, a := range s.processed.seniorActions.seniorBondRedeems {
		_, err = stmt.Exec(a.SYAddress, a.OwnerAddress, a.SeniorBondID.String(), a.Fee.String(), a.TransactionHash, a.TransactionIndex, a.LogIndex, s.processed.blockTimestamp, s.processed.blockNumber)
		if err != nil {
			return err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	err = stmt.Close()
	if err != nil {
		return err
	}

	return nil
}
