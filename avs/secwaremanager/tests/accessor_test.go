package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goplus/avs/secwaremanager"
	"goplus/shared/pkg/types"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newSecwareAccessorImpl() *secwaremanager.SecwareAccessorImpl {
	return &secwaremanager.SecwareAccessorImpl{}
}

func TestSecwareAccessorImpl_HandleTask(t *testing.T) {
	operator, _ := types.NewHexBytesFromString("0x1111111111111111111111111111111111111111")
	signTx, _ := types.NewHexBytesFromString("0xf86e01851191460ee38252089497e542ec6b81dea28f212775ce8ac436ab77a7df880de0b6b3a764000082307826a0963edb6e57c2c6bd0d4a8d827a53f6f9e164f09dd1bc5d1f8580c020abad56b5a04343639754d9d9662e1bb6cdeee65f5c315d27fe183119763c9c235876a3d2f8")
	sigGateway, _ := types.NewHexBytesFromString("0x1c604a17c286707f35631553d9b98344ba2a277d57d8b1ca73e8fb9ab3a2b9117e3250f77a26d170d672631082ae02f63c271afd259fff8b071f49d2ebaca9ac00")
	sigSecware, _ := types.NewHexBytesFromString("0x396d36bfcb464a624fdddea887400fbd10e47f2ed9daa2b9a76d669c3a047dae")
	signTask := types.SignedSecwareTask{
		Operator: operator,
		Task: types.SecwareTask{
			SecwareId:      1,
			SecwareVersion: 1,
			SignedTx:       signTx,
			StartTime:      0,
			EndTime:        0,
			Args:           `{\"args\": \"TEST\"}`,
		},
		SigGateway: sigGateway,
	}
	signTaskBody, _ := json.Marshal(signTask)

	secwareResult := types.SignedSecwareResult{
		Result: types.SecwareResult{
			Code:           0,
			Message:        "ok",
			Details:        `{"test": "ok"}`,
			Operator:       operator,
			SecwareId:      1,
			SecwareVersion: 1,
		},
		SigSecware: sigSecware,
	}
	secwareResultBody, _ := json.Marshal(secwareResult)

	t.Log(string(signTaskBody))
	testSecware := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		requestBody := make([]byte, r.ContentLength)
		r.Body.Read(requestBody)
		if !bytes.Equal(requestBody, signTaskBody) {
			t.Errorf("request body not match")
			return
		}

		fmt.Fprintln(w, string(secwareResultBody))
	}))

	defer testSecware.Close()
	testServerPort := testSecware.Listener.Addr().(*net.TCPAddr).Port

	mockState := secwaremanager.SecwareStatus{
		SecwareId:          1,
		SecwareVersion:     1,
		Port:               testServerPort,
		State:              "Available",
		ComposeProjectName: fmt.Sprintf("testsecware-%d-%d-%d", 1, 1, testServerPort),
	}

	sa := newSecwareAccessorImpl()
	result, err := sa.HandleTask(&mockState, &signTask)
	if err != nil {
		t.Errorf("HandleTask() error = %v", err)
		return
	}
	t.Log(result)
	if !bytes.Equal(result.Result.Operator, operator) {
		t.Errorf("HandleTask() = %v, want %v", result.Result.Operator, operator)
	}

	if !bytes.Equal(result.SigSecware, sigSecware) {
		t.Errorf("HandleTask() = %v, want %v", result.SigSecware, sigSecware)
	}
}
