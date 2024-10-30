package handlers

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"goplus/mock_secware/pkg/config"
)

type SecwareSetting struct {
	PkGateway      *ecdsa.PublicKey
	SkSecwareBytes []byte
	SecwareId      int
	SecwareVersion int
}

func BuildContext(ctx context.Context) (context.Context, error) {
	skGateway, err := crypto.HexToECDSA(config.SkGatewayHex)
	if err != nil {
		return nil, fmt.Errorf("bad gateway private key")
	}

	skSecwareBytes, err := hex.DecodeString(config.SkSecwareHex)
	if err != nil {
		return nil, fmt.Errorf("bad secware private key")
	}

	setting := SecwareSetting{
		PkGateway:      skGateway.Public().(*ecdsa.PublicKey),
		SkSecwareBytes: skSecwareBytes,
		SecwareId:      config.SecwareId,
		SecwareVersion: config.SecwareVersion,
	}
	return context.WithValue(ctx, "Setting", &setting), nil
}
