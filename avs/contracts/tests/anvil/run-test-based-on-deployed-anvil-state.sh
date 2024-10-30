#!/bin/bash

# cd to the directory of this script so that this can be run from anywhere
parent_path=$(
    cd "$(dirname "${BASH_SOURCE[0]}")"
    pwd -P
)
# At this point we are in tests/anvil
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

# start an empty anvil chain in the background and dump its state to a json file upon exit
start_anvil_docker "${parent_path}/avs-and-eigenlayer-deployed-anvil-state.json" ""

cd ../../

forge test --fork-url http://localhost:8545 -vv
