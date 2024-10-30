package chainio

import (
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	contractRegistryCoordinator2 "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"goplus/avs/contracts/bindings/GoPlusServiceManager"
)

type AvsReader struct {
	ethClient      eth.Client
	ServiceManager *contractGoPlusServiceManager.GoPlusServiceManager
}

func NewAvsReader(registryCoordinatorAddr common.Address, ethClient eth.Client) (*AvsReader, error) {
	contractRegistryCoordinator, err := contractRegistryCoordinator2.NewContractRegistryCoordinator(registryCoordinatorAddr, ethClient)
	if err != nil {
		return nil, err
	}

	serverManagerAddr, err := contractRegistryCoordinator.ServiceManager(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	serverManager, err := contractGoPlusServiceManager.NewGoPlusServiceManager(serverManagerAddr, ethClient)
	if err != nil {
		return nil, err
	}

	return &AvsReader{
		ethClient:      ethClient,
		ServiceManager: serverManager,
	}, nil
}

func (r *AvsReader) GatewayConfig() (common.Address, string, error) {
	addr, err := r.ServiceManager.GatewayAddr(&bind.CallOpts{})
	if err != nil {
		return common.Address{}, "", err
	}

	url, err := r.ServiceManager.GatewayURI(&bind.CallOpts{})
	if err != nil {
		return common.Address{}, "", err
	}

	return addr, url, nil
}
