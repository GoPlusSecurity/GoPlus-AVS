package config

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"goplus/avs/chainio"
	"goplus/shared/pkg/types"
	"path/filepath"
	"strings"
)

var (
	ConfigFileFlag = "config-file"
)

type RawConfig struct {
	ComposeFilePath            string `mapstructure:"COMPOSE_FILE_PATH"`
	OperatorSecretKey          string `mapstructure:"OPERATOR_SECRET_KEY"`
	OperatorBLSSecretKey       string `mapstructure:"OPERATOR_BLS_SECRET_KEY"`
	APIPort                    int    `mapstructure:"API_PORT"`
	OperatorURL                string `mapstructure:"OPERATOR_URL"`
	NodeClass                  string `mapstructure:"NODE_CLASS"`
	ETHRpc                     string `mapstructure:"ETH_RPC"`
	RegCoordinatorAddr         string `mapstructure:"REGISTRY_COORDINATOR_ADDR"`
	OperatorStateRetrieverAddr string `mapstructure:"OPERATOR_STATE_RETRIEVER"`
	QuorumNums                 []int  `mapstructure:"QUORUM_NUMS"`
}

func (r *RawConfig) isValid() error {
	if r.ComposeFilePath == "" {
		return fmt.Errorf("compose file path is required")
	}
	if r.OperatorSecretKey == "" {
		return fmt.Errorf("operator secret key is required")
	}
	if r.OperatorBLSSecretKey == "" {
		return fmt.Errorf("operator bls secret key is required")
	}
	if r.APIPort == 0 {
		return fmt.Errorf("api port is required")
	}
	if r.OperatorURL == "" {
		return fmt.Errorf("api host is required")
	}
	if !strings.HasPrefix(r.OperatorURL, "http://") && !strings.HasPrefix(r.OperatorURL, "https://") {
		return fmt.Errorf("api host must start with http:// or https://")
	}
	if r.NodeClass == "" {
		return fmt.Errorf("node class is required")
	}
	if r.NodeClass != "s" && r.NodeClass != "m" && r.NodeClass != "l" && r.NodeClass != "xl" {
		return fmt.Errorf("node class must be one of s, m, l, xl")
	}
	if r.ETHRpc == "" {
		return fmt.Errorf("eth rpc is required")
	}
	if r.RegCoordinatorAddr == "" {
		return fmt.Errorf("registry coordinator address is required")
	}
	if r.OperatorStateRetrieverAddr == "" {
		return fmt.Errorf("operator state retriever address is required")
	}
	if len(r.QuorumNums) == 0 {
		return fmt.Errorf("quorum nums is required")
	}

	return nil
}

type Config struct {
	Logger                     sdklogging.Logger
	ComposeFilePath            string
	GatewayUrl                 string
	RegCoordinatorAddr         common.Address
	OperatorStateRetrieverAddr common.Address

	ETHRpc            string
	EthHttpClient     eth.Client
	PromMetricsIpPort string

	BLSKeypair      *bls.KeyPair
	AddressGateway  common.Address
	AddressOperator common.Address
	SkOperator      *ecdsa.PrivateKey

	NodeClass   string
	OperatorURL string
	APIPort     int

	QuorumNums []int
}

func getRawConfigFromFile(filePath string) (RawConfig, error) {
	rawConfig := RawConfig{}
	if filePath != "" {
		viper.SetConfigFile(filePath)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return RawConfig{}, err
	}

	err = viper.Unmarshal(&rawConfig)
	if err != nil {
		return RawConfig{}, err
	}

	return rawConfig, nil
}

func getRawConfigFromEnv() (RawConfig, error) {
	rawConfig := RawConfig{}
	err := viper.BindEnv("COMPOSE_FILE_PATH")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("OPERATOR_SECRET_KEY")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("OPERATOR_BLS_SECRET_KEY")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("API_PORT")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("API_HOST")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("NODE_CLASS")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("ETH_RPC")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("REGISTRY_COORDINATOR_ADDR")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("OPERATOR_STATE_RETRIEVER")
	if err != nil {
		return RawConfig{}, err
	}
	err = viper.BindEnv("QUORUM_NUMS")
	if err != nil {
		return RawConfig{}, err
	}

	err = viper.Unmarshal(&rawConfig)
	if err != nil {
		return RawConfig{}, err
	}

	return rawConfig, nil
}

func NewConfig(ctx *cli.Context) (Config, error) {
	configFilePath := ctx.String(ConfigFileFlag)

	var rawConfig RawConfig
	var err error
	if configFilePath != "" {
		rawConfig, err = getRawConfigFromFile(configFilePath)
		if err != nil {
			return Config{}, err
		}
	} else {
		rawConfig, err = getRawConfigFromEnv()
		if err != nil {
			return Config{}, err
		}
	}

	if err := rawConfig.isValid(); err != nil {
		return Config{}, err
	}

	logger, err := sdklogging.NewZapLogger(sdklogging.Production)
	if err != nil {
		return Config{}, err
	}

	ethRpcClient, err := eth.NewClient(rawConfig.ETHRpc)
	if err != nil {
		return Config{}, err
	}

	regCoordinatorAddr := common.HexToAddress(rawConfig.RegCoordinatorAddr)
	avsReader, err := chainio.NewAvsReader(regCoordinatorAddr, ethRpcClient)
	if err != nil {
		return Config{}, err
	}

	gatewayAddr, gatewayUrl, err := avsReader.GatewayConfig()
	if err != nil {
		return Config{}, err
	}

	logger.Infof("Gateway address: %s", gatewayAddr.Hex())
	logger.Infof("Gateway url: %s", gatewayUrl)

	skOperatorString := rawConfig.OperatorSecretKey
	skOperatorString = strings.TrimPrefix(skOperatorString, "0x")
	if skOperatorString == "" {
		return Config{}, fmt.Errorf("operator secret key is required")
	}

	skOperator, err := crypto.HexToECDSA(skOperatorString)
	if err != nil {
		return Config{}, err
	}

	addressOperator := crypto.PubkeyToAddress(skOperator.PublicKey)
	addressOperatorString := addressOperator.Hex()
	logger.Infof("Operator address: %s", addressOperatorString)

	if !filepath.IsAbs(rawConfig.ComposeFilePath) {
		return Config{}, fmt.Errorf("compose file path must be absolute")
	}

	skBLSString := rawConfig.OperatorBLSSecretKey
	if skBLSString == "" {
		return Config{}, fmt.Errorf("operator bls secret key is required")
	}
	skBLS, err := types.NewHexBytesFromString(skBLSString)
	if err != nil {
		return Config{}, fmt.Errorf("cannot parse operator's BLS private key")
	}
	blsKeyPair := bls.NewKeyPair(new(fr.Element).SetBytes(skBLS))

	return Config{
		Logger:          logger,
		ComposeFilePath: rawConfig.ComposeFilePath,

		RegCoordinatorAddr:         common.HexToAddress(rawConfig.RegCoordinatorAddr),
		OperatorStateRetrieverAddr: common.HexToAddress(rawConfig.OperatorStateRetrieverAddr),

		ETHRpc:        rawConfig.ETHRpc,
		EthHttpClient: ethRpcClient,

		GatewayUrl:      gatewayUrl,
		AddressGateway:  gatewayAddr,
		AddressOperator: addressOperator,
		SkOperator:      skOperator,
		BLSKeypair:      blsKeyPair,

		NodeClass:   rawConfig.NodeClass,
		OperatorURL: rawConfig.OperatorURL,
		APIPort:     rawConfig.APIPort,

		QuorumNums: rawConfig.QuorumNums,
	}, nil
}
