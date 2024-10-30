// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import "@eigenlayer/contracts/interfaces/IAVSDirectory.sol";
import "@eigenlayer-middleware/src/ServiceManagerBase.sol";


contract GoPlusServiceManager is ServiceManagerBase {
    address public gatewayAddr;
    string public gatewayURI;

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
        gatewayAddr = _gatewayAddr;
        gatewayURI = _gatewayURI;
    }

    function updateGatewayAddress(address _gatewayAddr) external onlyOwner {
        gatewayAddr = _gatewayAddr;
    }

    function updateGatewayURI(string memory _gatewayURI) external onlyOwner {
        gatewayURI = _gatewayURI;
    }
}
