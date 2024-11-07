package signature

import (
	"encoding/hex"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"goplus/shared/pkg/types"
	"testing"
)

func TestBLS(t *testing.T) {
	operatorAddr, err := hex.DecodeString("15fbbC47a244aE2A38071A106dCfcF3D57C9D939")
	if err != nil {
		t.Fatal(err)
	}

	secwareResult := types.SecwareResult{
		Code:           0,
		Message:        "Message",
		Details:        "Details",
		Operator:       operatorAddr,
		SecwareId:      1,
		SecwareVersion: 1,
	}

	secwareKey, err := hex.DecodeString("1111111111111111111111111111111111111111111111111111111111111111")
	if err != nil {
		t.Fatal(err)
	}

	signedSecwareResult, err := SignSecwareResult(&secwareResult, secwareKey)
	if err != nil {
		t.Fatal(err)
	}

	skBLS, err := hex.DecodeString("2222222222222222222222222222222222222222222222222222222222222222")
	if err != nil {
		t.Fatal(err)
	}

	blsKeyPair := bls.NewKeyPair(new(fr.Element).SetBytes(skBLS))
	signedOperatorResult, err := SignBLSOperatorResult(signedSecwareResult, blsKeyPair)
	if err != nil {
		t.Fatal(err)
	}

	b := VerifyBLSOperatorSignatureWithBLSPubKey(signedOperatorResult, blsKeyPair.GetPubKeyG2())
	if b != true {
		t.Fatal("bad BLS signature")
	}
}
