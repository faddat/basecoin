package paytovote

import (
	"fmt"

	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/basecoin/state"
	"github.com/tendermint/basecoin/types"
	"github.com/tendermint/go-wire"
)

type P2VPlugin struct {
	name string
}

func New() *P2VPlugin {
	return &P2VPlugin{
		name: "paytovote",
	}
}

///////////////////////////////////////////////////

const (
	TypeByteTxCreate byte = 0x01
	TypeByteTxVote   byte = 0x02

	TypeByteVoteFor     byte = 0x01
	TypeByteVoteAgainst byte = 0x02
)

type createIssueTx struct {
	Issue           string      //Issue to be created
	FeePerVote      types.Coins //Cost to vote for the issue
	Fee2CreateIssue types.Coins //Cost to create a new issue
}

type voteTx struct {
	Issue        string //Issue being voted for
	VoteTypeByte byte   //How is the vote being cast
}

func NewCreateIssueTxBytes(issue string, feePerVote, fee2CreateIssue types.Coins) []byte {
	data := wire.BinaryBytes(
		createIssueTx{
			Issue:           issue,
			FeePerVote:      feePerVote,
			Fee2CreateIssue: fee2CreateIssue,
		})
	data = append([]byte{TypeByteTxCreate}, data...)
	return data
}

func NewVoteTxBytes(issue string, voteTypeByte byte) []byte {
	data := wire.BinaryBytes(
		voteTx{
			Issue:        issue,
			VoteTypeByte: voteTypeByte,
		})
	data = append([]byte{TypeByteTxVote}, data...)
	return data
}

///////////////////////////////////////////////////

type P2VIssue struct {
	Issue        string
	FeePerVote   types.Coins
	votesFor     int
	votesAgainst int
}

func newP2VIssue(issue string, feePerVote types.Coins) P2VIssue {
	return P2VIssue{
		Issue:        issue,
		FeePerVote:   feePerVote,
		votesFor:     0,
		votesAgainst: 0,
	}
}

func (p2v *P2VPlugin) IssueKey(issue string) []byte {
	//The state key is defined as only being affected by effected issue
	// aka. if multiple paytovote plugins are initialized
	// then all will have access to the same issue vote counts
	return []byte(fmt.Sprintf("P2VPlugin{issue=%v}.State", issue))
}

///////////////////////////////////////////////////

func (p2v *P2VPlugin) Name() string {
	return p2v.name
}

func (p2v *P2VPlugin) SetOption(store types.KVStore, key string, value string) (log string) {
	return ""
}

func (p2v *P2VPlugin) RunTx(store types.KVStore, ctx types.CallContext, txBytes []byte) (res abci.Result) {

	defer func() {
		//Return the ctx coins to the wallet if there is an error
		if res.IsErr() {
			acc := ctx.CallerAccount
			acc.Balance = acc.Balance.Plus(ctx.Coins)       // add the context transaction coins
			state.SetAccount(store, ctx.CallerAddress, acc) // save the new balance
		}
	}()

	//Determine the transaction type and then send to the appropriate transaction function
	if len(txBytes) < 1 {
		return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: no tx bytes")
	}

	//Note that the zero position of txBytes contains the type-byte for the tx type
	switch txBytes[0] {
	case TypeByteTxCreate:
		return p2v.runTxCreateIssue(store, ctx, txBytes[1:])
	case TypeByteTxVote:
		return p2v.runTxVote(store, ctx, txBytes[1:])
	default:
		return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: bad prepended bytes")
	}
}

func chargeFee(store types.KVStore, ctx types.CallContext, fee types.Coins) {

	//Charge the Fee from the context coins
	leftoverCoins := ctx.Coins.Minus(fee)
	if !leftoverCoins.IsZero() {
		lc := "leftoverCoins: "
		for i := 0; i < len(leftoverCoins); i++ {
			lc += " " + leftoverCoins[i].String()
		}
		fmt.Println(lc)
		lc = "fee: "
		for i := 0; i < len(fee); i++ {
			lc += " " + fee[i].String()
		}
		fmt.Println(lc)

		acc := ctx.CallerAccount
		lc = "acc b4: "
		for i := 0; i < len(acc.Balance); i++ {
			lc += " " + acc.Balance[i].String()
		}
		fmt.Println(lc)

		//return leftover coins
		acc.Balance = acc.Balance.Plus(leftoverCoins)   // subtract fees
		state.SetAccount(store, ctx.CallerAddress, acc) // save the new balance
		lc = "acc aftr: "
		for i := 0; i < len(acc.Balance); i++ {
			lc += " " + acc.Balance[i].String()
		}
		fmt.Println(lc)

	}
}

