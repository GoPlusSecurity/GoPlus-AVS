// 启动 AVS 的服务组
package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/json"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	sdkecdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	eigenSdkTypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/urfave/cli/v2"
	"goplus/avs/config"
	"goplus/avs/metrics"
	"goplus/avs/secwaremanager"
	"goplus/avs/server"
	"goplus/avs/state"
	"goplus/shared/pkg/signature"
	"goplus/shared/pkg/types"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func getOperatorECDSAKey() (*ecdsa.PrivateKey, error) {
	ecdsaKeyStorePath, _ := config.GetOperatorECDSAKeyStorePath()
	if ecdsaKeyStorePath == "" {
		log.Fatal("Operator ECDSA Key Store Path not provided")
	}
	ecdsaKeyPassword, _ := config.GetOperatorECDSAKeyPassword()

	skOperator, err := sdkecdsa.ReadKey(ecdsaKeyStorePath, ecdsaKeyPassword)
	if err != nil {
		log.Fatalf("Failed to read operator ECDSA key: %v", err)
		return nil, err
	}
	return skOperator, nil
}

func registerWithAVS(cliCtx *cli.Context) error {
	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		log.Fatal(err)
	}

	skOperator, err := getOperatorECDSAKey()
	if err != nil {
		cfg.Logger.Fatalf("Failed to read operator ECDSA key: %v", err)
		return err
	}

	opAddr, err := signature.GetAddressFromPrivateKey(skOperator)
	if err != nil {
		cfg.Logger.Fatal("Failed to calculate operator address.")
		return err
	}

	chainId, err := cfg.EthHttpClient.ChainID(cliCtx.Context)
	if err != nil {
		cfg.Logger.Fatal("Failed to get ChainID.")
		return err
	}

	signerV2 := func(ctx context.Context, address common.Address) (bind.SignerFn, error) {
		return signerv2.PrivateKeySignerFn(skOperator, chainId)
	}

	skWallet, err := wallet.NewPrivateKeyWallet(cfg.EthHttpClient, signerV2, opAddr, cfg.Logger)
	if err != nil {
		panic(err)
	}

	txMgr := txmgr.NewSimpleTxManager(skWallet, cfg.EthHttpClient, cfg.Logger, opAddr)

	var quorumNumbers eigenSdkTypes.QuorumNums
	for _, q := range cfg.QuorumNums {
		quorumNumbers = append(quorumNumbers, eigenSdkTypes.QuorumNum(q))
	}

	if len(quorumNumbers) != 1 || quorumNumbers[0] != 0 {
		panic("Only quorum[0] is available now.")
	}

	socket := types.OperatorSocket{
		NodeClass: cfg.NodeClass,
		URL:       cfg.OperatorURL,
	}
	socketBytes, err := json.Marshal(socket)
	if err != nil {
		cfg.Logger.Fatalf("Failed to json dump socket: %#v", socket)
		return err
	}

	var operatorToAvsRegistrationSigSalt [32]byte
	_, err = rand.Read(operatorToAvsRegistrationSigSalt[:])
	if err != nil {
		log.Fatalf("Failed to generate random salt: %v", err)
	}

	curBlockNum, err := cfg.EthHttpClient.BlockNumber(context.Background())
	if err != nil {
		cfg.Logger.Errorf("Unable to get current block number")
		return err
	}
	curBlock, err := cfg.EthHttpClient.HeaderByNumber(context.Background(), big.NewInt(int64(curBlockNum)))
	if err != nil {
		cfg.Logger.Errorf("Unable to get current block")
		return err
	}
	sigValidForSeconds := int64(1_000_000)
	operatorToAvsRegistrationSigExpiry := big.NewInt(int64(curBlock.Time) + sigValidForSeconds)

	avsRegistryWriter, err := avsregistry.NewWriterFromConfig(avsregistry.Config{
		RegistryCoordinatorAddress:    cfg.RegCoordinatorAddr,
		OperatorStateRetrieverAddress: cfg.OperatorStateRetrieverAddr,
	}, cfg.EthHttpClient, txMgr, cfg.Logger)
	if err != nil {
		cfg.Logger.Fatal("Failed to crete avsRegistryWriter")
		return err
	}

	_, err = avsRegistryWriter.RegisterOperatorInQuorumWithAVSRegistryCoordinator(
		context.Background(),
		skOperator,
		operatorToAvsRegistrationSigSalt,
		operatorToAvsRegistrationSigExpiry,
		cfg.BLSKeypair,
		quorumNumbers,
		string(socketBytes),
	)
	if err != nil {
		cfg.Logger.Errorf("Unable to register operator with avs registry coordinator")
		return err
	}
	cfg.Logger.Infof("Registered operator with avs registry coordinator.")

	return nil
}

