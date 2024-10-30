package state

import "goplus/avs/config"

type AvsDbState struct{}

func NewAvsDbState(cfg config.Config) (*AvsDbState, error) {
	return &AvsDbState{}, nil
}

var _ AvsStateInterface = (*AvsDbState)(nil)
