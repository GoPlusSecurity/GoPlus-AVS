package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type HexBytes []byte
type HexInt64 int64

// 用户对Secware的设定
type SecwareSetting struct {
	SecwareId int    `json:"secware_id"`
	Args      string `json:"args"` // json string 形式，提供具体 Secware 所需的额外参数
}

// SecwareTask 即将由指定的Secware执行的任务
type SecwareTask struct {
	SecwareId      int      `json:"secware_id"`
	SecwareVersion int      `json:"secware_version"`
	SignedTx       HexBytes `json:"signed_tx"`  // 即将发往目标链的已签名交易
	StartTime      HexInt64 `json:"start_time"` // 任务的开始时间
	EndTime        HexInt64 `json:"end_time"`   // 任务的截止时间
	Args           string   `json:"args"`       // json string 形式，提供具体 Secware 所需的额外参数
}

// SignedSecwareTask 被 Gateway 签名的 SecwareTask，是 Secware 的完整输入
type SignedSecwareTask struct {
	Operator   HexBytes    `json:"operator,omitempty"` // Operator 地址
	Task       SecwareTask `json:"task"`
	SigGateway HexBytes    `json:"sig_gateway"` // Gateway 对 Task 的签名
}

// SecwareResult 的各个字段正常情况下由 Secware 填写，Timeout/Crash 时由 Operator 填写。
type SecwareResult struct {
	Code           int      `json:"code"`            // 状态码 0: 正常，1: 超时，2: Crash，>=3: Secware自由使用，表示此交易不安全的各种状态
	Message        string   `json:"message"`         // 状态描述
	Details        string   `json:"result"`          // json string 形式，Secware 输出详细结果。即使没有，也要填写空 JSON `{}`
	Operator       HexBytes `json:"operator"`        // Operator 地址。用于 Secware 生成 HMAC
	SecwareId      int      `json:"secware_id"`      // 表示结果来自哪个 Secware
	SecwareVersion int      `json:"secware_version"` // 表示结果来自哪个版本
}

// SignedSecwareResult 是加入了 Secware 计算的 HMAC 后的完整结果
type SignedSecwareResult struct {
	Result     SecwareResult `json:"result"`
	SigSecware HexBytes      `json:"sig_secware,omitempty"` // 由 SecwareResult 和 Secware私钥 计算出的 HMAC-SHA256
}

// SignedOperatorResult 的 Code，Message 正常时由 Secware 填写，Timeout/Crash 时由 Operator 填写。
// 而 SigSecware，Details 在 Timeout/Crash 时可忽略
type SignedOperatorResult struct {
	Result      SignedSecwareResult `json:"result"`
	SigOperator HexBytes            `json:"sig_operator"`
}

//func (task SecwareTask) MarshalJSON() ([]byte, error) {
//
//	s := fmt.Sprintf(`{"signed_tx":"0x%s","start_time":"0x%x","end_time":"0x%x","args":"0x%s"}`,
//		hex.EncodeToString(task.SignedTx),
//		task.StartTime,
//		task.EndTime,
//		hex.EncodeToString(task.Args))
//	return []byte(s), nil
//}

func NewHexBytesFromString(s string) (HexBytes, error) {
	s = strings.TrimPrefix(s, "0x")
	return hex.DecodeString(s)
}

func (hb HexBytes) MarshalJSON() ([]byte, error) {
	if len(hb) == 0 {
		return []byte(`"0x"`), nil
	}
	return json.Marshal("0x" + hex.EncodeToString(hb))
}

func (hb *HexBytes) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	decoded, err := NewHexBytesFromString(s)
	if err != nil {
		return err
	}
	*hb = decoded
	return nil
}

func (hb HexBytes) String() string {
	return "0x" + hex.EncodeToString(hb)
}

func (h HexInt64) MarshalJSON() ([]byte, error) {
	r := fmt.Sprintf(`"0x%x"`, uint64(h))
	return []byte(r), nil
}

func (h *HexInt64) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	s = strings.TrimPrefix(s, "0x")
	r, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return err
	}
	*h = HexInt64(r)
	return nil
}
