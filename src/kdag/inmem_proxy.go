package kdag

import (
	"github.com/Kdag-K/kdag/src/kdag"
	"github.com/Kdag-K/evm/src/service"
	"github.com/Kdag-K/evm/src/state"
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