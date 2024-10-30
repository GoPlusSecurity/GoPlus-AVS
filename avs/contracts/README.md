# GoPlus AVS contracts

## 目录

- abi: 存放必要的 ABI 文件
- bindings: 存放 ABI 对应的 Go binding 文件
- goplus-avs: 是对 https://github.com/Layr-Labs/hello-world-avs 的改进版本

## 用法

1. 找一个空目录，git clone https://github.com/Layr-Labs/hello-world-avs
2. 将 goplus-avs 目录中的所有内容覆盖到 hello-world-avs 中，并进入此目录。
3. `make create-accounts` 创建两个账号 `goplus-deployer` 和 `goplus-operator`。
3. `make start-from-el-deployed` 启动 anvil 本地测试网，部署好 EigenLayer 合约。
4. `make deploy-avs` 会使用 `goplus-deployer` 账号部署 AVS 相关的合约。
5. `make register-operator` 会使用 `goplus-operator` 账号向 EigenLayer 和 AVS 注册为合法的 Operator。
