# GoPlus AVS

## Directory Structure

Uses a multi-module (`mod`) pattern, described by `./go.work`.

- `./shared`: An independent Go module defining data structures and common functionalities for component interactions, such as signing and verification processes. Please pay special attention to the struct definitions within.
- `./avs`: An independent Go module defining the functionalities of AVS.
- `./mock_secware`: An independent Go module defining the hypothetical functions of SecWare.
- `./mock_gateway`: An independent Go module defining the hypothetical functions of Gateway. Provides interfaces identical to those in the production GoPlus Gateway environment.
- `./mock_fanout_service`: An independent Go module defining the hypothetical functions of Fanout Service. Provides interfaces identical to those in the production GoPlus Fanout Service environment.
- `./scripts`: Stores scripts for building and managing Docker, as well as running and managing various service components.

## Important Module-Level Interfaces

### Gateway

The Gateway inputs a signed transaction `shared/pkg/types.SignedTx` from the end user, reads the Secware configuration information of the end user for this transaction, and constructs multiple `shared/pkg/types.SecwareTask`, each corresponding to a Secware.

```go
// SecwareTask is the task to be executed by the specified Secware
type SecwareTask struct {
    SecwareId      int      `json:\"secware_id\"`
    SecwareVersion int      `json:\"secware_version\"`
    SignedTx       HexBytes `json:\"signed_tx\"`  // The signed transaction to be sent to the target chain
    StartTime      HexInt64 `json:\"start_time\"` // Task start time
    EndTime        HexInt64 `json:\"end_time\"`   // Task deadline
    Args           string   `json:\"args\"`       // Additional parameters required by specific Secware, in JSON string format
}
```

The Gateway sends `[]SecwareTask` to the Fanout Service.

### Fanout Service

After receiving `[]SecwareTask` from the Gateway, the Fanout Service uses the Gateway's private key to sign each element in the array, resulting in a set of `shared/pkg/types.SignedSecwareTask`, where the `Operator` field is left empty.

```go
// SignedSecwareTask is the SecwareTask signed by the Gateway, which is the complete input for the Secware
type SignedSecwareTask struct {
    Operator   HexBytes    `json:\"operator,omitempty\"` // Operator's address
    Task       SecwareTask `json:\"task\"`
    SigGateway HexBytes    `json:\"sig_gateway\"` // Gateway's signature on the Task
}
```

The Fanout Service needs to select a group of Operators based on the Fanout policy and send `[]SignedSecwareTask` to each Operator in the group.

### Operator

After receiving `[]SignedSecwareTask` from the Fanout Service, the Operator executes the following for each element:

1. Fills its own on-chain Operator address into the `SignedSecwareTask.Operator` field.
2. According to the `(SecwareId, SecwareVersion)` described in `SignedSecwareTask.Task`, sends each `SignedSecwareTask` to the corresponding Secware.

### Secware

After receiving `SignedSecwareTask` from the Operator, the Secware executes:

1. Verifies whether `SignedSecwareTask.SigGateway` is valid.
2. Verifies whether `(SecwareId, SecwareVersion)` matches itself.
3. Parses the extra parameters in JSON format from `SecwareTask.Args`.
4. Executes the Task.
5. Writes the result into `shared/pkg/types.SecwareResult`.
6. Combines the Operator address and execution result to compute an HMAC on `SecwareResult` using its own private key, and fills it into the `SecwareResult.SigSecware` field.
7. Returns `SecwareResult` to the Operator.

```go
// SecwareResult fields are normally filled by the Secware; in case of Timeout/Crash, they are filled by the Operator.
type SecwareResult struct {
    Code     int      `json:\"code\"`     // Status code: 0 - Normal, 1 - Timeout, 2 - Crash, >=3 - Used freely by Secware to indicate various unsafe states of the transaction
    Message  string   `json:\"message\"`  // Status description
    Details  string   `json:\"result\"`   // Detailed result output by Secware in JSON string format. Even if empty, should be filled with an empty JSON '{}'
    Operator HexBytes `json:\"operator\"` // Operator's address. Used by Secware to generate HMAC
}

