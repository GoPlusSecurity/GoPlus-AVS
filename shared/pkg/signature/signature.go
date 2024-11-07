package signature

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"goplus/shared/pkg/types"
)

func HashJSON(obj interface{}) (common.Hash, error) {
	encoded, err := json.Marshal(obj)
	if err != nil {
		return common.Hash{}, err
	}
	return crypto.Keccak256Hash(encoded), nil
}

func SignHash(hash common.Hash, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	return crypto.Sign(hash.Bytes(), privateKey)
}

func GetPublicKeyBytes(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key")
	}

	// compressed form
	// return crypto.CompressPubkey(publicKeyECDSA), nil

	// uncompressed form
	return crypto.FromECDSAPub(publicKeyECDSA), nil
}

func GetAddressFromPrivateKey(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	// 获取公钥
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("error casting public key")
	}

	// 计算以太坊地址
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func VerifySignature(hash []byte, signature []byte, pkBytes []byte) bool {
	recoveredPkBytes, err := crypto.Ecrecover(hash, signature)
	if err != nil {
		return false
	}
	return bytes.Equal(recoveredPkBytes, pkBytes)
}

func VerifySignatureWithAddress(hash []byte, signature []byte, address common.Address) bool {
	recoveredPubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return false
	}

	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey)

	if bytes.Equal(address.Bytes(), recoveredAddr.Bytes()) {
		return true
	}
	return false
}

// SignSecwareTask 将 SecwareTask 转化为 SignedSecwareTask
// - skGateway: Gateway 的私钥
func SignSecwareTask(task *types.SecwareTask, skGateway *ecdsa.PrivateKey) (*types.SignedSecwareTask, error) {
	hashTask, err := HashJSON(task)
	if err != nil {
		return nil, err
	}
	sigGateway, err := SignHash(hashTask, skGateway)
	if err != nil {
		return nil, err
	}
	return &types.SignedSecwareTask{
		Task:       *task,
		SigGateway: sigGateway,
	}, nil
}

// VerifySignedSecwareTask 验证签名是否合法
func VerifySignedSecwareTask(task *types.SignedSecwareTask, pkGateway *ecdsa.PublicKey) bool {
	hashTask, err := HashJSON(task.Task)
	if err != nil {
		return false
	}
	return VerifySignature(hashTask.Bytes(), task.SigGateway, crypto.FromECDSAPub(pkGateway))
}

func VerifySignedSecwareTaskWithAddress(task *types.SignedSecwareTask, address common.Address) bool {
	hashTask, err := HashJSON(task.Task)
	if err != nil {
		return false
	}
	return VerifySignatureWithAddress(hashTask.Bytes(), task.SigGateway, address)
}

// ComputeHMAC256 生成 HMAC
// - result: SecwareResult.Details 的内容
// - operator: operator 的地址
// - key: Secware 的私钥
func ComputeHMAC256(msg, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(msg)
	return h.Sum(nil)
}

// SignSecwareResult 将 SecwareResult 转化为 SignedSecwareResult
// - key: Secware 的私钥
func SignSecwareResult(result *types.SecwareResult, key []byte) (*types.SignedSecwareResult, error) {
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	sigSecware := ComputeHMAC256(resultBytes, key)

	// 构造返回结果
	return &types.SignedSecwareResult{
		Result:     *result,
		SigSecware: sigSecware,
	}, nil
}

// VerifySecwareSignature 验证 Secware 的签名是否合法
// - key: Secware 的私钥
func VerifySecwareSignature(signedResult *types.SignedSecwareResult, key []byte) bool {
	expect, err := SignSecwareResult(&signedResult.Result, key)
	if err != nil {
		return false
	}
	return bytes.Equal(expect.SigSecware, signedResult.SigSecware)
}

func SignOperatorResult(result types.SignedSecwareResult, skOperator *ecdsa.PrivateKey) (*types.SignedOperatorResult, error) {
	hashTask, err := HashJSON(result)
	if err != nil {
		return nil, err
	}
	sigOperator, err := SignHash(hashTask, skOperator)
	if err != nil {
		return nil, err
	}

	return &types.SignedOperatorResult{
		Result:      result,
		SigOperator: sigOperator,
	}, nil
}

func VerifyOperatorSignature(signedResult *types.SignedOperatorResult, pkOperator *ecdsa.PublicKey) bool {
	hashTask, err := HashJSON(signedResult.Result)
	if err != nil {
		return false
	}
	return VerifySignature(hashTask.Bytes(), signedResult.SigOperator, crypto.FromECDSAPub(pkOperator))
}

func VerifyOperatorSignatureWithAddress(signedResult *types.SignedOperatorResult, address common.Address) bool {
	hashTask, err := HashJSON(signedResult.Result)
	if err != nil {
		return false
	}
	return VerifySignatureWithAddress(hashTask.Bytes(), signedResult.SigOperator, address)
}

func SignBLSOperatorResult(result *types.SignedSecwareResult, blsKeyPair *bls.KeyPair) (*types.SignedOperatorResult, error) {
	hashTask, err := HashJSON(result)
	if err != nil {
		return nil, err
	}
	sigOperator := blsKeyPair.SignMessage(hashTask)

	return &types.SignedOperatorResult{
		Result:      *result,
		SigOperator: sigOperator.Marshal(),
	}, nil
}

func VerifyBLSOperatorSignatureWithBLSPubKey(signedResult *types.SignedOperatorResult, blsPubKey *bls.G2Point) bool {
	hashTask, err := HashJSON(signedResult.Result)
	if err != nil {
		return false
	}

	sigBLS := bls.Signature{G1Point: new(bls.G1Point).Deserialize(signedResult.SigOperator)}
	b, err := sigBLS.Verify(blsPubKey, hashTask)
	if err != nil {
		return false
	}
	return b
}
