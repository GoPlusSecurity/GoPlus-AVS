// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import "@eigenlayer/contracts/permissions/PauserRegistry.sol";
import {DelegationManager} from "@eigenlayer/contracts/core/DelegationManager.sol";
import {AVSDirectory} from "@eigenlayer/contracts/core/AVSDirectory.sol";
import {StrategyManager} from "@eigenlayer/contracts/core/StrategyManager.sol";
import {IStrategyManager, IStrategy} from "@eigenlayer/contracts/interfaces/IStrategyManager.sol";
import {ISlasher} from "@eigenlayer/contracts/interfaces/ISlasher.sol";
import {StrategyBaseTVLLimits} from "@eigenlayer/contracts/strategies/StrategyBaseTVLLimits.sol";
import "@eigenlayer/test/mocks/EmptyContract.sol";

import "@eigenlayer-middleware/src/RegistryCoordinator.sol" as regcoord;
import {BLSApkRegistry} from "@eigenlayer-middleware/src/BLSApkRegistry.sol";
import {IndexRegistry} from "@eigenlayer-middleware/src/IndexRegistry.sol";
import {StakeRegistry} from "@eigenlayer-middleware/src/StakeRegistry.sol";
import "@eigenlayer-middleware/src/OperatorStateRetriever.sol";

import {GoPlusServiceManager, IServiceManager} from "../src/GoPlusServiceManager.sol";
import "../src/ERC20Mock.sol";

import {Utils} from "./utils/Utils.sol";

import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "forge-std/console2.sol";

