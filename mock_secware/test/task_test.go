package mock_secware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/ethereum/go-ethereum/crypto"
	"goplus/mock_secware/pkg/config"
	"goplus/mock_secware/pkg/handlers"
	"goplus/shared/pkg/signature"
	"goplus/shared/pkg/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_OnSecwareTask_Normal(t *testing.T) {
	expectResult := `{"safe":"yes"}`

	// 构造输入参数
	args := handlers.SecwareArgs{
		Sleep:  1,
		Crash:  false,
		Result: expectResult,
	}

	argsBytes, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	task := types.SecwareTask{
		SecwareId:      1,
		SecwareVersion: 2,
		SignedTx:       []byte{0xab, 0xcd},
		StartTime:      0x12345678,
		EndTime:        0x87654321,
		Args:           string(argsBytes),
	}

	// 签名
	skGateway, err := crypto.HexToECDSA(config.SkGatewayHex)
	if err != nil {
		t.Fatal(err)
	}

	signedTask, err := signature.SignSecwareTask(&task, skGateway)
	if err != nil {
		t.Fatal(err)
	}

	operator := bytes.Repeat([]byte{0x11}, 20)
	signedTask.Operator = operator

	jsonData, err := json.Marshal(signedTask)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	body := bytes.NewBuffer(jsonData)

	// 创建一个请求
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 注入 context
	ctx, cancelCtx := context.WithCancel(context.Background())
	newCtx, err := handlers.BuildContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	setting, ok := newCtx.Value("Setting").(*handlers.SecwareSetting)
	if !ok {
		t.Fatalf("could not get Setting")
	}

	req = req.WithContext(newCtx)

	// 创建一个 ResponseRecorder (它实现了 http.ResponseWriter) 来记录响应
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.OnSecwareTask)

	// 调用服务器的处理函数直接处理我们的请求
	handler.ServeHTTP(rr, req)

	// 检查状态码是否正确
	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 检查响应体是否正确
	var signedResult types.SignedSecwareResult
	if err := json.Unmarshal(rr.Body.Bytes(), &signedResult); err != nil {
		t.Fatalf("could not unmarshal result: %v", err)
	}

	// 检查签名
	if !signature.VerifySecwareSignature(&signedResult, setting.SkSecwareBytes) {
		t.Fatalf("invalid secware signature")
	}

	// 检查结果内容
	result := &signedResult.Result
	if result.Code != 0 || result.Message != "ok" {
		t.Fatalf("handler returned wrong result: %v", result)
	}

	if result.Details != expectResult {
		t.Fatalf("handler returned wrong result: %v, expect: %v", result.Details, expectResult)
	}

	//expected := `{"operator":"0x1111111111111111111111111111111111111111","task":{"signed_tx":"0xabcd","start_time":"0x12345678","end_time":"0x87654321","args":"0x7b226b6579223a2276616c7565227d"},"sig_gateway":"0x0ad860ce92cc0df623b2a6425cce46d6314525b98951ab2c547a1640bafd4dc83adf54243228acff3a6bc7177b45675b242cb5d3ef39da87330f36ee450d3a6601"}`
	//if rr.Body.String() != expected+"\n" {
	//	t.Errorf("handler returned unexpected body: got %v want %v",
	//		rr.Body.String(), expected)
	//}

	cancelCtx()
}

func Test_OnSecwareTask_BadSig(t *testing.T) {
	expectResult := `{"safe":"yes"}`

	// 构造输入参数
	args := handlers.SecwareArgs{
		Sleep:  1,
		Crash:  false,
		Result: expectResult,
	}

	argsBytes, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	task := types.SecwareTask{
		SecwareId:      1,
		SecwareVersion: 2,
		SignedTx:       []byte{0xab, 0xcd},
		StartTime:      0x12345678,
		EndTime:        0x87654321,
		Args:           string(argsBytes),
	}

	// 签名
	skGateway, err := crypto.HexToECDSA(config.SkGatewayHex)
	if err != nil {
		t.Fatal(err)
	}

	signedTask, err := signature.SignSecwareTask(&task, skGateway)
	if err != nil {
		t.Fatal(err)
	}

	// 篡改签名
	signedTask.SigGateway[10] += 1

	operator := bytes.Repeat([]byte{0x11}, 20)
	signedTask.Operator = operator

	jsonData, err := json.Marshal(signedTask)
	if err != nil {
		t.Fatalf("failed to marshal JSON: %v", err)
	}

	body := bytes.NewBuffer(jsonData)

	// 创建一个请求
	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 注入 context
	ctx, cancelCtx := context.WithCancel(context.Background())
	newCtx, err := handlers.BuildContext(ctx)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(newCtx)

	// 创建一个 ResponseRecorder (它实现了 http.ResponseWriter) 来记录响应
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.OnSecwareTask)

	// 调用服务器的处理函数直接处理我们的请求
	handler.ServeHTTP(rr, req)

	// 检查状态码是否正确
	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// 检查响应体是否正确
	expected := "signature verification failed\n"
	if rr.Body.String() != expected {
		t.Fatalf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}

	cancelCtx()
}
