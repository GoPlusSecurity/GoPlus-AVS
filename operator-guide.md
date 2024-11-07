# Operator Guide

This guide contains the steps needed to set up and register your node for GoPlus AVS (testnet/mainnet).

# Minimal system requirements
- 4 CPU
- 8GB Memory
- 20GB Hard disk (Amazon EBS st1)
- Ubuntu 22.04 LTS
- Docker v24 and above
- [Docker compose](https://docs.docker.com/compose/install)
- [Golang 1.23](https://go.dev/dl)
- [EigenLayer CLI](https://github.com/Layr-Labs/eigenlayer-cli)

# Minimal stake requirements
1. GoAltLayer MACH AVS Mainnet - 1 ETH
2. GoAltLayer MACH AVS Testnet - 1 wei

# Supported token strategy
Beacon Chain Ether and all ETH-based LSTs supported by EigenLayer are supported by our AVS.

💡 **Currently, only `quorum[0]` is available. Other quorums will be opened in the future.**

# Currently active AVS
1. [GoPlus AVS Mainnet](https://app.eigenlayer.xyz/avs/0xa3f64d3102a035db35c42a9001bbc83e08c7a366)
2. [GoPlus AVS Testnet](https://holesky.eigenlayer.xyz/avs/0x6e0e0479e177c7f5111682c7025b4412613cd9de)

# Operator setup

## Key generation and wallet funding

1. Follow EigenLayer guide and Install EigenLayer CLI
2. Generate ECDSA and BLS keypair using the following command
    ```bash
    eigenlayer operator keys create --key-type ecdsa [keyname]
    eigenlayer operator keys create --key-type bls [keyname]
    ```

💡 **Please ensure you backup your private keys to a safe location. By default, the encrypted keys will be stored in ~/.eigenlayer/operator_keys/**. Fund at least 0.3 ETH to the ECDSA address generated. It will be required for node registration in the later steps.

## Register on EigenLayer as an operator

💡 You may skip the following steps if you are already a registered operator on the EigenLayer testnet and mainnet.

**You will need to do it once for testnet and once for mainnet.**

1. Create the configuration files needed for operator registration using the following commands. Follow the step-by-step prompt. Once completed, `operator.yaml` and `metadata.json` will be created.
    ```bash
    eigenlayer operator config create
    ```

2. Edit `metadata.json` and fill in your operator's details.
    ```json
    {
      "name": "Example Operator",
      "website": "<https://example.com/>",
      "description": "Example description",
      "logo": "<https://example.com/logo.png>",
      "twitter": "<https://twitter.com/example>"
    }
    ```
3. Upload `metadata.json` to a public URL. Then update the `operator.yaml` file with the url (`metadata_url`). If you need hosting service to host the metadata, you can consider uploading the metadata [gist](https://gist.github.com/) and get the `raw` url.
4. If this is your first time registering this operator, run the following command to register and update your operator
    ```bash
    eigenlayer operator register operator.yaml
    ```
   Upon successful registration, you should see
    ```
    ✅ Operator is registered successfully to EigenLayer
    ```
   If you need to edit the metadata in the future, simply update metadata.json and run the following command
    ```bash
    eigenlayer operator update operator.yaml
    ```
5. After your operator has been registered, it will be reflected on the EigenLayer operator page.
   Testnet: https://holesky.eigenlayer.xyz/operator
   Mainnet: https://app.eigenlayer.xyz/operator

You can also check the operator registration status using the following command.

```bash
eigenlayer operator status operator.yaml
```

# Joining GoPlus AVS

## GoPlus AVS Setup

### Clone the GoPlus AVS repository

Run the following command to clone the [GoPlus AVS operator repository](https://github.com/GoPlusSecurity/GoPlus-AVS)

```bash
git clone https://github.com/GoPlusSecurity/GoPlus-AVS
```

Inside this repository, we have configurations for various GoPlus AVS. Different `.env` configurations determine whether AVS runs on Mainnet or Testnet.

### Prepare Configuration File

- Run `make copy-config`; this command will create an `.env` configuration file in the project's root directory.
- Fill in the configuration settings in `.env`:
    - `COMPOSE_FILE_PATH`: Path where AVS stores Docker Compose files; replace with an empty and permission-appropriate folder path.
    - `OPERATOR_SECRET_KEY, OPERATOR_BLS_SECRET_KEY`: Fill in the prepared private keys.
    - `NODE_CLASS`: AVS node class, defaults to \"xl\" and does not need modification.
    - `API_PORT`: Port for communication with Gateway; any available port is acceptable.
    - `OPERATOR_URL`: URL path for Gateway access. If not using DNS, set it to `http://{Host IP}:{API_PORT}`, for example, `http://8.8.8.8:7890`.
    - `QUORUM_NUMS`: 0
    - `ETH_RPC`: RPC address. The program uses the RPC address to distinguish between the testnet and mainnet. You can use RPC addresses from providers like Alchemy.
    - `REGISTRY_COORDINATOR_ADDR, OPERATOR_STATE_RETRIEVER`: Copy the deployment addresses for the corresponding network from the [README.md](./README.md).

Example `.env` for Mainnet:

```
COMPOSE_FILE_PATH=/home/user/secwares
OPERATOR_SECRET_KEY=<hexstr>
OPERATOR_BLS_SECRET_KEY=<hexstr>
NODE_CLASS=xl
API_PORT=7776
OPERATOR_URL=http://your_operator_ip:7776
ETH_RPC=https://eth-mainnet.g.alchemy.com/v2/<apikey>
QUORUM_NUMS=0
REGISTRY_COORDINATOR_ADDR=0x91228C6361997a5a4da1a01EdDB2F6B604536A32
OPERATOR_STATE_RETRIEVER=0xD5D7fB4647cE79740E6e83819EFDf43fa74F8C31
```

Example `.env` for Testnet:

```
COMPOSE_FILE_PATH=/home/user/secwares
OPERATOR_SECRET_KEY=<hexstr>
OPERATOR_BLS_SECRET_KEY=<hexstr>
NODE_CLASS=xl
API_PORT=7776
OPERATOR_URL=http://your_operator_ip:7776
ETH_RPC=https://eth-holesky.g.alchemy.com/v2/<apikey>
QUORUM_NUMS=0
REGISTRY_COORDINATOR_ADDR=0x61AA80e5891DbfCebD0B78a704F3de996E449FdE
OPERATOR_STATE_RETRIEVER=0x5ce26317F7edCBCBD1a569629af5DC41c1622045
```

💡 Pay attention to `.env` file permissions so that you don't leak private keys.

### To opt-in

💡 Before you opt-in to GoPlus AVS, please ensure that you have the right infrastructure to keep the operator up and running. Non-performing AVS operators may be subjected to ejection out of GoPlus AVS.

1. Run `make build-avs` to compile AVS.
2. Run `make reg-with-avs` to register.

💡 It may take a few minutes for EigenLayer AVS and operator page to be updated This is an automatic process.

### To opt-out

If you no longer want to run the AVS, you can opt out by running `make dereg-with-avs`.


## Start AVS
> The AVS runtime environment can be either as a standalone process or using Docker Compose; Docker Compose is recommended.

1. Start with Docker Compose:
    1. Run `make build-avs-docker` to build the image.
    2. Run `sudo make run-avs-docker` to start. This also starts Prometheus and Grafana. All components use the **Host** network mode, so make sure the API_PORT, 9090, and 3000 ports in the configuration are not in use.

2. Start as a standalone process:
    1. Run `sudo docker login -u goplusavs -p dckr_pat_wRhsTj4U7REe7IFnrgFkAOswjaM` to log in.
    2. Run `make run-avs` to start.


## Check AVS Running Status
1. Connectivity Check
    - Send a request to `{OPERATOR_URL}/avs/ping` to check the connectivity of the AVS web service.

2. Secware Running Status
    - AVS will periodically request Secware configuration from the Gateway and run Secware in Docker. It also regularly reports Secware's health status to the Gateway.
    - Run `sudo docker compose ls` to view Secware’s running status.

If AVS is running in a Docker Compose environment, you can access Grafana at `http://{OPERATOR_URL}:3000` to view monitoring data. The default username and password are `goplus_avs/admin`."
