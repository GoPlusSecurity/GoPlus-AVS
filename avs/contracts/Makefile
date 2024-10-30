############################# HELP MESSAGE #############################
# Make sure the help command stays first, so that it's printed by default when `make` is called without arguments
.PHONY: help clean
help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init-repo:  ## download and setup repo and submodules
	git submodule update --init --recursive && cd eigenlayer-middleware && git checkout mainnet && cd lib/eigenlayer-contracts && git checkout mainnet && cd ../../../

build-contracts:  ## builds all contracts
	forge build

deploy-eigenlayer-to-anvil-and-save-state: build-contracts  ## Deploy eigenlayer
	bash ./tests/anvil/deploy-eigenlayer-save-anvil-state.sh

deploy-avs-to-anvil-and-save-state: build-contracts  ## Deploy avs
	bash ./tests/anvil/deploy-avs-save-anvil-state.sh

start-anvil-with-contracts-deployed:  ## starts anvil from a saved state file (with el and avs contracts deployed)
	bash ./tests/anvil/start-anvil-with-el-and-avs-deployed.sh

start-anvil-with-prepared-operators:  ## register 3 test operators to EL.
	bash ./tests/anvil/start-anvil-with-prepared-operators.sh

test-avs:  ## start GoPlus AVS based on eigenlayer and AVS deployed anvil state
	bash ./tests/anvil/run-test-based-on-deployed-anvil-state.sh

clean:
	forge clean
	rm -rf cache/ out/ broadcast/ goplus-avs/script/output
	mkdir -p goplus-avs/script/output/31337
	sudo rm -rf tests/anvil/avs-and-eigenlayer-deployed-anvil-state.json
	sudo rm -rf tests/anvil/eigenlayer-deployed-anvil-state.json