func deregisterWithAVS(cliCtx *cli.Context) error {
	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		log.Fatal(err)
	}

	skOperator, err := getOperatorECDSAKey()
	if err != nil {
		cfg.Logger.Fatalf("Failed to read operator ECDSA key: %v", err)
		return err
	}

	opAddr, err := signature.GetAddressFromPrivateKey(skOperator)
	if err != nil {
		cfg.Logger.Fatal("Failed to calculate operator address.")
		return err
	}
	chainId, err := cfg.EthHttpClient.ChainID(cliCtx.Context)
	if err != nil {
		cfg.Logger.Fatal("Failed to get ChainID.")
		return err
	}

	signerV2 := func(ctx context.Context, address common.Address) (bind.SignerFn, error) {
		return signerv2.PrivateKeySignerFn(skOperator, chainId)
	}

	skWallet, err := wallet.NewPrivateKeyWallet(cfg.EthHttpClient, signerV2, opAddr, cfg.Logger)
	if err != nil {
		panic(err)
	}

	txMgr := txmgr.NewSimpleTxManager(skWallet, cfg.EthHttpClient, cfg.Logger, opAddr)

	avsRegistryWriter, err := avsregistry.NewWriterFromConfig(avsregistry.Config{
		RegistryCoordinatorAddress:    cfg.RegCoordinatorAddr,
		OperatorStateRetrieverAddress: cfg.OperatorStateRetrieverAddr,
	}, cfg.EthHttpClient, txMgr, cfg.Logger)
	if err != nil {
		cfg.Logger.Fatal("Failed to crete avsRegistryWriter")
		return err
	}

	var quorumNumbers eigenSdkTypes.QuorumNums
	for _, q := range cfg.QuorumNums {
		quorumNumbers = append(quorumNumbers, eigenSdkTypes.QuorumNum(q))
	}

	g1Point := cfg.BLSKeypair.GetPubKeyG1()
	bn254G1Point := regcoord.BN254G1Point{
		X: g1Point.X.BigInt(big.NewInt(0)),
		Y: g1Point.Y.BigInt(big.NewInt(0)),
	}

	_, err = avsRegistryWriter.DeregisterOperator(
		context.Background(),
		quorumNumbers,
		bn254G1Point,
	)

	if err != nil {
		cfg.Logger.Errorf("Unable to deregister operator with avs registry coordinator")
		return err
	}
	cfg.Logger.Infof("Deregistered operator with avs registry coordinator.")

	return nil
}

func start(cliCtx *cli.Context) error {
	cfg, err := config.NewConfig(cliCtx)
	if err != nil {
		log.Fatal(err)
	}

	mt, err := metrics.NewAvsMetrics(cfg)
	if err != nil {
		log.Fatal(err)
	}

	st, err := state.NewAvsDbState(cfg)
	if err != nil {
		log.Fatal(err)
	}

	manager, err := secwaremanager.NewSecwareManager(cfg, mt)
	if err != nil {
		log.Fatal(err)
	}

	err = manager.Init()
	if err != nil {
		log.Fatal(err)
	}

	svr, err := server.New(cfg, mt, manager, st)
	if err != nil {
		log.Fatal(err)
	}

	err = svr.Init()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(cliCtx.Context)
	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := manager.Start(ctx); err != nil {
			cancel()
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		if err := svr.Start(ctx); err != nil {
			cancel()
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		cfg.Logger.Info("Received signal, shutting down")
	case <-ctx.Done():
		cfg.Logger.Info("Context cancelled, shutting down")
	}

	cancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer timeoutCancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-timeoutCtx.Done():
	case <-done:
	}

	return nil
}