// SignedSecwareResult is the complete result after adding the HMAC computed by Secware
type SignedSecwareResult struct {
    Result     SecwareResult `json:\"result\"`
    SigSecware HexBytes      `json:\"sig_secware,omitempty\"` // HMAC-SHA256 computed from SecwareResult and Secware's private key
}
```

### Operator Submitting Results

Finally, the Operator will, before `SecwareTask.EndTime`:

1. Summarize the execution results from each Secware (including Timeout/Crash), where each Secware corresponds to a `SecwareResult`.
2. Sign the contents of `[]SecwareResult` using its own private key.
3. Fill the `shared/pkg/types.SignedOperatorResult` structure and send it to the Gateway.

### Gateway Summarizing Results

The Gateway waits and summarizes the execution results submitted by Operators before each `SecwareTask.EndTime`. Then it performs:

1. Completes the results of Operators that did not respond.
2. Computes consensus.
3. If the consensus result is safe, it will forward `SignedTx` to the target network, wait for and collect the execution results from the target network.
4. Summarizes the execution results of each step and returns them to the end user.

## Signing Methods

### Gateway Signature

Serialize `SecwareTask` in JSON according to the order of field declarations, without spaces between fields. Fields of types `HexBytes` and `HexInt64` are encoded as strings starting with `0x`. Even if `HexBytes` is empty, it retains `0x`.

Compute the SHA3 hash of the serialized `[]byte`, then perform an ECDSA signature. This process is the same as the signing process in go-ethereum.

### Operator Signature

Same as the Gateway signature, but the object is `SignedSecwareResult`.

### Secware HMAC

Serialize `SecwareResult` in JSON according to the order of field declarations to get `msg`, use the private key held by this version of Secware as `key`, compute `HMAC-SHA256(msg, key)`.

## Secware Version Update Process

Operators regularly request the currently enabled `(SecwareId, SecwareVersion)` pairs from the Gateway, and find:

- If there are pairs they don't hold, they pull the Docker images.
- If there are pairs they've already started but are not in the list returned by the Gateway, they shut down the corresponding Docker compose.

## mock_secware

- **Configuration**: `./mock_secware/pkg/config/config.go`
- **Entry point**: `./mock_secware/cmd/main.go`
- **Test**: `./mock_secware/test/task_test.go`

The input of the HTTP RPC is `defs.SignedSecwareTask`, and you can predefine the behavior of `mock_secware` through `defs.SignedSecwareTask.Args`.

The `Args` field is the JSON serialization of `mock_secware/handlers.SecwarerArgs`.

```go
// SecwarerArgs are additional parameters provided to mock_secware, coming from SignedSecwareTask.Task.Args (JSON)
// - result: string, indicates the predefined security audit result returned by mock_secware, used for debugging, such as \"Yes\", \"No\", etc.
// - sleep: int, indicates the duration (in seconds) that mock_secware waits before executing actions
// - crash: bool, indicates whether mock_secware ignores the return and actively crashes
type SecwarerArgs struct {
    Result defs.HexBytes `json:\"result,omitempty\"`
    Sleep  int           `json:\"sleep,omitempty\"`
    Crash  bool          `json:\"crash,omitempty\"`
}
```

---

## Holesky Deployment

```
Deployer: 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939
GoPlusProxyAdmin: 0xdf9EE7B28fef9aEe47f52DeA24e6eBEfECc9EaC2
    owner: 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939
GoPlusServiceManager: 0xC3c5934686254A59C3B9ce40CFa9F36c1a0BeFf9
    Impl: 0x7e186EbCe5c31D297A478639E7dF8bfe4b075b3A
    owner: 0xdf9EE7B28fef9aEe47f52DeA24e6eBEfECc9EaC2
RegistryCoordinator: 0x3C503C651e3BD82C7AD169411E674d8ea6ad07e6
    Impl: 0x07510c426105e3023aB859764532e687543f1Fe2
BLSApkRegistry: 0xf89d6536994682260b8D98349218eF6cb0159824
    Impl: 0xA01c421b15345085b55b1d9c0Be01104e6DdA85F
IndexRegistry: 0xCce02fb16b1F9893951DD49Ecd5941BcC4Ef8D5A
    Impl: 0x1e43Dc9046f2D1C00Fa8bb506F19B3103c4f567A
StakeRegistry: 0x0965C97ED9DBB76a102b4F1fa1A5dBA2cBd802f0
    Impl: 0xde530efef354F6D2cA18867adf271610Dd22b1f5
OperatorStateRetriever: 0x5ce26317F7edCBCBD1a569629af5DC41c1622045
PauserRegistry: 0xc2284B80Cf95BaD900dd0c31d0a4660b3A4Bb8cC
    Constructor parameters: [0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939], 0x15fbbC47a244aE2A38071A106dCfcF3D57C9D939
```

`ServiceManager` and `ECDSAStakeRegistry` are TUP contracts; the former handles the creation and response of Tasks and interacts with EigenLayer. The latter provides the joining and exiting of Operators.

Both of these TUP contracts have `ProxyAdmin` as the owner, and the owner of `ProxyAdmin` is the `deployer`."
                    