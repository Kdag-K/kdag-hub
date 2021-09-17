package kdag

import (
	"github.com/Kdag-K/kdag/src/kdag"
	kdagconf "github.com/Kdag-K/kdag/src/config"
	"github.com/Kdag-K/evm/src/service"
	"github.com/Kdag-K/evm/src/state"
	"github.com/sirupsen/logrus"
)

// InmemBabble implementes EVM's Consensus interface.
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