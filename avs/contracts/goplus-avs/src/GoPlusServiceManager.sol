// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import "@eigenlayer/contracts/interfaces/IAVSDirectory.sol";
import "@eigenlayer-middleware/src/ServiceManagerBase.sol";


contract GoPlusServiceManager is ServiceManagerBase {
    address public gatewayAddr;
    string public gatewayURI;

    event GatewayAddressUpdated(address indexed oldAddr, address indexed newAddr);
    event GatewayURIUpdated(string oldURI, string newURI);

    constructor(
        IAVSDirectory _avsDirectory,
        IRegistryCoordinator _registryCoordinator,
        IStakeRegistry _stakeRegistry
    )
        ServiceManagerBase(
            _avsDirectory,
            _registryCoordinator,
            _stakeRegistry
        )
    {
    }

    function initialize(address initialOwner, address _gatewayAddr, string memory _gatewayURI) public virtual initializer {
        __ServiceManagerBase_init(initialOwner);
        require(_gatewayAddr != address(0), "Gateway address cannot be zero");
        gatewayAddr = _gatewayAddr;
        gatewayURI = _gatewayURI;

        emit GatewayAddressUpdated(address(0), _gatewayAddr);
        emit GatewayURIUpdated("", _gatewayURI);
    }

    function updateGatewayAddress(address _gatewayAddr) external onlyOwner {
        require(_gatewayAddr != address(0), "Gateway address cannot be zero");
        address oldAddr = gatewayAddr;
        gatewayAddr = _gatewayAddr;
        emit GatewayAddressUpdated(oldAddr, _gatewayAddr);
    }

    function updateGatewayURI(string memory _gatewayURI) external onlyOwner {
        string memory oldURI = gatewayURI;
        gatewayURI = _gatewayURI;
        emit GatewayURIUpdated(oldURI, _gatewayURI);
    }
}
