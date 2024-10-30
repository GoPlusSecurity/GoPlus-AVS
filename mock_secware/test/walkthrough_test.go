package mock_secware_test

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"goplus/shared/pkg/signature"
	"goplus/shared/pkg/types"
	"log"
	"testing"
)

const NumOfSecware = 2
const SkGatewayHex = "fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"  // Gateway的私钥
const SkOperatorHex = "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" // Operator的私钥
const Operator = "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"                            // Operator 的ETH地址
const SkSecwareHex1 = "B66EC94BBF121BA0416925AB02A97E4A2EF9BB9835C6170B83A79B96F793340A" // 用于生成 HMAC
const SkSecwareHex2 = "32D9699552C32204AE4A497A7EADA0F21BCA97205D63BC1BA979137C8AB256D7" // 用于生成 HMAC

func Test_Walkthrough(t *testing.T) {
	// 假设用户预先设定要使用两个 Secware 做安全检测，他们的配置为：
	// SkSecwareHex := [NumOfSecware]string{SkSecwareHex1, SkSecwareHex2}
	skGateway, _ := crypto.HexToECDSA(SkGatewayHex)
	pkGateway := skGateway.Public().(*ecdsa.PublicKey)
	skOperator, _ := crypto.HexToECDSA(SkOperatorHex)
	pkOperator := skOperator.Public().(*ecdsa.PublicKey)

	user_settings := [NumOfSecware]types.SecwareSetting{
		{
			SecwareId: 1,
			Args:      "",
		},
		{
			SecwareId: 2,
			Args:      "{\"level\": 3}",
		},
	}

	// 用户输入，就是一个已签名交易的 bytes 表示
	signedTx, err := types.NewHexBytesFromString("")
	if err != nil {
		log.Fatal(err)
	}
	// 用户将其发给 Gateway。后者根据用户配置构造一组 SecwareTask

	// Gateway 会构造出两个 SecwareTask
	rawTaskList := make([]types.SecwareTask, NumOfSecware)
	for i := 0; i < NumOfSecware; i++ {
		rawTaskList[i] = types.SecwareTask{
			SecwareId:      user_settings[i].SecwareId,
			SecwareVersion: 1, // 从配置中心获得, 目前假设都是第一版
			SignedTx:       signedTx,
			StartTime:      0x12345678,
			EndTime:        0x87654321,
			Args:           user_settings[i].Args,
		}
	}

	// Gateway 对每个 SecwareTask 签名
	signedTaskList := make([]types.SignedSecwareTask, NumOfSecware)
	for i := 0; i < NumOfSecware; i++ {
		signedTask, _ := signature.SignSecwareTask(&rawTaskList[i], skGateway)
		signedTaskList[i] = *signedTask
	}

	// Gateway 把 `signedTaskList` 发送给 Fanout service，后者进一步分发给多个 Operator 的 AVS HTTP Endpoint。
	// 假设一个 Task 有 2 个 Secware 需要执行，有 3 个 Operator 提供 AVS，那么 Fanout service 就发起 2x3=6 次 HTTP 请求。
	// 每个 HTTP 请求在指定时间窗口内等待 Operator 的返回。

	// AVS 会收到 signedTask。（可选）检查签名是否来自 Gateway，然后填入 Operator 的地址。
	// 然后调用具体的 Secware 获得结果，并签名。
	// 最后在当前 HTTP 连接中将结果返回给 Fanout service。
	signedOperatorResultList := make([]types.SignedOperatorResult, NumOfSecware)
	for idx, signedTask := range signedTaskList {
		signedTaskBody, _ := json.Marshal(signedTask)
		log.Printf("signedTask: %s", string(signedTaskBody))

		if !signature.VerifySignedSecwareTask(&signedTask, pkGateway) {
			log.Fatal("bad SignedSecwareTask")
		}
		skOp, err := types.NewHexBytesFromString(Operator)
		if err != nil {
			log.Fatal(err)
		}
		signedTask.Operator = skOp

		// 调用具体的 Secware
		signedSecwareResult := callSecware(signedTask)

		// Operator 对结果签名并作为 HTTP Response 返回给 Fanout service
		signedOperatorResult, _ := signature.SignOperatorResult(signedSecwareResult, skOperator)

		// Fanout service 收集结果
		signedOperatorResultList[idx] = *signedOperatorResult
	}

	// Fanout service 收到 Operator 的结果后对签名和各个 Secware 的 HMAC 进行验证
	for _, signedOperatorResult := range signedOperatorResultList {
		if !signature.VerifyOperatorSignature(&signedOperatorResult, pkOperator) {
			log.Fatal("bad SignedOperatorResult")
		}

		signedSecwareResult := signedOperatorResult.Result
		var skSecware []byte
		if signedSecwareResult.Result.SecwareId == 1 {
			skSecware = []byte(SkSecwareHex1)
		} else {
			skSecware = []byte(SkSecwareHex2)
		}

		if !signature.VerifySecwareSignature(&signedSecwareResult, skSecware) {
			log.Fatal("bad SignedSecwareResult")
		}
	}

	// Fanout service 在时间窗口内，汇总各个 Operator 的结果，执行共识并返回给 Gateway
	// Gateway 在将结果返回或进一步转发 SignedTx 到目标网络。
}

func callSecware(signedTask types.SignedSecwareTask) types.SignedSecwareResult {
	// Secware 检查签名
	skGateway, _ := crypto.HexToECDSA(SkGatewayHex)
	pkGateway := skGateway.Public().(*ecdsa.PublicKey)

	if !signature.VerifySignedSecwareTask(&signedTask, pkGateway) {
		log.Fatal("bad SignedSecwareTask")
	}

	// Secware 解析参数
	task := signedTask.Task
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(task.Args), &args); err != nil {
		log.Fatal(err)
	}

	// 执行具体任务
	// 构造返回
	result := types.SecwareResult{
		Code:           0,
		Message:        "ok",
		Details:        fmt.Sprintf("{\"my_id\": %d}", task.SecwareId), // 随便返回什么
		Operator:       signedTask.Operator,
		SecwareId:      task.SecwareId,
		SecwareVersion: 1,
	}

	var skSecware []byte
	if task.SecwareId == 1 {
		skSecware = []byte(SkSecwareHex1)
	} else {
		skSecware = []byte(SkSecwareHex2)
	}

	// Secware 对返回结果生成 HMAC
	signedResult, _ := signature.SignSecwareResult(&result, skSecware)

	// 填充自身标识
	return *signedResult
}
