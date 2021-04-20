package supernova

import (
	"math/big"
)

type Withdraw struct {
	AmountWithdrew *big.Int
	AmountLeft     *big.Int
	User           string
}

type Deposit struct {
	Amount     *big.Int
	NewBalance *big.Int
	User       string
}

type BaseLog struct {
	LoggedBy         string
	TransactionHash  string
	TransactionIndex int64
	LogIndex         int64
}

type Lock struct {
	BaseLog

	User      string
	Timestamp *big.Int
}

type StakingAction struct {
	BaseLog

	UserAddress  string
	ActionType   ActionType
	Amount       string
	BalanceAfter string
}

type DelegateAction struct {
	BaseLog

	Sender     string
	Receiver   string
	ActionType ActionType
}

type DelegateChange struct {
	BaseLog

	ActionType          ActionType
	Sender              string
	Receiver            string
	Amount              *big.Int `abi:"amount"`
	ToNewDelegatedPower *big.Int `abi:"to_newDelegatedPower"`
}
