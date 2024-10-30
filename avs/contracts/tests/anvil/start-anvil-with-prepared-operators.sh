#!/bin/bash
RPC_URL=http://localhost:8545
PRIVATE_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# cd to the directory of this script so that this can be run from anywhere
parent_path=$(
    cd "$(dirname "${BASH_SOURCE[0]}")"
    pwd -P
)
cd "$parent_path"

set -a
source ./utils.sh
set +a

cleanup() {
    echo "Executing cleanup function..."
    set +e
    docker rm -f anvil
    exit_status=$?
    if [ $exit_status -ne 0 ]; then
        echo "Script exited due to set -e on line $1 with command '$2'. Exit status: $exit_status"
    fi
}
trap 'cleanup $LINENO "$BASH_COMMAND"' EXIT

# start an anvil instance in the background that has eigenlayer contracts deployed
# we start anvil in the background so that we can run the below script
# anvil --load-state avs-and-eigenlayer-deployed-anvil-state.json &
# FIXME: bug in latest foundry version, so we use this pinned version instead of latest
start_anvil_docker "$parent_path/avs-and-eigenlayer-deployed-anvil-state.json" ""

rm -f ~/.foundry/keystores/test-op-1 ~/.foundry/keystores/test-op-2 ~/.foundry/keystores/test-op-3
# Account #1: 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 (10000 ETH)
cast wallet import test-op-1 --private-key 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d --unsafe-password test
# Account #2: 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC (10000 ETH)
cast wallet import test-op-2 --private-key 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a --unsafe-password test
# Account #3: 0x90F79bf6EB2c4f870365E785982E1f101E93b906 (10000 ETH)
cast wallet import test-op-3 --private-key 0x7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6 --unsafe-password test

CHAIN_ID=$(cast chain-id)
DELEGATION_ADDR=$(jq -r '.addresses.delegation' ../../goplus-avs/script/output/"$CHAIN_ID"/eigenlayer_deployment_output.json)

for i in {1..3}; do
    # 构造账户名称
    account_name="test-op-${i}"

    # 获取 operator 地址
    OPERATOR=$(cast wallet address --account "$account_name" --password "test")

    echo "Register operator to EigenLayer... $OPERATOR"

    cast send \
      --json \
      --account $account_name --password test \
      "$DELEGATION_ADDR" \
      "registerAsOperator((address,address,uint32),string)" \
      "($OPERATOR,0x0000000000000000000000000000000000000000,1)" \
      'https://raw.githubusercontent.com/user_name/repo_name/main/metadata.json'

    go run ../../../cmd -- register-with-avs --config-file ../../../config/local.env
done

# we need to restart the anvil chain at the correct block, otherwise the indexRegistry has a quorumUpdate at the block number
# at which it was deployed (aka quorum was created/updated), but when we start anvil by loading state file it starts at block number 0
# so calling getOperatorListAtBlockNumber reverts because it thinks there are no quorums registered at block 0
# advancing chain manually like this is a current hack until https://github.com/foundry-rs/foundry/issues/6679 is merged
cast rpc anvil_mine 100 --rpc-url $RPC_URL
echo "advancing chain... current block-number:" $(cast block-number)

# Bring Anvil back to the foreground
docker attach anvil
