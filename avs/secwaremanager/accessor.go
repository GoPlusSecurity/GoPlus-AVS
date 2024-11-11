// Accessor 是对所有 Secware 接口的统一封装
package secwaremanager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"goplus/avs/config"
	"goplus/shared/pkg/signature"
	"goplus/shared/pkg/types"
	"io"
	"net/http"
	"time"
)

// SecwareAccessorInterface 对Secware的统一交互界面
type SecwareAccessorInterface interface {
	GetSecwareMeta(*SecwareStatus) (SecwareMeta, error)
	GetSecwareHealth(*SecwareStatus) (SecwareHealth, error)
	HandleTask(*SecwareStatus, *types.SignedSecwareTask) (types.SignedSecwareResult, error)
}

type SecwareAccessorImpl struct {
}

type SecwareMeta struct {
	SecwareId      int
	SecwareVersion int
}

type SecwareRequest struct {
	Id      string        `json:"id"`
	Version string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

// SecwareHealth 对 Secware 健康状况的描述
type SecwareHealth struct {
	Health bool
}

func (h *SecwareHealth) IsHealthy() bool {
	return h.Health
}

type secwareMetaResponse struct {
	SecwareId      int `json:"secware_id"`
	SecwareVersion int `json:"secware_version"`
}

// GetSecwareMeta 获取 secware 的描述信息
func (s *SecwareAccessorImpl) GetSecwareMeta(state *SecwareStatus) (SecwareMeta, error) {

	host := "localhost"
	url := fmt.Sprintf("http://%s:%d/meta", host, state.Port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SecwareMeta{}, err
	}

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return SecwareMeta{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SecwareMeta{}, err
	}

	var metaResp secwareMetaResponse
	err = json.Unmarshal(body, &metaResp)
	if err != nil {
		return SecwareMeta{}, err
	}

	return SecwareMeta{
		SecwareId:      metaResp.SecwareId,
		SecwareVersion: metaResp.SecwareVersion,
	}, nil
}

type secwareHealthResponse struct {
	Health bool `json:"health"`
}

// GetSecwareHealth 获取 Secware 的健康状况
func (s *SecwareAccessorImpl) GetSecwareHealth(state *SecwareStatus) (SecwareHealth, error) {
	host := "localhost"
	url := fmt.Sprintf("http://%s:%d/health", host, state.Port)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return SecwareHealth{}, err
	}

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return SecwareHealth{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SecwareHealth{}, err
	}

	var healthResp secwareHealthResponse
	err = json.Unmarshal(body, &healthResp)
	if err != nil {
		return SecwareHealth{}, err
	}

	return SecwareHealth{Health: healthResp.Health}, nil
}

// HandleTask 处理具体的交易安全检测任务
func (s *SecwareAccessorImpl) HandleTask(state *SecwareStatus, signTask *types.SignedSecwareTask) (types.SignedSecwareResult, error) {
	data, err := json.Marshal(signTask)
	if err != nil {
		return types.SignedSecwareResult{}, err
	}

	host := "localhost"
	url := fmt.Sprintf("http://%s:%d/secware", host, state.Port)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return types.SignedSecwareResult{}, err
	}

	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)

	if err != nil {
		return types.SignedSecwareResult{}, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.SignedSecwareResult{}, err
	}

	result := types.SignedSecwareResult{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return types.SignedSecwareResult{}, err
	}

	return result, nil
}

// GatewayAccessorInterface 对Gateway的统一交互界面
type GatewayAccessorInterface interface {
	GetSecwareConfig() ([]SecwareConfig, error)
	ReportHealth([]SecwareHealthResult) error
}

type GatewayAccessorImpl struct {
	AddressOperator common.Address
	BLSKeypair      *bls.KeyPair
	GatewayUrl      string
}

