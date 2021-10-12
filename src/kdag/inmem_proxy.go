package kdag

import (
	"github.com/Kdag-K/evm/src/service"
	"github.com/Kdag-K/evm/src/state"
	"github.com/Kdag-K/kdag/src/crypto/keys"
	"github.com/Kdag-K/kdag/src/hashgraph"
	"github.com/Kdag-K/kdag/src/kdag"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

// InmemProxy implements the kdag AppProxy interface
type InmemProxy struct {
	service  *service.Service
	state    *state.State
	kdag     *kdag.Kdag
	submitCh chan []byte
	logger   *logrus.Entry
}

// NewInmemProxy initializes and return a new InmemProxy
func NewInmemProxy(state *state.State,
	service *service.Service,
	kdag *kdag.Kdag,
	submitCh chan []byte,
	logger *logrus.Entry) *InmemProxy {

	return &InmemProxy{
		service:  service,
		state:    state,
		kdag:     kdag,
		submitCh: submitCh,
		logger:   logger,
	}
}

// SubmitCh is the channel through which the Service sends
// transactions to the node.
func (p *InmemProxy) SubmitCh() chan []byte {
	return p.submitCh
}

// getCoinbase returns the coinbase address which will receive all the
// transaction fees from the block. It is meant to be a safe and fair selection
// process from the current Kdag validator-set. We use the block hash, which
// is pseudo-random, but equal for all validators, to select a validator from
// the current validator-set.
func (p *InmemProxy) getCoinbase(block hashgraph.Block) (ethCommon.Address, error) {
	coinbaseAddress := ethCommon.Address{}

	if p.kdag != nil {
		kdagValidators, err := p.kdag.Node.GetValidatorSet(block.RoundReceived())
		if err != nil {
			return coinbaseAddress, err
		}

		blockHash, _ := block.Hash()
		blockRand := keys.PublicKeyID(blockHash)

		coinbaseValidator := kdagValidators[blockRand%uint32(len(kdagValidators))]

		coinbasePubKey, err := crypto.UnmarshalPubkey(coinbaseValidator.PubKeyBytes())
		if err != nil {
			return coinbaseAddress, err
		}

		coinbaseAddress = crypto.PubkeyToAddress(*coinbasePubKey)
	}

	return coinbaseAddress, nil
}
