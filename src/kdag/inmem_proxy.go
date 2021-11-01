package kdag

import (
	"github.com/Kdag-K/evm/src/service"
	"github.com/Kdag-K/evm/src/state"
	"github.com/Kdag-K/kdag/src/crypto/keys"
	"github.com/Kdag-K/kdag/src/hashgraph"
	"github.com/Kdag-K/kdag/src/kdag"
	"github.com/Kdag-K/kdag/src/proxy"
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
// processInternalTransactions decides if InternalTransactions should be
// accepted. For PEER_ADD transactions, it checks if the peer is authorised in
// the POA smart-contract. All PEER_REMOVE transactions are accepted.
func (p *InmemProxy) processInternalTransactions(internalTransactions []hashgraph.InternalTransaction) []hashgraph.InternalTransactionReceipt {
	receipts := []hashgraph.InternalTransactionReceipt{}
	
	for _, tx := range internalTransactions {
		switch tx.Body.Type {
		case hashgraph.PEER_ADD:
			pk, err := crypto.UnmarshalPubkey(tx.Body.Peer.PubKeyBytes())
			if err != nil {
				p.logger.Warningf("couldn't unmarshal pubkey bytes: %v", err)
			}
			
			addr := crypto.PubkeyToAddress(*pk)
			
			ok, err := p.state.CheckAuthorised(addr)
			
			if err != nil {
				p.logger.WithError(err).Error("Error in checkAuthorised")
				receipts = append(receipts, tx.AsRefused())
			} else {
				if ok {
					p.logger.WithField("addr", addr.String()).Info("Accepted peer")
					receipts = append(receipts, tx.AsAccepted())
				} else {
					p.logger.WithField("addr", addr.String()).Info("Rejected peer")
					receipts = append(receipts, tx.AsRefused())
				}
			}
		case hashgraph.PEER_REMOVE:
			receipts = append(receipts, tx.AsAccepted())
		}
	}
	
	return receipts
}

// processEvictions compares the current validator-set to the whitelist and
// creates InternalTransactionReceipts to evict any current validator which is
// not in the whitelist.
func (p *InmemProxy) processEvictions(block hashgraph.Block) []hashgraph.InternalTransactionReceipt {
	receipts := []hashgraph.InternalTransactionReceipt{}
	
	if p.kdag != nil {
		kdagValidators, err := p.kdag.Node.GetValidatorSet(block.RoundReceived())
		if err != nil {
			p.logger.WithError(err).Error("Error GetValidatorSet")
			return receipts
		}
		
		for _, val := range kdagValidators {
			pk, err := crypto.UnmarshalPubkey(val.PubKeyBytes())
			if err != nil {
				p.logger.Warningf("couldn't unmarshal pubkey bytes: %v", err)
				continue
			}
			
			addr := crypto.PubkeyToAddress(*pk)
			
			ok, err := p.state.CheckAuthorised(addr)
			
			if err != nil {
				p.logger.WithError(err).Error("Error in checkAuthorised")
			} else {
				if !ok {
					p.logger.WithField("addr", addr.String()).Info("Ejected peer")
					receipts = append(receipts,
						hashgraph.InternalTransactionReceipt{
							InternalTransaction: hashgraph.InternalTransaction{
								Body: hashgraph.InternalTransactionBody{
									Type: hashgraph.PEER_REMOVE,
									Peer: *val,
								},
							},
							Accepted: true,
						})
				}
			}
		}
	}
	
	return receipts
}
// CommitBlock applies the block's transactions to the state and commits. All
// transaction fees are sent to the coinbase address, which is computed from the
// block and the current validator-set. It also checks the block's internal
// transactions against the POA smart-contract to verify if joining peers are
// authorised to become validators in kdag. It returns the resulting
// state-hash and internal transaction receips.
func (p *InmemProxy) CommitBlock(block hashgraph.Block) (proxy.CommitResponse, error) {
	
	coinbaseAddress, err := p.getCoinbase(block)
	if err != nil {
		return proxy.CommitResponse{}, err
	}
	
	p.logger.WithFields(logrus.Fields{
		"coinbase": coinbaseAddress.String(),
		"block":    block.Index(),
	}).Info("Commit")
	
	blockHashBytes, err := block.Hash()
	blockHash := ethCommon.BytesToHash(blockHashBytes)
	
	for i, tx := range block.Transactions() {
		if err := p.state.ApplyTransaction(tx, i, blockHash, coinbaseAddress); err != nil {
			p.logger.WithError(err).Errorf("Failed to apply tx %d of %d", i+1, len(block.Transactions()))
		}
	}
	
	hash, err := p.state.Commit()
	if err != nil {
		return proxy.CommitResponse{}, err
	}
	
	internalTransactionReceipts := p.processInternalTransactions(block.InternalTransactions())
	
	evictionReceipts := p.processEvictions(block)
	
	receipts := append(internalTransactionReceipts, evictionReceipts...)
	
	res := proxy.CommitResponse{
		StateHash:                   hash.Bytes(),
		InternalTransactionReceipts: receipts,
	}
	
	return res, nil
}
//TODO - Implement these two functions
//GetSnapshot will generate a snapshot
func (p *InmemProxy) GetSnapshot(blockIndex int) ([]byte, error) {
	return []byte{}, nil
}

//Restore will restore a snapshot
func (p *InmemProxy) Restore(snapshot []byte) error {
	return nil
}