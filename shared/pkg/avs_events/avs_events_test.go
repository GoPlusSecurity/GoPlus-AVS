package avs_events

import (
	"context"
	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	svrmgr "github.com/Layr-Labs/eigensdk-go/contracts/bindings/ServiceManagerBase"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"testing"
)

func TestEventScanner(t *testing.T) {
	client, err := ethclient.Dial("http://localhost:8545") // 假设您的 Anvil 节点在本地运行
	if err != nil {
		log.Fatal(err)
	}

	blockNum, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	contracts := make(map[common.Address]ContractInfo)

	serviceManagerABI, err := svrmgr.ContractServiceManagerBaseMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}
	contracts[common.HexToAddress("0xa82fF9aFd8f496c3d6ac40E2a0F282E47488CFc9")] = ContractInfo{
		Name: "ServiceManager",
		ABI:  serviceManagerABI,
	}

	registryCoordinatorABI, err := regcoord.ContractRegistryCoordinatorMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}
	contracts[common.HexToAddress("0x1613beB3B2C4f22Ee086B2b38C1476A3cE7f78E8")] = ContractInfo{
		Name: "RegistryCoordinator",
		ABI:  registryCoordinatorABI,
	}

	if _, err := ScanEvents(client, big.NewInt(0), big.NewInt(int64(blockNum)), contracts); err != nil {
		log.Fatal(err)
	}
}
