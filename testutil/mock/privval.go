package mock

import (
	cmtproto "github.com/cometbft/cometbft/api/cometbft/types/v2"
	"github.com/cometbft/cometbft/v2/crypto"
	cmttypes "github.com/cometbft/cometbft/v2/types"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

var _ cmttypes.PrivValidator = PV{}

// PV implements PrivValidator without any safety or persistence.
// Only use it for testing.
type PV struct {
	PrivKey cryptotypes.PrivKey
}

func NewPV() PV {
	return PV{ed25519.GenPrivKey()}
}

// GetPubKey implements PrivValidator interface
func (pv PV) GetPubKey() (crypto.PubKey, error) {
	return cryptocodec.ToCmtPubKeyInterface(pv.PrivKey.PubKey())
}

// SignVote implements PrivValidator interface
func (pv PV) SignVote(chainID string, vote *cmtproto.Vote, signExtension bool) error {
	signBytes := cmttypes.VoteSignBytes(chainID, vote)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig
	if signExtension {
		extSignBytes, _ := cmttypes.VoteExtensionSignBytes(chainID, vote)
		extSig, err := pv.PrivKey.Sign(extSignBytes)
		if err != nil {
			return err
		}
		vote.ExtensionSignature = extSig
	}
	return nil
}

// SignProposal implements PrivValidator interface
func (pv PV) SignProposal(chainID string, proposal *cmtproto.Proposal) error {
	signBytes := cmttypes.ProposalSignBytes(chainID, proposal)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	proposal.Signature = sig
	return nil
}

func (pv PV) SignBytes(bytes []byte) ([]byte, error) {
	return pv.PrivKey.Sign(bytes)
}
