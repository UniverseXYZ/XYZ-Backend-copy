package types

import (
	"strconv"

	web3types "github.com/alethio/web3-go/types"
	"github.com/pkg/errors"
)

type Event struct {
	LoggedBy         string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

func (e *Event) Build(log web3types.Log) (*Event, error) {
	var err error
	e.LoggedBy = log.Address
	e.TransactionHash = log.TransactionHash

	e.TransactionIndex, err = strconv.ParseInt(log.TransactionIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert transactionIndex from bond contract to int64")
	}

	e.LogIndex, err = strconv.ParseInt(log.LogIndex, 0, 64)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert logIndex from  bond contract to int64")
	}

	return e, nil
}
