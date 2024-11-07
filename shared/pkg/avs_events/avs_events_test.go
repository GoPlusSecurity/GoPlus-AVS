package avs_events

import (
	"context"
	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	svrmgr "goplus/avs/contracts/bindings/GoPlusServiceManager"
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

	serviceManagerABI, err := svrmgr.GoPlusServiceManagerMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}
	contracts[common.HexToAddress("0x6E0e0479e177c7F5111682C7025b4412613cd9dE")] = ContractInfo{
		Name: "ServiceManager",
		ABI:  serviceManagerABI,
	}

	registryCoordinatorABI, err := regcoord.ContractRegistryCoordinatorMetaData.GetAbi()
	if err != nil {
		log.Fatal(err)
	}
	contracts[common.HexToAddress("0x61AA80e5891DbfCebD0B78a704F3de996E449FdE")] = ContractInfo{
		Name: "RegistryCoordinator",
		ABI:  registryCoordinatorABI,
	}

	if _, err := ScanEvents(client, big.NewInt(2688942), big.NewInt(int64(blockNum)), contracts); err != nil {
		log.Fatal(err)
	}
}
