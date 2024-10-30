#!/bin/bash

root_path=$(
    cd "$(dirname "${BASH_SOURCE[0]}")/../../"
    pwd -P
)
cd "$root_path"

cleanup() {
    echo "Executing cleanup function..."
    set +e
    exit_status=$?
    if [ $exit_status -ne 0 ]; then
        echo "Script exited due to set -e on line $1 with command '$2'. Exit status: $exit_status"
    fi
}
trap 'cleanup $LINENO "$BASH_COMMAND"' EXIT

CHAIN_ID=$(cast chain-id)
DELEGATION_ADDR=$(jq -r '.addresses.delegation' contracts/script/output/"$CHAIN_ID"/eigenlayer_deployment_output.json)
AVSDIRECTORY_ADDR=$(jq -r '.addresses.avsDirectory' contracts/script/output/"$CHAIN_ID"/eigenlayer_deployment_output.json)
SERVICE_MANAGER_ADDR=$(jq -r '.addresses.HelloWorldServiceManagerProxy' contracts/script/output/"$CHAIN_ID"/hello_world_avs_deployment_output.json)
STAKER_REGISTRY_ADDR=$(jq -r '.addresses.ECDSAStakeRegistry' contracts/script/output/"$CHAIN_ID"/hello_world_avs_deployment_output.json)
OPERATOR=$(cast wallet address --account goplus-operator --password test)
SALT="0x0000000000000000000000000000000000000000000000000000000000000000"
EXPIRY=9999999999

echo "Register operator to EigenLayer..."
cast send \
	--json \
	--account goplus-operator --password test \
	"$DELEGATION_ADDR" \
	"registerAsOperator((address,address,uint32),string)" \
	"(${OPERATOR},0x0000000000000000000000000000000000000000,1)" \
	'https://raw.githubusercontent.com/user_name/repo_name/main/metadata.json'

DIGEST=$(cast call \
  "$AVSDIRECTORY_ADDR" \
  "calculateOperatorAVSRegistrationDigestHash(address,address,bytes32,uint256)(bytes32)" \
  "$OPERATOR" \
  "$SERVICE_MANAGER_ADDR" \
  "$SALT" \
  "$EXPIRY")

SIG=$(cast wallet sign --account goplus-operator --password test --no-hash "$DIGEST")

echo "Register operator to AVS..."
cast send \
	--json \
	--account goplus-operator --password test \
	"$STAKER_REGISTRY_ADDR" \
	"registerOperatorWithSignature((bytes,bytes32,uint256),address)" \
	"(${SIG},${SALT},${EXPIRY})" \
	"$OPERATOR"

