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