func NewGatewayAccessorImpl(cfg config.Config) (*GatewayAccessorImpl, error) {
	return &GatewayAccessorImpl{
		AddressOperator: cfg.AddressOperator,
		BLSKeypair:      cfg.BLSKeypair,
		GatewayUrl:      cfg.GatewayUrl,
	}, nil
}

type secwareConfigRequest struct {
	Time     int64          `json:"time"`
	Operator types.HexBytes `json:"operator"`
}

type signedSecwareConfigRequest struct {
	Args        secwareConfigRequest `json:"args"`
	SigOperator types.HexBytes       `json:"sig_operator"`
}

type GatewayResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

func (g *GatewayAccessorImpl) GetSecwareConfig() ([]SecwareConfig, error) {
	sc := secwareConfigRequest{
		Time:     time.Now().Unix(),
		Operator: g.AddressOperator[:],
	}

	signedSc, err := signSecwareConfigRequest(&sc, g.BLSKeypair)
	if err != nil {
		return nil, err
	}

	signedScBytes, err := json.Marshal(signedSc)
	if err != nil {
		return nil, err
	}

	var body []byte
	err = retry.Do(func() error {
		configUrl := fmt.Sprintf("%s/pull/operator-config", g.GatewayUrl)
		req, err := http.NewRequest("POST", configUrl, bytes.NewBuffer(signedScBytes))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		client := http.Client{Timeout: time.Second * 10}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return nil
	}, retry.Attempts(3), retry.Delay(1*time.Second))

	if err != nil {
		return nil, err
	}

	var configResp GatewayResponse[map[string][]SecwareConfig]
	err = json.Unmarshal(body, &configResp)
	if err != nil {
		return nil, err
	}

	if configResp.Code != 200 {
		return nil, fmt.Errorf("secware config bad response code %d", configResp.Code)
	}

	result, ok := configResp.Result["secware"]
	if !ok {
		return nil, fmt.Errorf("no secware in response")
	}

	return result, nil
}

func signSecwareConfigRequest(req *secwareConfigRequest, blsKeyPair *bls.KeyPair) (*signedSecwareConfigRequest, error) {
	hashReq, err := signature.HashJSON(req)
	if err != nil {
		return nil, err
	}

	sigReq := blsKeyPair.SignMessage(hashReq)

	return &signedSecwareConfigRequest{
		Args:        *req,
		SigOperator: sigReq.Marshal(),
	}, nil
}

type ReportHealthRequest struct {
	Time     int64                 `json:"time"`
	Operator types.HexBytes        `json:"operator"`
	Secware  []SecwareHealthResult `json:"secware"`
}

type SignedReportHealthRequest struct {
	Heartbeat   ReportHealthRequest `json:"heartbeat"`
	SigOperator types.HexBytes      `json:"sig_operator"`
}

func (g *GatewayAccessorImpl) ReportHealth(health []SecwareHealthResult) error {
	rh := ReportHealthRequest{
		Time:     time.Now().Unix(),
		Operator: g.AddressOperator[:],
		Secware:  health,
	}

	srh, err := signReportHealthRequest(&rh, g.BLSKeypair)
	if err != nil {
		return err
	}

	srhBytes, err := json.Marshal(srh)
	if err != nil {
		return err
	}

	healthUrl := fmt.Sprintf("%s/report/heartbeat", g.GatewayUrl)
	req, err := http.NewRequest("POST", healthUrl, bytes.NewBuffer(srhBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var respBody GatewayResponse[interface{}]
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return err
	}

	if respBody.Code != 200 {
		return fmt.Errorf("report health bad response code %d", respBody.Code)
	}

	return nil

}

func signReportHealthRequest(req *ReportHealthRequest, blsKeyPair *bls.KeyPair) (*SignedReportHealthRequest, error) {
	hashReq, err := signature.HashJSON(req)
	if err != nil {
		return nil, err
	}

	sigReq := blsKeyPair.SignMessage(hashReq)
	return &SignedReportHealthRequest{
		Heartbeat:   *req,
		SigOperator: sigReq.Marshal(),
	}, nil
}
