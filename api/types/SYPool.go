package types

import (
	"time"

	"github.com/shopspring/decimal"
)

type SYPool struct {
	ProtocolId         string `json:"protocolId"`
	ControllerAddress  string `json:"controllerAddress"`
	ModelAddress       string `json:"modelAddress"`
	ProviderAddress    string `json:"providerAddress"`
	SmartYieldAddress  string `json:"smartYieldAddress"`
	OracleAddress      string `json:"oracleAddress"`
	JuniorBondAddress  string `json:"juniorBondAddress"`
	SeniorBondAddress  string `json:"seniorBondAddress"`
	CTokenAddress      string `json:"cTokenAddress"`
	UnderlyingAddress  string `json:"underlyingAddress"`
	UnderlyingSymbol   string `json:"underlyingSymbol"`
	UnderlyingDecimals int64  `json:"underlyingDecimals"`
	RewardPoolAddress  string `json:"rewardPoolAddress"`

	State SYPoolState `json:"state"`
}

type SYPoolState struct {
	BlockNumber           int64           `json:"blockNumber"`
	BlockTimestamp        time.Time       `json:"blockTimestamp"`
	SeniorLiquidity       decimal.Decimal `json:"seniorLiquidity"`
	JuniorLiquidity       decimal.Decimal `json:"juniorLiquidity"`
	JTokenPrice           decimal.Decimal `json:"jTokenPrice"`
	SeniorAPY             float64         `json:"seniorApy"`
	JuniorAPY             float64         `json:"juniorApy"`
	OriginatorApy         float64         `json:"originatorApy"`
	OriginatorNetApy      float64         `json:"originatorNetApy"`
	AvgSeniorMaturityDays float64         `json:"avgSeniorMaturityDays"`
	NumberOfSeniors       int64           `json:"numberOfSeniors"`
	NumberOfJuniors       int64           `json:"numberOfJuniors"`
	NumberOfJuniorsLocked int64           `json:"numberOfJuniorsLocked"`
	JuniorLiquidityLocked decimal.Decimal `json:"juniorLiquidityLocked"`
}

type SYRewardPool struct {
	PoolAddress        string `json:"poolAddress"`
	PoolTokenAddress   string `json:"poolTokenAddress"`
	RewardTokenAddress string `json:"rewardTokenAddress"`
	PoolTokenDecimals  int64  `json:"poolTokenDecimals"`
	ProtocolID         string `json:"protocolId"`
	UnderlyingSymbol   string `json:"underlyingSymbol"`
}
