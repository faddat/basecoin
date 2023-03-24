package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the staking Querier
const (
	QueryValidators                    = "validators"
	QueryValidator                     = "validator"
	QueryDelegatorDelegations          = "delegatorDelegations"
	QueryDelegatorUnbondingDelegations = "delegatorUnbondingDelegations"
	QueryRedelegations                 = "redelegations"
	QueryValidatorDelegations          = "validatorDelegations"
	QueryValidatorRedelegations        = "validatorRedelegations"
	QueryValidatorUnbondingDelegations = "validatorUnbondingDelegations"
	QueryDelegation                    = "delegation"
	QueryUnbondingDelegation           = "unbondingDelegation"
	QueryDelegatorValidators           = "delegatorValidators"
	QueryDelegatorValidator            = "delegatorValidator"
	QueryPool                          = "pool"
	QueryParameters                    = "parameters"
	QueryHistoricalInfo                = "historicalInfo"
)

type (
	// defines the params for the following queries:
	// - 'custom/staking/delegatorDelegations'
	// - 'custom/staking/delegatorUnbondingDelegations'
	// - 'custom/staking/delegatorValidators'
	QueryDelegatorParams struct {
		DelegatorAddr sdk.AccAddress
	}

	// QueryValidatorsParams defines the params for the following queries:
	// - 'custom/staking/validators'
	QueryValidatorsParams struct {
		Page, Limit int
		Status      string
	}

	// defines the params for the following queries:
	// - 'custom/staking/validator'
	// - 'custom/staking/validatorDelegations'
	// - 'custom/staking/validatorUnbondingDelegations'
	QueryValidatorParams struct {
		ValidatorAddr sdk.ValAddress
		Page, Limit   int
	}

	// defines the params for the following queries:
	// - 'custom/staking/redelegation'
	QueryRedelegationParams struct {
		DelegatorAddr    sdk.AccAddress
		SrcValidatorAddr sdk.ValAddress
		DstValidatorAddr sdk.ValAddress
	}
)

func NewQueryDelegatorParams(delegatorAddr sdk.AccAddress) QueryDelegatorParams {
	return QueryDelegatorParams{
		DelegatorAddr: delegatorAddr,
	}
}

func NewQueryValidatorParams(validatorAddr sdk.ValAddress, page, limit int) QueryValidatorParams {
	return QueryValidatorParams{
		ValidatorAddr: validatorAddr,
		Page:          page,
		Limit:         limit,
	}
}

func NewQueryRedelegationParams(delegatorAddr sdk.AccAddress, srcValidatorAddr, dstValidatorAddr sdk.ValAddress) QueryRedelegationParams {
	return QueryRedelegationParams{
		DelegatorAddr:    delegatorAddr,
		SrcValidatorAddr: srcValidatorAddr,
		DstValidatorAddr: dstValidatorAddr,
	}
}

func NewQueryValidatorsParams(page, limit int, status string) QueryValidatorsParams {
	return QueryValidatorsParams{page, limit, status}
}
