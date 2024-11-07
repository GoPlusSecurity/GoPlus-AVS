// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import "@eigenlayer/contracts/permissions/PauserRegistry.sol";
import {IDelegationManager} from "@eigenlayer/contracts/interfaces/IDelegationManager.sol";
import {IAVSDirectory} from "@eigenlayer/contracts/interfaces/IAVSDirectory.sol";
import {IStrategyManager, IStrategy} from "@eigenlayer/contracts/interfaces/IStrategyManager.sol";
import {ISlasher} from "@eigenlayer/contracts/interfaces/ISlasher.sol";
import {StrategyBaseTVLLimits} from "@eigenlayer/contracts/strategies/StrategyBaseTVLLimits.sol";
import "@eigenlayer/test/mocks/EmptyContract.sol";

import "@eigenlayer-middleware/src/RegistryCoordinator.sol" as regcoord;
import {IBLSApkRegistry, IIndexRegistry, IStakeRegistry} from "@eigenlayer-middleware/src/RegistryCoordinator.sol";
import {BLSApkRegistry} from "@eigenlayer-middleware/src/BLSApkRegistry.sol";
import {IndexRegistry} from "@eigenlayer-middleware/src/IndexRegistry.sol";
import {StakeRegistry} from "@eigenlayer-middleware/src/StakeRegistry.sol";
import "@eigenlayer-middleware/src/OperatorStateRetriever.sol";

import {GoPlusServiceManager, IServiceManager} from "../src/GoPlusServiceManager.sol";

import {Utils} from "./utils/Utils.sol";

import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "forge-std/console.sol";