func (p2v *P2VPlugin) runTxCreateIssue(store types.KVStore, ctx types.CallContext, txBytes []byte) (res abci.Result) {

	// Decode tx
	var tx createIssueTx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}

	//Validate Tx
	switch {
	case len(tx.Issue) == 0:
		return abci.ErrInternalError.AppendLog("P2VTx.Issue must have a length greater than 0")
	case !tx.FeePerVote.IsValid():
		return abci.ErrInternalError.AppendLog("P2VTx.Fee2Vote is not sorted or has zero amounts")
	case !tx.FeePerVote.IsNonnegative():
		return abci.ErrInternalError.AppendLog("P2VTx.Fee2Vote must be nonnegative")
	case !tx.Fee2CreateIssue.IsValid():
		return abci.ErrInternalError.AppendLog("P2VTx.Fee2CreateIssue is not sorted or has zero amounts")
	case !tx.Fee2CreateIssue.IsNonnegative():
		return abci.ErrInternalError.AppendLog("P2VTx.Fee2CreateIssue must be nonnegative")
	case !ctx.Coins.IsGTE(tx.Fee2CreateIssue): // Did the caller provide enough coins?
		return abci.ErrInsufficientFunds.AppendLog("Tx Funds insufficient for creating a new issue")
	}

	// Load P2VIssue
	var p2vIssue P2VIssue
	p2vIssueBytes := store.Get(p2v.IssueKey(tx.Issue))

	//Return if the issue already exists
	if len(p2vIssueBytes) > 0 {
		return abci.ErrInsufficientFunds.AppendLog("Cannot create an already existing issue")
	}

	// Create and Save P2VIssue, charge fee, return
	newP2VIssue := newP2VIssue(tx.Issue, tx.FeePerVote)
	store.Set(p2v.IssueKey(tx.Issue), wire.BinaryBytes(newP2VIssue))
	chargeFee(store, ctx, tx.Fee2CreateIssue)
	return abci.NewResultOK(wire.BinaryBytes(p2vIssue), "")
}

func (p2v *P2VPlugin) runTxVote(store types.KVStore, ctx types.CallContext, txBytes []byte) (res abci.Result) {

	// Decode tx
	var tx voteTx
	err := wire.ReadBinaryBytes(txBytes, &tx)
	if err != nil {
		return abci.ErrBaseEncodingError.AppendLog("Error decoding tx: " + err.Error())
	}

	//Validate Tx
	if len(tx.Issue) == 0 {
		return abci.ErrInternalError.AppendLog("transaction issue must have a length greater than 0")
	}

	// Load P2VIssue
	var p2vIssue P2VIssue
	p2vIssueBytes := store.Get(p2v.IssueKey(tx.Issue))

	//Determine if the issue already exists and load
	if len(p2vIssueBytes) > 0 { //is there a record of the issue existing?
		err = wire.ReadBinaryBytes(p2vIssueBytes, &p2vIssue)
		if err != nil {
			return abci.ErrInternalError.AppendLog("Error decoding state: " + err.Error())
		}
	} else {
		return abci.ErrInsufficientFunds.AppendLog("Tx Issue not found")
	}

	// Did the caller provide enough coins?
	if !ctx.Coins.IsGTE(p2vIssue.FeePerVote) {
		return abci.ErrInsufficientFunds.AppendLog("Tx Funds insufficient for voting")
	}

	//Transaction Logic
	switch tx.VoteTypeByte {
	case TypeByteVoteFor:
		p2vIssue.votesFor += 1
	case TypeByteVoteAgainst:
		p2vIssue.votesAgainst += 1
	default:
		return abci.ErrInternalError.AppendLog("P2VTx.ActionTypeByte was not recognized")
	}

	// Save P2VIssue, charge fee, return
	store.Set(p2v.IssueKey(tx.Issue), wire.BinaryBytes(p2vIssue))
	chargeFee(store, ctx, p2vIssue.FeePerVote)
	return abci.NewResultOK(wire.BinaryBytes(p2vIssue), "")
}

func (p2v *P2VPlugin) InitChain(store types.KVStore, vals []*abci.Validator) {}
func (p2v *P2VPlugin) BeginBlock(store types.KVStore, height uint64)         {}
func (p2v *P2VPlugin) EndBlock(store types.KVStore, height uint64) []*abci.Validator {
	return nil
}
