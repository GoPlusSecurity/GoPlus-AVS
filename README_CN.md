# GoPlus AVS

## 总览
- 如果要在 Holesky testnet 上启动 GoPlus AVS，请参考文档 `deployment.md`
- 如果要开发 GoPlus AVS 的后台与合约，请参考本文档。
- 如果要部署 GoPlus AVS 到 mainnet，需要明确一些基础设定：
  - Deployer
  - ProxyAdmin owner
  - Pauser group
  - ChurnApprover
  - Ejector
  - Quorum
    - OperatorSetParam
    - StrategyParam
    - MinimumStake

## 目录结构
采用多 mod 的模式。由 `./go.work` 描述。

- `./shared` 一个独立的Go模块，定义了各组件交互的数据结构和通用功能，比如签名和验证过程。请重点关注其中的结构体定义。
- `./avs` 一个独立的Go模块，定义了AVS的功能。
- `./mock_secware` 一个独立的Go模块，定义了假想的 SecWare 职能。
- `./mock_gateway` 一个独立的Go模块，定义了假想的 Gateway 职能。提供与生产环境中 GoPlus Gateway 相同的接口。
- `./mock_fanout_service` 一个独立的Go模块，定义了假想的 Fanout Service 职能。提供与 prod 环境中 GoPlus Fanout Service 相同的接口。
- 
- `./scripts` 存放建构和管理 docker 的脚本，以及各个服务组件的运行管理

## 重要的模块级接口

### Gateway

Gateway 输入来自 End user 的已签名交易 `shared/pkg/types.SignedTx`，读取 End user 对于此交易的 Secware 配置信息，构造多个 `shared/pkg/types.SecwareTask`，每个对应于一个 Secware。 

```go
// SecwareTask 即将由指定的Secware执行的任务
type SecwareTask struct {
	SecwareId      int      `json:"secware_id"`
	SecwareVersion int      `json:"secware_version"`
	SignedTx       HexBytes `json:"signed_tx"`  // 即将发往目标链的已签名交易
	StartTime      HexInt64 `json:"start_time"` // 任务的开始时间
	EndTime        HexInt64 `json:"end_time"`   // 任务的截止时间
	Args           string   `json:"args"`       // json string 形式，提供具体 Secware 所需的额外参数
}
```

Gateway 将 `[]SecwareTask` 发送给 Fanout service。

### Fanout service

Fanout service 接到来自 Gateway 的 `[]SecwareTask` 后，使用 Gateway 的私钥进行对数组中的每一个元素进行签名，从而得到一组 `shared/pkg/types.SignedSecwareTask`，其中 Operator 字段不填。

```go
// SignedSecwareTask 被 Gateway 签名的 SecwareTask，是 Secware 的完整输入
type SignedSecwareTask struct {
	Operator   HexBytes    `json:"operator,omitempty"` // Operator 的地址
	Task       SecwareTask `json:"task"`
	SigGateway HexBytes    `json:"sig_gateway"` // Gateway 对 Task 的签名
}
```
Fanout service 需要基于 Fanout policy 选择出一组 Operator 进行，将 `[]SignedSecwareTask` 发送给组中的每一个 Operator。

### Operator

Operator 收到来自 Fanout service 的 `[]SignedSecwareTask` 后，对其中每个元素执行：

1. 填写自身的 Operator 链上地址到 `SignedSecwareTask.Operator` 字段；
2. 依据 `SignedSecwareTask.Task` 中描述的 `(SecwareId, SecwareVersion)`，将每个 `SignedSecwareTask` 发送给对应的 Secware。


### Secware

Secware 收到来自 Operator 的 `SignedSecwareTask` 后，执行：

1. 验证 `SignedSecwareTask.SigGateway` 是否合法；
2. 验证 `(SecwareId, SecwareVersion)` 是否与自身一致；
3. 解析 JSON 格式的额外参数 `SecwareTask.Args`；
4. 执行 Task；
5. 将结果写入 `shared/pkg/types.SecwareResult`；
6. 结合 Operator 地址以及执行结果，对 `SecwareResult` 用自身私钥计算 HMAC，填入 `SecwareResult.SigSecware` 字段；
7. 将 `SecwareResult` 返回给 Operator。

