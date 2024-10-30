# AVS 部署

## 1. 环境准备
- 准备服务器，需要公网ip，最小配置要求，4核CPU 8GB内存
- 安装docker compose 
  - https://docs.docker.com/compose/install/
- 安装golang 1.23
  - https://go.dev/dl/

---
## 2. 准备私钥
- 准备私钥
  1. 完成 EigenLayer Operator 注册
  2. 完成资产质押，质押量应满足最小值
  3. 确保余额充足，用于支付AVS注册的交易费用
- 准备BLS私钥：此私钥目前暂时不起作用，但依旧要妥善保管，因为会用于生成 OperatorId

---
## 3. 准备配置文件
- 执行`make copy-config`，此命令将在项目根目录下创建 `.env` 配置文件
- 填充`.env`中各项配置
  - `COMPOSE_FILE_PATH` AVS存放docker compose文件的路径，准备空且有权限的文件夹路径替换
  - `OPERATOR_SECRET_KEY, OPERATOR_BLS_SECRET_KEY` 将准备好的私钥分别填充
  - `NODE_CLASS`AVS节点规格, 默认为xl, 无需修改
  - `API_PORT`与Gateway的通信端口, 任意空闲端口即可
  - `OPERATOR_URL`Gateway访问的路径, 如果不使用DNS, 则填写http://{主机ip}:{API_PORT}即可, 例如http://8.8.8.8:7890
  - `ETH_RPC`rpc地址, 程序使用rpc地址区分测试网和主网, 可以使用alchemy等供应商提供的rpc地址 
  - `REGISTRY_COORDINATOR_ADDR, OPERATOR_STATE_RETRIEVER` 复制 [Holesky部署位置](#Holesky部署位置) 中对应的地址

---
## 4. 注册AVS
1. 执行 `make build-avs` 编译avs
2. 执行 `make reg-with-avs` 注册

---
## 5. 启动AVS
> AVS执行环境分为裸进程和docker compose两种方式，推荐使用docker compose

1. docker compose启动
   1. 执行 `make build-avs-docker` 构建镜像
   2. 执行 `sudo make run-avs-docker` 启动，同时还会启动 Prometheus 和 Grafana，所有组件均使用**Host网络模式**启动，所以需确保配置中的 `API_PORT`, `9090` 和 `3000` 端口未被占用
2. 裸进程启动
   1. 执行 `sudo docker login -u goplusavs -p dckr_pat_wRhsTj4U7REe7IFnrgFkAOswjaM` 登录
   2. 执行 `make run-avs` 启动

---
## 6.查看AVS运行状态
1. 连通性检查
   - 请求 `{OPERATOR_URL}/avs/ping` 接口检查 AVS 的WEB服务连通性
2. Secware运行状态
   - AVS 会定时向 Gateway 请求 Secware 配置，并在 docker 中运行 Secware，同时会定时向 Gateway 报告 Secware 的健康状态
   - 执行 `sudo docker compose ls` 查看 Secware 运行状态

如果使用 docker compose 为 AVS 的运行环境，可以通过 `http://{OPERATOR_URL}:3000` 访问 Grafana 查看监控数据，默认用户名密码为 `goplus_avs/admin`

---
## Holesky部署位置

目前 GoPlus AVS 已经部署在 Holesky testnet 上了，各个合约的部署地址如下：

```
Deployer: 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939
GoPlusProxyAdmin: 0xdf9EE7B28fef9aEe47f52DeA24e6eBEfECc9EaC2
GoPlusServiceManager: 0xC3c5934686254A59C3B9ce40CFa9F36c1a0BeFf9
RegistryCoordinator: 0x3C503C651e3BD82C7AD169411E674d8ea6ad07e6
BLSApkRegistry: 0xf89d6536994682260b8D98349218eF6cb0159824
IndexRegistry: 0xCce02fb16b1F9893951DD49Ecd5941BcC4Ef8D5A
StakeRegistry: 0x0965C97ED9DBB76a102b4F1fa1A5dBA2cBd802f0
OperatorStateRetriever: 0x5ce26317F7edCBCBD1a569629af5DC41c1622045
PauserRegistry: 0xc2284B80Cf95BaD900dd0c31d0a4660b3A4Bb8cC
```

---
## 本地部署合约

> 由于 GoPlus AVS 的合约已经部署到了 Holesky testnet 上，默认情况下，AVS 后台直接与已部署合约通信。
> 因此仅在需要调试开发 GoPlus AVS 合约时才需要本地部署合约。

1. 从 Holesky testnet fork 一条本地链：

    ```bash
    anvil --fork-url <rpc> --fork-block-number 2031600
    ```

2. 向本地链部署 GoPlus AVS 合约：

    ```bash
    cd avs/contracts
    # 用于获取 eigenlayer-middleware submodule 中的各种文件
    make init-repo
   
    # 执行部署脚本
    forge script --rpc-url 'http://localhost:8545' --sig 'runHolesky()' --keystore <keyfile> --broadcast -vvv goplus-avs/script/GoPlusDeployer.s.sol
    ```

3. 更新配置文件

    请参考 `deployment.md` 文档，更新 `.env` 与本地部署参数对应。
