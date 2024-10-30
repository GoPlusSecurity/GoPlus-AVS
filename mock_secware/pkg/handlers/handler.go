package handlers

import (
	"encoding/json"
	"fmt"
	"goplus/shared/pkg/signature"
	"goplus/shared/pkg/types"
	"net/http"
	"os"
	"time"
)

//func OnHome(w http.ResponseWriter, r *http.Request) {
//	_, err := io.WriteString(w, "Home")
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}

// SecwareArgs 是提供给 Secware 的额外参数，来自 SignedSecwareTask.Task.Args (JSON)
// - result: string 表示预先指定 Secware 返回的安全审查结果, 供调试使用, 比如 "Yes", "No", ...
// - sleep: int  表示 Secware 执行动作前等待的时长（秒）
// - crash: bool 表示 Secware 是否忽略 result 去主动崩溃
type SecwareArgs struct {
	Result string `json:"result,omitempty"`
	Sleep  int    `json:"sleep,omitempty"`
	Crash  bool   `json:"crash,omitempty"`
}

// OnSecwareTask 处理 Operator 发来的请求
func OnSecwareTask(w http.ResponseWriter, r *http.Request) {
	var task *types.SignedSecwareTask

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, fmt.Sprintf("could not unmarshal json: %s\n", err), http.StatusBadRequest)
		return
	}

	if input, err := json.Marshal(task); err != nil {
		fmt.Printf("input: %s\n", input)
	}

	if len(task.Operator) != 20 {
		http.Error(w, "invalid operator address", http.StatusBadRequest)
		return
	}

	setting, ok := r.Context().Value("Setting").(*SecwareSetting)
	if !ok {
		http.Error(w, "could not get Setting", http.StatusInternalServerError)
	}

	// check signature
	if !signature.VerifySignedSecwareTask(task, setting.PkGateway) {
		http.Error(w, "signature verification failed", http.StatusBadRequest)
		return
	}

	// 解析额外参数
	var args SecwareArgs
	if err := json.Unmarshal([]byte(task.Task.Args), &args); err != nil {
		http.Error(w, fmt.Sprintf("could not unmarshal args: %s\n", err), http.StatusBadRequest)
		return
	}

	if args.Sleep > 0 {
		time.Sleep(time.Duration(args.Sleep) * time.Second)
	}

	if args.Crash {
		fmt.Printf("crashing")
		os.Exit(1)
	}

	result, err := signature.SignSecwareResult(&types.SecwareResult{
		Code:     0,
		Message:  "ok",
		Details:  args.Result,
		Operator: task.Operator,
	}, setting.SkSecwareBytes)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not sign result: %s\n", err), http.StatusInternalServerError)
		return
	}

	output, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not marshal result: %s\n", err), http.StatusInternalServerError)
	}
	fmt.Printf("output: %s\n", output)

	// 发送结果
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, fmt.Sprintf("could not encode json: %s\n", err), http.StatusInternalServerError)
	}
}

func OnGetMeta(w http.ResponseWriter, r *http.Request) {
	setting, ok := r.Context().Value("Setting").(*SecwareSetting)
	if !ok {
		http.Error(w, "could not get Setting", http.StatusInternalServerError)
	}

	meta := struct {
		SecwareId      int `json:"secware_id"`
		SecwareVersion int `json:"secware_version"`
	}{
		SecwareId:      setting.SecwareId,
		SecwareVersion: setting.SecwareVersion,
	}

	if err := json.NewEncoder(w).Encode(meta); err != nil {
		http.Error(w, fmt.Sprintf("could not encode json: %s\n", err), http.StatusInternalServerError)
	}
}

func OnGetHealth(w http.ResponseWriter, r *http.Request) {
	health := struct {
		Health bool `json:"health"`
	}{
		Health: true,
	}

	if err := json.NewEncoder(w).Encode(health); err != nil {
		http.Error(w, fmt.Sprintf("could not encode json: %s\n", err), http.StatusInternalServerError)
	}

}
