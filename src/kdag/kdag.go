package kdag

import (
	"github.com/Kdag-K/evm/src/service"
	"github.com/Kdag-K/evm/src/state"
	kdagconf "github.com/Kdag-K/kdag/src/config"
	"github.com/Kdag-K/kdag/src/kdag"
	"github.com/sirupsen/logrus"
)

// InmemKdag implementes EVM's Consensus interface.
// It uses an inmemory Kdag node.
type InmemKdag struct {
	config      *kdagconf.Config
	kdag       *kdag.Kdag
	ethService *service.Service
	ethState   *state.State
	logger     *logrus.Entry
}

// NewInmemKdag instantiates a new InmemKdag consensus system
func NewInmemKdag(config *kdagconf.Config, logger *logrus.Entry) *InmemKdag {
	return &InmemKdag{
		config: config,
		logger: logger,
	}
}
/*******************************************************************************
IMPLEMENT CONSENSUS INTERFACE
*******************************************************************************/

// Init instantiates a Kdag inmemory node.
//
// XXX - Normally, the Kdag object takes a reference to the InmemProxy via its
// config. Here, we need the InmemProxy to have a reference to the Kdag object
// as well; a sort of circular reference, which is quite ugly. This is necessary
// because the InmemProxy calls the Kdag object directly to retrieve the list
// of validators. We will change this when Blocks are modified to contain the
// validator-set. cf. work on Kdag merkleize branch.
func (ik *InmemKdag) Init(state *state.State, service *service.Service) error {
	ik.ethState = state
	ik.ethService = service
	
	kdag := kdag.NewKdag(ik.config)
	
	inmemProxy := NewInmemProxy(state,
		service,
		kdag,
		service.GetSubmitCh(),
		ik.logger)
	
	ik.config.Proxy = inmemProxy
	
	err := kdag.Init()
	if err != nil {
		return err
	}
	
	ik.kdag = kdag
	
	return nil
}

// Run starts the Kdag node
func (ik *InmemKdag) Run() error {
	ik.kdag.Run()
	return nil
}

// Info returns Kdag stats
func (ik *InmemKdag) Info() (map[string]string, error) {
	info := ik.kdag.Node.GetStats()
	info["type"] = "kdag"
	return info, nil
}