// # To deploy and verify our contract
// forge script script/GoPlusDeployer.s.sol:HelloWorldDeployer --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract HoleskyGoPlusChecker is Script, Utils, Test {
    regcoord.RegistryCoordinator public constant reg = regcoord.RegistryCoordinator(0x3C503C651e3BD82C7AD169411E674d8ea6ad07e6);
    GoPlusServiceManager public immutable serviceMgr;
    StakeRegistry public immutable stakeReg;
    BLSApkRegistry public immutable blsReg;
    IndexRegistry public immutable indexReg;

    AVSDirectory public immutable avsDirectory;
    DelegationManager public immutable delegationMgr;


    constructor() {
        serviceMgr = GoPlusServiceManager(address(reg.serviceManager()));
        stakeReg = StakeRegistry(address(reg.stakeRegistry()));
        blsReg = BLSApkRegistry(address(reg.blsApkRegistry()));
        indexReg = IndexRegistry(address(reg.indexRegistry()));
        avsDirectory = AVSDirectory(serviceMgr.avsDirectory());
        delegationMgr = DelegationManager(address(avsDirectory.delegation()));
    }

    function depositToken() public {
        StrategyManager mgr = StrategyManager(0xdfB5f6CE42aAA7830E94ECFCcAd411beF4d4D5b6);
        IStrategy realt = IStrategy(0xAD76D205564f955A9c18103C4422D1Cd94016899);
        IERC20 token = IERC20(realt.underlyingToken());

        vm.startBroadcast();
        token.approve(address(mgr), 1 ether);
        mgr.depositIntoStrategy(realt, token, 50);
        vm.stopBroadcast();
    }

    function depositEther() public {
        StrategyManager mgr = StrategyManager(0xdfB5f6CE42aAA7830E94ECFCcAd411beF4d4D5b6);
        IStrategy realt = IStrategy(0xAD76D205564f955A9c18103C4422D1Cd94016899);
        IERC20 token = IERC20(realt.underlyingToken());

        vm.startBroadcast();
        token.approve(address(mgr), 1 ether);
        mgr.depositIntoStrategy(realt, token, 50);
        vm.stopBroadcast();
    }

    function changeQuorum() public {
        // 删除 0 号
        uint256[] memory indicesToRemove = new uint256[](2);
        // 必须从后向前删除 strategy
        indicesToRemove[0] = 1;
        indicesToRemove[1] = 0;

        uint8 quorumIdx = 0;
        StakeRegistry.StrategyParams[] memory params = new StakeRegistry.StrategyParams[](3);
        params[0].strategy = IStrategy(0x43252609bff8a13dFe5e057097f2f45A24387a84);
        params[1].strategy = IStrategy(0xAD76D205564f955A9c18103C4422D1Cd94016899);
        params[2].strategy = IStrategy(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0);
        params[0].multiplier = 1 ether;
        params[1].multiplier = 1 ether;
        params[2].multiplier = 1 ether;

        vm.startBroadcast();
        stakeReg.removeStrategies(quorumIdx, indicesToRemove);
        stakeReg.addStrategies(quorumIdx, params);
        stakeReg.setMinimumStakeForQuorum(quorumIdx, 1);

        address[] memory opList = new address[](7);
        opList[0] = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;
        opList[1] = 0x2D6B10600Ddd0B96bdd34346d1ab07236d01BE6E;
        opList[2] = 0x694B1da0b159289B7218D2AEc7FdbA98102436f7;
        opList[3] = 0x6e9df115BDdc029F31F2ed6901e405ABcAbf5E94;
        opList[4] = 0x5dE2805968a2cB2318Fe77fC44C39722b74118f6;
        opList[5] = 0xbE2d195D57217941fAb5bC8B554ad60899e99a0F;
        opList[6] = 0xD6e418e4E3c7a290750c0B9F60cea3cb0D635929;
        reg.updateOperators(opList);
        vm.stopBroadcast();
        // run();
    }

    function queueWithdrawals() public {
        address op = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;

        IStrategy[] memory strategies = new IStrategy[](1);
        strategies[0] = IStrategy(0x31B6F59e1627cEfC9fA174aD03859fC337666af7);
        // strategies[1] = IStrategy(0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3);

        uint256[] memory shares = new uint256[](1);
        shares[0] = strategies[0].shares(op);
        // shares[1] = strategies[1].shares(op);
        console.log("shares: %s", shares[0]);

        DelegationManager.QueuedWithdrawalParams[] memory params = new DelegationManager.QueuedWithdrawalParams[](1);
        params[0].strategies = strategies;
        params[0].shares = shares;
        params[0].withdrawer = op;

        DelegationManager d = DelegationManager(0xA44151489861Fe9e3055d95adC98FbD462B948e7);

        vm.broadcast();
        d.queueWithdrawals(params);
    }

    function completeQueuedWithdrawal() public {
        address op = 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939;
        IStrategy s = IStrategy(0x31B6F59e1627cEfC9fA174aD03859fC337666af7);
        DelegationManager d = DelegationManager(0xA44151489861Fe9e3055d95adC98FbD462B948e7);
        IStrategy[] memory strategies = new IStrategy[](1);
        strategies[0] = IStrategy(0x31B6F59e1627cEfC9fA174aD03859fC337666af7);
        uint256[] memory shares = new uint256[](1);
        shares[0] = 3990000000000000000;

        uint256 nonce = d.stakerNonce(op);
        console.log("staker nonce: %s", nonce);

        IERC20[] memory tokens = new IERC20[](1);
        tokens[0] = s.underlyingToken();
        DelegationManager.Withdrawal memory w;
        w.staker = op;
        w.delegatedTo = op;
        w.withdrawer = op;
        w.nonce = nonce;
        w.startBlock = 2519101;
        w.strategies = strategies;
        w.shares = shares;

        vm.broadcast();
        d.completeQueuedWithdrawal(w, tokens, 0, true);
    }

    function run() view public {
        console.log("==== chain id: %s ====", block.chainid);
        console.log("AVSDirectory: %s", address(avsDirectory));
        console.log("DelegationManager: %s", address(delegationMgr));

        console.log("");
        console.log("========");
        console.log("RegistryCoordinator: %s", address(reg));
        console.log("GoPlusServiceManager: %s", address(serviceMgr));
        console.log("StakeRegistry: %s", address(stakeReg));
        console.log("BLSApkRegistry: %s", address(blsReg));
        console.log("IndexRegistry: %s", address(indexReg));
        console.log("RegistryCoordinator's owner: %s", address(reg.owner()));
        console.log("churnApprover: %s", reg.churnApprover());
        console.log("ejector: %s", reg.ejector());

        console.log("");
        console.log("========");

        uint8 quorumCount = reg.quorumCount();
        {
            console.log("quorumCount: %s", quorumCount);
            for (uint8 quorumIdx = 0; quorumIdx < reg.quorumCount(); quorumIdx++) {
                regcoord.RegistryCoordinator.OperatorSetParam memory param = reg.getOperatorSetParams(quorumIdx);
                console.log("");
                console.log("quorum[%s]", quorumIdx);
                console.log("  maxOperatorCount: %s", param.maxOperatorCount);
                console.log("  kickBIPsOfOperatorStake: %s", param.kickBIPsOfOperatorStake);
                console.log("  kickBIPsOfTotalStake: %s", param.kickBIPsOfTotalStake);
                console.log("  minimumStakeForQuorum: %s", stakeReg.minimumStakeForQuorum(quorumIdx));
                uint256 len = stakeReg.strategyParamsLength(quorumIdx);
                console.log("  Num of strategies in quorum: %s", len);
                for (uint256 j = 0; j < len; j++) {
                    StakeRegistry.StrategyParams memory params = stakeReg.strategyParamsByIndex(quorumIdx, j);
                    console.log("    strategy[%s] address: %s, multiplier: %s", j, address(params.strategy), params.multiplier);
                }
                console.log("  CurrentTotalStake: %s", stakeReg.getCurrentTotalStake(quorumIdx));
            }
        }
        console.log("========");
        {
            showOperator(0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939);
            console.log("");
            showOperator(0x5dE2805968a2cB2318Fe77fC44C39722b74118f6);
            console.log("");
            showOperator(0xD6e418e4E3c7a290750c0B9F60cea3cb0D635929);
        }
//        console.log("  Total stake history:");
//        len = stakeReg.getTotalStakeHistoryLength(i);
//        for (uint256 j = 0; j < len; j++) {
//            StakeRegistry.StakeUpdate memory update = stakeReg.getTotalStakeUpdateAtIndex(i, j);
//            console.log("    Update block number: %s", update.updateBlockNumber);
//            console.log("    Stake: %s", update.stake);
//        }
    }

    function showOperator(address operator) view internal {
        regcoord.RegistryCoordinator.OperatorInfo memory opInfo = reg.getOperator(operator);
        uint8 status = uint8(opInfo.status);
        bytes32 operatorId = opInfo.operatorId;
        console.log("operator: %s", operator);
        if (status == 0) {
            console.log("status: NEVER_REGISTERED");
        } else if (status == 1) {
            console.log("status: REGISTERED");
        } else if (status == 2) {
            console.log("status: DEREGISTERED");
        } else {
            require(false, "bad operator status");
        }
        console.log("operatorId: %s", vm.toString(operatorId));

        uint192 bitmap = reg.getCurrentQuorumBitmap(operatorId);
        console.log("current quorum bitmap: %s", bitmap);
        for (uint8 quorumIdx = 0; quorumIdx < 192; quorumIdx++) {
            if (bitmap & (1 << quorumIdx) == 0) {
                continue;
            }
            console.log("weight in quorum[%s]: %s", quorumIdx, stakeReg.weightOfOperatorForQuorum(quorumIdx, operator));
        }


//        uint256 hisLen = reg.getQuorumBitmapHistoryLength(operatorId);
//        for (uint256 i = 0; i < hisLen; i++) {
//            regcoord.RegistryCoordinator.QuorumBitmapUpdate memory update = reg.getQuorumBitmapUpdateByIndex(operatorId, i);
//            console.log("updateBlockNumber[%s]: %s", i, update.updateBlockNumber);
//            console.log("quorumBitmap[%s]: %s", i, update.quorumBitmap);
//            console.log("");
//        }
    }
}
