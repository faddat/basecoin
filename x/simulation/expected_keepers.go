package simulation

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type (
	// AccountKeeper defines the expected account keeper used for simulations (noalias)
	AccountKeeper interface {
		GetAccount(ctx sdk.Context, addr sdk.AccAddress) sdk.AccountI
	}

	// BankKeeper defines the expected interface needed to retrieve account balances.
	BankKeeper interface {
		SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	}
)
