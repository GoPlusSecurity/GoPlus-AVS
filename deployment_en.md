# AVS Deployment

## 1. Environment Preparation
- Prepare a server with a public IP; minimum configuration required is 4 CPUs and 8GB RAM.
- Install Docker Compose 
  - https://docs.docker.com/compose/install/
- Install Golang 1.23
  - https://go.dev/dl/

---

## 2. Prepare Private Key
- Prepare the private key
  1. Complete the registration as an Eigen Layer operator.
  2. Complete staking, ensuring the staked amount meets the minimum requirement.
  3. Ensure sufficient balance to cover the transaction fees for AVS registration.
- Prepare the BLS private key.

---

## 3. Prepare Configuration File
- Run `make copy-config`; this command will create an `.env` configuration file in the project's root directory.
- Fill in the configuration settings in `.env`:
  - `COMPOSE_FILE_PATH`: Path where AVS stores Docker Compose files; replace with an empty and permission-appropriate folder path.
  - `OPERATOR_SECRET_KEY, OPERATOR_BLS_SECRET_KEY`: Fill in the prepared private keys.
  - `NODE_CLASS`: AVS node class, defaults to \"xl\" and does not need modification.
  - `API_PORT`: Port for communication with Gateway; any available port is acceptable.
  - `OPERATOR_URL`: URL path for Gateway access. If not using DNS, set it to `http://{Host IP}:{API_PORT}`, for example, `http://8.8.8.8:7890`.
  - `ETH_RPC`: RPC address. The program uses the RPC address to distinguish between the testnet and mainnet. You can use RPC addresses from providers like Alchemy.
  - `REGISTRY_COORDINATOR_ADDR, OPERATOR_STATE_RETRIEVER`: Copy the deployment addresses for the corresponding network from the [README.md](./README.md).

---

## 4. Register AVS
1. Run `make build-avs` to compile AVS.
2. Run `make reg-with-avs` to register.

---

## 5. Start AVS
> The AVS runtime environment can be either as a standalone process or using Docker Compose; Docker Compose is recommended.

1. Start with Docker Compose:
   1. Run `make build-avs-docker` to build the image.
   2. Run `sudo make run-avs-docker` to start. This also starts Prometheus and Grafana. All components use the **Host** network mode, so make sure the API_PORT, 9090, and 3000 ports in the configuration are not in use.
   
2. Start as a standalone process:
   1. Run `sudo docker login -u goplusavs -p dckr_pat_wRhsTj4U7REe7IFnrgFkAOswjaM` to log in.
   2. Run `make run-avs` to start.

---

## 6. Check AVS Running Status
1. Connectivity Check
   - Send a request to `{OPERATOR_URL}/avs/ping` to check the connectivity of the AVS web service.
   
2. Secware Running Status
   - AVS will periodically request Secware configuration from the Gateway and run Secware in Docker. It also regularly reports Secware's health status to the Gateway.
   - Run `sudo docker compose ls` to view Secware’s running status.

If AVS is running in a Docker Compose environment, you can access Grafana at `http://{OPERATOR_URL}:3000` to view monitoring data. The default username and password are `goplus_avs/admin`."
                    ]