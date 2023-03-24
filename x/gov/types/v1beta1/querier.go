package v1beta1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// query endpoints supported by the governance Querier
const (
	QueryParams    = "params"
	QueryProposals = "proposals"
	QueryProposal  = "proposal"
	QueryDeposits  = "deposits"
	QueryDeposit   = "deposit"
	QueryVotes     = "votes"
	QueryVote      = "vote"
	QueryTally     = "tally"

	ParamDeposit  = "deposit"
	ParamVoting   = "voting"
	ParamTallying = "tallying"
)

type (
	// QueryProposalParams Params for queries:
	// - 'custom/gov/proposal'
	// - 'custom/gov/deposits'
	// - 'custom/gov/tally'
	QueryProposalParams struct {
		ProposalID uint64
	}

	// QueryProposalVotesParams used for queries to 'custom/gov/votes'.
	QueryProposalVotesParams struct {
		ProposalID uint64
		Page       int
		Limit      int
	}

	// QueryDepositParams params for query 'custom/gov/deposit'
	QueryDepositParams struct {
		ProposalID uint64
		Depositor  sdk.AccAddress
	}

	// QueryProposalsParams Params for query 'custom/gov/proposals'
	QueryProposalsParams struct {
		Page           int
		Limit          int
		Voter          sdk.AccAddress
		Depositor      sdk.AccAddress
		ProposalStatus ProposalStatus
	}

	// QueryVoteParams Params for query 'custom/gov/vote'
	QueryVoteParams struct {
		ProposalID uint64
		Voter      sdk.AccAddress
	}
)

// NewQueryProposalParams creates a new instance of QueryProposalParams
func NewQueryProposalParams(proposalID uint64) QueryProposalParams {
	return QueryProposalParams{
		ProposalID: proposalID,
	}
}

// NewQueryProposalVotesParams creates new instance of the QueryProposalVotesParams.
func NewQueryProposalVotesParams(proposalID uint64, page, limit int) QueryProposalVotesParams {
	return QueryProposalVotesParams{
		ProposalID: proposalID,
		Page:       page,
		Limit:      limit,
	}
}

// NewQueryDepositParams creates a new instance of QueryDepositParams
func NewQueryDepositParams(proposalID uint64, depositor sdk.AccAddress) QueryDepositParams {
	return QueryDepositParams{
		ProposalID: proposalID,
		Depositor:  depositor,
	}
}

// NewQueryVoteParams creates a new instance of QueryVoteParams
func NewQueryVoteParams(proposalID uint64, voter sdk.AccAddress) QueryVoteParams {
	return QueryVoteParams{
		ProposalID: proposalID,
		Voter:      voter,
	}
}

// NewQueryProposalsParams creates a new instance of QueryProposalsParams
func NewQueryProposalsParams(page, limit int, status ProposalStatus, voter, depositor sdk.AccAddress) QueryProposalsParams {
	return QueryProposalsParams{
		Page:           page,
		Limit:          limit,
		Voter:          voter,
		Depositor:      depositor,
		ProposalStatus: status,
	}
}