```go
// SecwareResult 的各个字段正常情况下由 Secware 填写，Timeout/Crash 时由 Operator 填写。
type SecwareResult struct {
	Code     int      `json:"code"`     // 状态码 0: 正常，1: 超时，2: Crash，>=3: Secware自由使用，表示此交易不安全的各种状态
	Message  string   `json:"message"`  // 状态描述
	Details  string   `json:"result"`   // json string 形式，Secware 输出详细结果。即使没有，也要填写空 JSON `{}`
	Operator HexBytes `json:"operator"` // Operator 地址。用于 Secware 生成 HMAC
}

// SignedSecwareResult 是加入了 Secware 计算的 HMAC 后的完整结果
type SignedSecwareResult struct {
	Result     SecwareResult `json:"result"`
	SigSecware HexBytes      `json:"sig_secware,omitempty"` // 由 SecwareResult 和 Secware私钥 计算出的 HMAC-SHA256
}
```

### Operator 提交结果

最终 Operator 会在 `SecwareTask.EndTime` 之前执行：

1. 将各个 Secware 执行的结果（含 Timeout / Crash）进行汇总，每个 Secware 对应一个 `SecwareResult`；
2. Operator 对 `[]SecwareResult` 的内容用自己私钥进行签名；
3. 填写 `shared/pkg/types.SignedOperatorResult` 结构体，将其发送给 Gateway。

### Gateway 汇总结果

Gateway 在每个 `SecwareTask.EndTime` 之前等待并汇总 Operator 提交的执行结果。然后执行：

1. 补全未响应的 Operator 的结果；
2. 计算共识；
3. 如果共识结果为安全，则会转发 `SignedTx` 到目标网络，等待并收集目标网络执行结果；
4. 将各个步骤的执行结果汇总返回给 End user。

## 签名方式

### Gateway 签名

按 `SecwareTask` 字段声明的顺序进行 JSON 序列化，字段间不留空格。同时 `HexBytes` 和 `HexInt64` 类型的字段都编码为 `0x` 开头的字符串。`HexBytes` 即使为空，也保留 `0x`。
对 JSON 序列化后的 `[]byte` 取 SHA3，然后做 ECDSA 签名。这个过程跟 go-ethereum 中的签名流程相同。 

### Operator 签名

同 Gateway 签名，只不过对象为 `SignedSecwareResult`。 

### Secware HMAC

将 `SecwareResult` 按字段声明的顺序进行 JSON 序列化为 `msg`，以此版本的 Secware 持有的私钥为 `key`，计算 `HMAC-SHA256(msg, key)`。

## Secware 版本更新流程

Operator 定期向 Gateway 请求目前启用的各个 `(SecwareId, SecwareVersion)`，发现
- 自己并不持有的 pair 时，去拉取 docker images；
- 自己已经启动但 Gateway 返回列表中没有的 pair 时，关闭对应的 docker compose。

## mock_secware

- 配置: `./mock_secware/pkg/config/config.go`
- 入口: `./mock_secware/cmd/main.go`
- 测试: `./mock_secware/test/task_test.go`

HTTP RPC 的输入为 `defs.SignedSecwareTask`，可通过 `defs.SignedSecwareTask.Args` 来提前定义 mock_secware 的行为。
`Args` 字段为 `mock_secware/handlers.SecwarerArgs` 的 JSON 序列化的只。 

```go
// SecwarerArgs 是提供给 mock_secware 的额外参数，来自 SignedSecwareTask.Task.Args (JSON)
// - result: string 表示预先指定 mock_secware 返回的安全审查结果, 供调试使用, 比如 "Yes", "No", ...
// - sleep: int  表示 mock_secware 执行动作前等待的时长（秒）
// - crash: bool 表示 mock_secware 是否忽略 return 而主动崩溃
type SecwarerArgs struct {
	Result defs.HexBytes `json:"result,omitempty"`
	Sleep  int           `json:"sleep,omitempty"`
	Crash  bool          `json:"crash,omitempty"`
}
```



