package api

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"

	"github.com/barnbridge/barnbridge-backend/api/types"
	"github.com/barnbridge/barnbridge-backend/utils"
)

func (a *API) handlePoolDetails(c *gin.Context) {
	pool := c.Param("address")

	poolAddress, err := utils.ValidateAccount(pool)
	if err != nil {
		BadRequest(c, errors.New("invalid pool address"))
		return
	}

	var p types.SYPool

	err = a.db.QueryRow(`
		select protocol_id,
			   controller_address,
			   model_address,
			   provider_address,
			   sy_address,
			   oracle_address,
			   junior_bond_address,
			   senior_bond_address,
			   receipt_token_address,
			   underlying_address,
			   underlying_symbol,
			   underlying_decimals
		from smart_yield_pools p
		where sy_address = $1
	`, poolAddress).Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals)
	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}
	if err == sql.ErrNoRows {
		NotFound(c)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	var state types.SYPoolState
	err = a.db.QueryRow(`
			select included_in_block,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   (select count(distinct buyer_address) from smart_yield_senior_buy where sy_address = pool_address ) as number_of_seniors,
				   coalesce((select sum(for_days*underlying_in)/sum(underlying_in) from smart_yield_senior_buy where sy_address = pool_address), 0) as avg_senior_buy,
				   (select count(distinct buyer_address) from smart_yield_token_buy where sy_address = pool_address ) as number_of_juniors,
					( select sum(case 
					    when (select count(*) from smart_yield_junior_redeem as r where r.junior_bond_address = b.junior_bond_address
								 																and r.junior_bond_id = b.junior_bond_id) = 0 then tokens_in else 0
				   		end)
					from smart_yield_junior_buy as b
					) as junior_liquidity_locked
			from smart_yield_state
			where pool_address = $1
			order by included_in_block desc
			limit 1
		`, p.SmartYieldAddress).Scan(&state.BlockNumber, &state.BlockTimestamp, &state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY, &state.OriginatorApy, &state.OriginatorNetApy, &state.NumberOfSeniors, &state.AvgSeniorMaturityDays, &state.NumberOfJuniors, &state.JuniorLiquidityLocked)
	if err != nil {
		Error(c, err)
		return
	}

	tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

	state.JuniorLiquidityLocked = state.JuniorLiquidityLocked.Div(tenPowDec)
	state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
	state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
	state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

	p.State = state

	OK(c, p)
}

func (a *API) handlePools(c *gin.Context) {
	protocols := strings.ToLower(c.DefaultQuery("originator", "all"))

	protocolsArray := strings.Split(protocols, ",")

	var pools []types.SYPool

	query := `
		select protocol_id,
			   controller_address,
			   model_address,
			   provider_address,
			   sy_address,
			   oracle_address,
			   junior_bond_address,
			   senior_bond_address,
			   receipt_token_address,
			   underlying_address,
			   underlying_symbol,
			   underlying_decimals
		from smart_yield_pools p
		where 1 = 1 %s
	`

	var parameters []interface{}

	if protocols == "all" {
		query = fmt.Sprintf(query, "")
	} else {
		protocolFilter := fmt.Sprintf("and protocol_id = ANY($1)")
		parameters = append(parameters, pq.Array(protocolsArray))
		query = fmt.Sprintf(query, protocolFilter)
	}
	rows, err := a.db.Query(query, parameters...)

	if err != nil && err != sql.ErrNoRows {
		Error(c, err)
		return
	}

	tenPow18 := decimal.NewFromInt(10).Pow(decimal.NewFromInt(18))

	for rows.Next() {
		var p types.SYPool

		err := rows.Scan(&p.ProtocolId, &p.ControllerAddress, &p.ModelAddress, &p.ProviderAddress, &p.SmartYieldAddress, &p.OracleAddress, &p.JuniorBondAddress, &p.SeniorBondAddress, &p.CTokenAddress, &p.UnderlyingAddress, &p.UnderlyingSymbol, &p.UnderlyingDecimals)
		if err != nil {
			Error(c, err)
			return
		}

		var state types.SYPoolState
		err = a.db.QueryRow(`
			select included_in_block,
				   block_timestamp,
				   senior_liquidity,
				   junior_liquidity,
				   jtoken_price,
				   senior_apy,
				   junior_apy,
				   originator_apy,
				   originator_net_apy,
				   (select count(distinct buyer_address) from smart_yield_senior_buy where sy_address = pool_address ) as number_of_seniors,
				   coalesce((select sum(for_days*underlying_in)/sum(underlying_in) from smart_yield_senior_buy where sy_address = pool_address), 0) as avg_senior_buy,
				   (select count(distinct buyer_address) from smart_yield_token_buy where sy_address = pool_address ) as number_of_juniors,
					( select sum(case 
					    when (select count(*) from smart_yield_junior_redeem as r where r.junior_bond_address = b.junior_bond_address
								 																and r.junior_bond_id = b.junior_bond_id) = 0 then tokens_in else 0
				   		end)
					from smart_yield_junior_buy as b
					) as junior_liquidity_locked
			from smart_yield_state
			where pool_address = $1
			order by included_in_block desc
			limit 1
		`, p.SmartYieldAddress).Scan(&state.BlockNumber, &state.BlockTimestamp, &state.SeniorLiquidity, &state.JuniorLiquidity, &state.JTokenPrice, &state.SeniorAPY, &state.JuniorAPY, &state.OriginatorApy, &state.OriginatorNetApy, &state.NumberOfSeniors, &state.AvgSeniorMaturityDays, &state.NumberOfJuniors, &state.JuniorLiquidityLocked)
		if err != nil && err != sql.ErrNoRows {
			Error(c, err)
			return
		}

		tenPowDec := decimal.NewFromInt(10).Pow(decimal.NewFromInt(p.UnderlyingDecimals))

		state.JuniorLiquidityLocked = state.JuniorLiquidityLocked.Div(tenPowDec)
		state.JTokenPrice = state.JTokenPrice.DivRound(tenPow18, 18)
		state.SeniorLiquidity = state.SeniorLiquidity.Div(tenPowDec)
		state.JuniorLiquidity = state.JuniorLiquidity.Div(tenPowDec)

		p.State = state

		pools = append(pools, p)
	}

	OK(c, pools)
}