// # To deploy and verify our contract
contract GoPlusDeployer is Script, Utils {
    // GoPlus contracts
    ProxyAdmin public goPlusProxyAdmin;
    PauserRegistry public goPlusPauserReg;

    regcoord.RegistryCoordinator public registryCoordinator;
    regcoord.IRegistryCoordinator public registryCoordinatorImplementation;

    IBLSApkRegistry public blsApkRegistry;
    IBLSApkRegistry public blsApkRegistryImplementation;

    IIndexRegistry public indexRegistry;
    IIndexRegistry public indexRegistryImplementation;

    IStakeRegistry public stakeRegistry;
    IStakeRegistry public stakeRegistryImplementation;

    OperatorStateRetriever public operatorStateRetriever;

    GoPlusServiceManager public goPlusServiceManager;
    IServiceManager public goPlusServiceManagerImplementation;

    // Deploy GoPlus AVS contracts to mainnet
    function run() external {
        // Define EL core contracts
        IAVSDirectory avsDirectory = IAVSDirectory(0x055733000064333CaDDbC92763c58BF0192fFeBf);
        IDelegationManager delegationManager = IDelegationManager(0xA44151489861Fe9e3055d95adC98FbD462B948e7);

        // Define GoPlus accounts
        // Replace them with multi-sig accounts
        address goPlusCommunityMultisig = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;
        address goPlusPauser = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;
        address ejector = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;
        address churnApprover = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;

        vm.startBroadcast();

        // deploy proxy admin for ability to upgrade proxy contracts
        {
            goPlusProxyAdmin = new ProxyAdmin();
        }

        // deploy pauser registry
        {
            address[] memory pausers = new address[](2);
            pausers[0] = goPlusPauser;
            pausers[1] = goPlusCommunityMultisig;
            goPlusPauserReg = new PauserRegistry(pausers, goPlusCommunityMultisig);
        }

        // deploy empty contract
        EmptyContract emptyContract = new EmptyContract();

        // deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
        // not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
        goPlusServiceManager = GoPlusServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(goPlusProxyAdmin),
                    ""
                )
            )
        );

        registryCoordinator = regcoord.RegistryCoordinator(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(goPlusProxyAdmin),
                    ""
                )
            )
        );
        blsApkRegistry = IBLSApkRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(goPlusProxyAdmin),
                    ""
                )
            )
        );
        indexRegistry = IIndexRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(goPlusProxyAdmin),
                    ""
                )
            )
        );
        stakeRegistry = IStakeRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(goPlusProxyAdmin),
                    ""
                )
            )
        );

        operatorStateRetriever = OperatorStateRetriever(0x5ce26317F7edCBCBD1a569629af5DC41c1622045);

        // Second, deploy the *implementation* contracts, using the *proxy contracts* as inputs
        {
            // setup StakeRegistry impl
            stakeRegistryImplementation = new StakeRegistry(
                registryCoordinator,
                delegationManager
            );

            goPlusProxyAdmin.upgrade(
                TransparentUpgradeableProxy(payable(address(stakeRegistry))),
                address(stakeRegistryImplementation)
            );
        }

        {
            // setup BLSApkRegistry impl
            blsApkRegistryImplementation = new BLSApkRegistry(
                registryCoordinator
            );

            goPlusProxyAdmin.upgrade(
                TransparentUpgradeableProxy(payable(address(blsApkRegistry))),
                address(blsApkRegistryImplementation)
            );
        }

        {
            // setup IndexRegistry impl
            indexRegistryImplementation = new IndexRegistry(
                registryCoordinator
            );

            goPlusProxyAdmin.upgrade(
                TransparentUpgradeableProxy(payable(address(indexRegistry))),
                address(indexRegistryImplementation)
            );
        }

        {
            // setup RegistryCoordinator impl
            registryCoordinatorImplementation = new regcoord.RegistryCoordinator(
                goPlusServiceManager,
                regcoord.IStakeRegistry(address(stakeRegistry)),
                regcoord.IBLSApkRegistry(address(blsApkRegistry)),
                regcoord.IIndexRegistry(address(indexRegistry))
            );

            regcoord.IRegistryCoordinator.OperatorSetParam[] memory quorumsOperatorSetParams = new regcoord.IRegistryCoordinator.OperatorSetParam[](3);
            quorumsOperatorSetParams[0] = regcoord.IRegistryCoordinator.OperatorSetParam({
                maxOperatorCount: 10000,
                kickBIPsOfOperatorStake: 15000,
                kickBIPsOfTotalStake: 100
            });
            quorumsOperatorSetParams[1] = regcoord.IRegistryCoordinator.OperatorSetParam({
                maxOperatorCount: 10000,
                kickBIPsOfOperatorStake: 15000,
                kickBIPsOfTotalStake: 100
            });
            quorumsOperatorSetParams[2] = regcoord.IRegistryCoordinator.OperatorSetParam({
                maxOperatorCount: 10000,
                kickBIPsOfOperatorStake: 15000,
                kickBIPsOfTotalStake: 100
            });

            uint96[] memory quorumsMinimumStakes = new uint96[](3);
            quorumsMinimumStakes[0] = 1;
            quorumsMinimumStakes[1] = 1;
            quorumsMinimumStakes[2] = 1;

            IStakeRegistry.StrategyParams[][] memory quorumsStrategyParams = new IStakeRegistry.StrategyParams[][](3);
            quorumsStrategyParams[0] = new IStakeRegistry.StrategyParams[](11);
            quorumsStrategyParams[1] = new IStakeRegistry.StrategyParams[](1);
            quorumsStrategyParams[2] = new IStakeRegistry.StrategyParams[](1);

            quorumsStrategyParams[0][0] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0),  // Beacon Chain ETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][1] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3),  // stETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][2] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x3A8fBdf9e77DFc25d09741f51d3E181b25d0c4E0),  // rETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][3] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x80528D6e9A2BAbFc766965E0E26d5aB08D9CFaF9),  // WETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][4] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x05037A81BD7B4C9E0F7B430f1F2A22c31a2FD943),  // lsETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][5] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x9281ff96637710Cd9A5CAcce9c6FAD8C9F54631c),  // sfrxETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][6] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x31B6F59e1627cEfC9fA174aD03859fC337666af7),  // ETHx
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][7] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x46281E3B7fDcACdBa44CADf069a94a588Fd4C6Ef),  // osETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][8] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x70EB4D3c164a6B4A5f908D4FBb5a9cAfFb66bAB6	),  // cbETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][9] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xaccc5A86732BE85b5012e8614AF237801636F8e5),  // mETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][10] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x7673a47463F80c6a3553Db9E54c8cDcd5313d0ac),  // ankrETH
                multiplier: 1 ether
            });

            quorumsStrategyParams[1][0] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x43252609bff8a13dFe5e057097f2f45A24387a84),  // EigenStrategy (EIGEN)
                multiplier: 1 ether
            });

            quorumsStrategyParams[2][0] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xAD76D205564f955A9c18103C4422D1Cd94016899),  // reALT
                multiplier: 1 ether
            });

            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinator))),
                address(registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    regcoord.RegistryCoordinator.initialize.selector,
                    goPlusCommunityMultisig,  // _initialOwner
                    churnApprover,  // _churnApprover
                    ejector,  // _ejector
                    goPlusPauserReg,  // _pauserRegistry
                    0,  // _initialPausedStatus
                    quorumsOperatorSetParams,
                    quorumsMinimumStakes,
                    quorumsStrategyParams
                )
            );
        }

        {
            // setup ServiceManager impl
            goPlusServiceManagerImplementation = new GoPlusServiceManager(
                avsDirectory,
                registryCoordinator,
                stakeRegistry
            );

            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(goPlusServiceManager))),
                address(goPlusServiceManagerImplementation),
                abi.encodeWithSelector(
                    GoPlusServiceManager.initialize.selector,
                    address(goPlusCommunityMultisig),
                    0x96216849c49358B10257cb55b28eA603c874b05E , // gatewayAddress
                    "https://test-go-avs.ansuzsecurity.com/api/v1", // gatewayURI
                    "https://static2.gopluslabs.io/avs_metadata.json" // metadataURI
                )
            );
        }
        goPlusProxyAdmin.transferOwnership(goPlusCommunityMultisig);
        vm.stopBroadcast();

        // write deployment results to logs
        console.log("Deployer: %s", msg.sender);
        console.log("GoPlusProxyAdmin: %s", address(goPlusProxyAdmin));
        console.log("GoPlusProxyAdmin owner: %s", goPlusProxyAdmin.owner());
        console.log("GoPlusServiceManager: %s, Impl: %s", address(goPlusServiceManager), address(goPlusServiceManagerImplementation));
        console.log("RegistryCoordinator: %s, Impl: %s", address(registryCoordinator), address(registryCoordinatorImplementation));
        console.log("BLSApkRegistry: %s, Impl: %s", address(blsApkRegistry), address(blsApkRegistryImplementation));
        console.log("IndexRegistry: %s, Impl: %s", address(indexRegistry), address(indexRegistryImplementation));
        console.log("StakeRegistry: %s, Impl: %s", address(stakeRegistry), address(stakeRegistryImplementation));
        console.log("OperatorStateRetriever: %s", address(operatorStateRetriever));
        console.log("GoPlusServiceManager owner: %s", goPlusServiceManager.owner());

        // write deployment results to JSON
        string memory parent_object = "addresses";
        vm.serializeAddress(
            parent_object,
            "deployer",
            address(msg.sender)
        );
        vm.serializeAddress(
            parent_object,
            "goPlusProxyAdmin",
            address(goPlusProxyAdmin)
        );
        vm.serializeAddress(
            parent_object,
            "goPlusProxyAdminOwner",
            address(goPlusProxyAdmin.owner())
        );
        vm.serializeAddress(
            parent_object,
            "goPlusServiceManager",
            address(goPlusServiceManager)
        );
        vm.serializeAddress(
            parent_object,
            "goPlusServiceManagerOwner",
            goPlusServiceManager.owner()
        );
        vm.serializeAddress(
            parent_object,
            "goPlusServiceManagerImplementation",
            address(goPlusServiceManagerImplementation)
        );
        vm.serializeAddress(
            parent_object,
            "goPlusCommunityMultisig",
            address(goPlusCommunityMultisig)
        );
        vm.serializeAddress(
            parent_object,
            "goPlusPauser",
            address(goPlusPauser)
        );
        vm.serializeAddress(
            parent_object,
            "goPlusPauserRegistry",
            address(goPlusPauserReg)
        );
        vm.serializeAddress(
            parent_object,
            "churnApprover",
            address(churnApprover)
        );
        vm.serializeAddress(
            parent_object,
            "Ejector",
            address(ejector)
        );
        vm.serializeAddress(
            parent_object,
            "registryCoordinator",
            address(registryCoordinator)
        );
        vm.serializeAddress(
            parent_object,
            "registryCoordinatorImplementation",
            address(registryCoordinatorImplementation)
        );
        vm.serializeAddress(
            parent_object,
            "blsApkRegistry",
            address(blsApkRegistry)
        );
        vm.serializeAddress(
            parent_object,
            "blsApkRegistryImplementation",
            address(blsApkRegistryImplementation)
        );
        vm.serializeAddress(
            parent_object,
            "indexRegistry",
            address(indexRegistry)
        );
        vm.serializeAddress(
            parent_object,
            "indexRegistryImplementation",
            address(indexRegistryImplementation)
        );
        vm.serializeAddress(
            parent_object,
            "stakeRegistry",
            address(stakeRegistry)
        );
        vm.serializeAddress(
            parent_object,
            "stakeRegistryImplementation",
            address(stakeRegistryImplementation)
        );
        string memory finalJson = vm.serializeAddress(
            parent_object,
            "operatorStateRetriever",
            address(operatorStateRetriever)
        );
        writeOutput(finalJson, "goplus_avs_deployment_output");
    }
}

