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


    // Deploy GoPlus AVS contracts to Holesky testnet
    function runHolesky() external {
        // Define EL core contracts
        IAVSDirectory avsDirectory = IAVSDirectory(0x055733000064333CaDDbC92763c58BF0192fFeBf);
        IDelegationManager delegationManager = IDelegationManager(0xA44151489861Fe9e3055d95adC98FbD462B948e7);

        // Define GoPlus accounts
        // Replace them with multi-sig accounts
        address goPlusCommunityMultisig = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;  // Hardhat account #0
        address goPlusPauser = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;  // Hardhat account #1

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

        operatorStateRetriever = new OperatorStateRetriever();

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

            // setup quorum[0]
            // replace with proper settings
            regcoord.IRegistryCoordinator.OperatorSetParam[] memory quorumsOperatorSetParams = new regcoord.IRegistryCoordinator.OperatorSetParam[](1);
            quorumsOperatorSetParams[0] = regcoord.IRegistryCoordinator.OperatorSetParam({
                        maxOperatorCount: 10000,
                        kickBIPsOfOperatorStake: 15000,
                        kickBIPsOfTotalStake: 100
            });

            uint96[] memory quorumsMinimumStakes = new uint96[](1);
            quorumsMinimumStakes[0] = 1 ether;

            IStakeRegistry.StrategyParams[][] memory quorumsStrategyParams = new IStakeRegistry.StrategyParams[][](1);
            quorumsStrategyParams[0] = new IStakeRegistry.StrategyParams[](3);
            quorumsStrategyParams[0][0] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x43252609bff8a13dFe5e057097f2f45A24387a84),  // EIGEN
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][1] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xAD76D205564f955A9c18103C4422D1Cd94016899),  // reALT
                multiplier: 1 ether
            });

            quorumsStrategyParams[0][2] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0),  // Ether
                multiplier: 1 ether
            });

            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinator))),
                address(registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    regcoord.RegistryCoordinator.initialize.selector,
                    goPlusCommunityMultisig,  // _initialOwner
                    goPlusCommunityMultisig,  // _churnApprover
                    goPlusCommunityMultisig,  // _ejector
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

            address gatewayAddress = 0xfa2b8075362d1c6cd7c48306b68546482dac72c4;
            string memory gatewayURI = "https://avs.gopluslabs.io/api/v1";
            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(goPlusServiceManager))),
                address(goPlusServiceManagerImplementation),
                abi.encodeWithSelector(
                    GoPlusServiceManager.initialize.selector,
                    address(goPlusProxyAdmin),
                    gatewayAddress,
                    gatewayURI
                )
            );

            // setup AVSMetadataURI
            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(goPlusServiceManager))),
                address(goPlusServiceManagerImplementation),
                abi.encodeWithSelector(
                    IServiceManager.updateAVSMetadataURI.selector,
                    "https://static2.gopluslabs.io/avs_metadata.json"
                )
            );
        }
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
        require(address(goPlusProxyAdmin) == goPlusServiceManager.owner(), "The owner of ServiceManager is incorrect.");

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
            address(goPlusCommunityMultisig)
        );
        vm.serializeAddress(
            parent_object,
            "Ejector",
            address(goPlusCommunityMultisig)
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
    
    // Deploy GoPlus AVS contracts to mainnet
    function run() external {
        // Define EL core contracts
        IAVSDirectory avsDirectory = IAVSDirectory(0x135dda560e946695d6f155dacafc6f1f25c1f5af);
        IDelegationManager delegationManager = IDelegationManager(0x39053D51B77DC0d36036Fc1fCc8Cb819df8Ef37A);

        // Define GoPlus accounts
        // Replace them with multi-sig accounts
        address goPlusCommunityMultisig = 0x0A33f7Ad41A2Ed3510EF5a65b6B4397c6307e410;
        address goPlusPauser = 0xdf40044d40ff0d8f34d43f5cfc3d89c42bbfbde2;  
        address ejector = 0xbc6ce40a4137f42d14c8cd1aff944000c8921a1d;
        address churnApprover = 0xa6abe31f70311b59b2f1f0adc9cabd9bdab3dc55;

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

        operatorStateRetriever = new OperatorStateRetriever();

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
            quorumsMinimumStakes[0] = 1 ether;
            quorumsMinimumStakes[1] = 1 ether;
            quorumsMinimumStakes[2] = 1 ether;

            IStakeRegistry.StrategyParams[][] memory quorumsStrategyParams = new IStakeRegistry.StrategyParams[][](3);
            quorumsStrategyParams[0] = new IStakeRegistry.StrategyParams[](13);
            quorumsStrategyParams[1] = new IStakeRegistry.StrategyParams[](1);
            quorumsStrategyParams[2] = new IStakeRegistry.StrategyParams[](1);

            quorumsStrategyParams[0][0] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0),  // Beacon Chain ETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][1] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xa4C637e0F704745D182e4D38cAb7E7485321d059),  // oETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][2] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x9d7eD45EE2E8FC5482fa2428f15C971e6369011d),  // ETHx
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][3] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x298aFB19A105D59E74658C4C334Ff360BadE6dd2),  // mETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][4] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x8CA7A5d6f3acd3A7A8bC468a8CD0FB14B6BD28b6),  // sfrxETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][5] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xAe60d8180437b5C34bB956822ac2710972584473),  // lsETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][6] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x54945180dB7943c0ed0FEE7EdaB2Bd24620256bc),  // cbETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][7] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x13760F50a9d7377e4F20CB8CF9e4c26586c658ff),  // ankrETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][8] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x93c4b944D05dfe6df7645A86cd2206016c51564D),  // stETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][9] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x57ba429517c3473B6d34CA9aCd56c0e735b94c02),  // osETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][10] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x7CA911E83dabf90C90dD3De5411a10F1A6112184),  // wBETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][11] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x0Fe4F44beE93503346A3Ac9EE5A26b130a5796d6),  // swETH
                multiplier: 1 ether
            });
            quorumsStrategyParams[0][12] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0x1BeE69b7dFFfA4E2d53C2a2Df135C388AD25dCD2),  // rETH
                multiplier: 1 ether
            });
            
            quorumsStrategyParams[1][0] = IStakeRegistry.StrategyParams({
                strategy: IStrategy(0xaCB55C530Acdb2849e6d4f36992Cd8c9D50ED8F7),  // EigenStrategy (EIGEN)	
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

            address gatewayAddress = 0xfa2b8075362d1c6cd7c48306b68546482dac72c4;
            string memory gatewayURI = "https://avs.gopluslabs.io/api/v1";
            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(goPlusServiceManager))),
                address(goPlusServiceManagerImplementation),
                abi.encodeWithSelector(
                    GoPlusServiceManager.initialize.selector,
                    address(goPlusProxyAdmin),
                    gatewayAddress,
                    gatewayURI
                )
            );

            // setup AVSMetadataURI
            goPlusProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(goPlusServiceManager))),
                address(goPlusServiceManagerImplementation),
                abi.encodeWithSelector(
                    IServiceManager.updateAVSMetadataURI.selector,
                    "https://avs.gopluslabs.io/avs_metadata.json"
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
        require(address(goPlusProxyAdmin) == goPlusServiceManager.owner(), "The owner of ServiceManager is incorrect.");

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

