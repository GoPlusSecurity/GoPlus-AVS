package secwaremanager

import (
	"context"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/avast/retry-go/v4"
	"github.com/ethereum/go-ethereum/common"
	"goplus/avs/config"
	"goplus/avs/metrics"
	"io"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

// SecwareManager 是 AVS 后台的组织者，用于管理各个内部组件。 包括 DockerRunner 和 SecwareMonitorImpl
type SecwareManager struct {
	logger             logging.Logger
	metricsIntf        metrics.AvsMetricsInterface
	composeFileDirPath string
	addressOperator    common.Address

	availableSecwares    []*SecwareStatus
	availableSecwaresMap map[string]*SecwareStatus
	rwLock               sync.RWMutex

	GatewayAccessorIntf GatewayAccessorInterface
	DockerRunnerIntf    DockerRunnerInterface
	SecwareMonitorIntf  SecwareMonitorInterface
}

type SecwareConfig struct {
	SecwareId      int    `json:"id"`
	SecwareVersion int    `json:"version"`
	ComposeFileUrl string `json:"docker_compose_file"`
	Name           string `json:"name"`
}

func NewSecwareManager(cfg config.Config, metrics metrics.AvsMetricsInterface) (*SecwareManager, error) {
	manager := &SecwareManager{
		logger:             cfg.Logger,
		metricsIntf:        metrics,
		composeFileDirPath: cfg.ComposeFilePath,
		addressOperator:    cfg.AddressOperator,

		availableSecwares:    make([]*SecwareStatus, 0),
		availableSecwaresMap: make(map[string]*SecwareStatus),

		DockerRunnerIntf:    nil,
		SecwareMonitorIntf:  nil,
		GatewayAccessorIntf: nil,
	}

	gatewayAccessor, err := NewGatewayAccessorImpl(cfg)
	if err != nil {
		return nil, err
	}
	runner, err := NewDockerRunnerImpl(cfg)
	if err != nil {
		return nil, err
	}
	monitor, err := NewSecwareMonitorImpl(cfg, manager, gatewayAccessor, metrics)
	if err != nil {
		return nil, err
	}

	manager.DockerRunnerIntf = runner
	manager.SecwareMonitorIntf = monitor
	manager.GatewayAccessorIntf = gatewayAccessor

	return manager, nil
}

func (mgr *SecwareManager) Init() error {
	mgr.logger.Info("Starting SecwareManager Init...")
	err := mgr.DockerRunnerIntf.CheckDockerCompose()
	if err != nil {
		return err
	}

	err = mgr.checkComposeFileDirPath()
	if err != nil {
		return err
	}

	errSecware, err := mgr.syncSecware()
	if err != nil {
		return err
	}

	if len(errSecware) > 0 {
		return errors.New("some secware failed to start")
	}

	mgr.logger.Info("SecwareManager Init Success")
	return nil
}

func (mgr *SecwareManager) Start(ctx context.Context) error {
	mgr.logger.Info("SecwareManager Start")

	monitorDone := make(chan struct{})
	go func() {
		if err := mgr.SecwareMonitorIntf.Start(ctx); err != nil {
			mgr.logger.Error(err.Error())
			close(monitorDone)
		}
	}()

	// 启动定时器，定期同步各个 secware 的 docker compose
	ticker := time.NewTicker(300 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			mgr.logger.Info("SecwareManager Exit")
			return nil
		case <-monitorDone:
			return errors.New("SecwareMonitorImpl is exiting with error")
		case <-ticker.C:
			if errSecware, err := mgr.syncSecware(); err != nil {
				mgr.logger.Error(err.Error())
				for _, i := range errSecware {
					mgr.logger.Errorf("secware %d-%d failed to start", i.SecwareId, i.SecwareVersion)
				}
			}
		}
	}
}

func (mgr *SecwareManager) checkComposeFileDirPath() error {
	_, err := os.Stat(mgr.composeFileDirPath)
	if os.IsNotExist(err) {
		return errors.New("compose file path not exist")
	}
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", mgr.composeFileDirPath, mgr.addressOperator.String())
	err = os.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

// GetAvailableSecwareList 获取当前可用的 secware 列表
func (mgr *SecwareManager) GetAvailableSecwareList() []*SecwareStatus {
	mgr.rwLock.RLock()
	defer mgr.rwLock.RUnlock()
	runningSecwares := make([]*SecwareStatus, len(mgr.availableSecwares))
	copy(runningSecwares, mgr.availableSecwares)
	return runningSecwares
}

func (mgr *SecwareManager) downloadComposeFile(secwareId int, secwareVersion int, composeFileUrl string) (string, error) {
	var body []byte
	err := retry.Do(func() error {
		req, err := http.NewRequest("GET", composeFileUrl, nil)
		if err != nil {
			return err
		}

		client := http.Client{Timeout: 10 * time.Second}
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

	}, retry.Attempts(3), retry.Delay(2*time.Second))

	path := fmt.Sprintf("%s/%s/secware-%d-%d.yml", mgr.composeFileDirPath, mgr.addressOperator.String(), secwareId, secwareVersion)
	err = os.WriteFile(path, body, 0644)
	if err != nil {
		return "", err
	}

	return path, nil
}

// syncSecware 用于同步 secware docker compose 的设定
func (mgr *SecwareManager) syncSecware() ([]SecwareConfig, error) {
	secwareCfg, err := mgr.GatewayAccessorIntf.GetSecwareConfig()
	mgr.logger.Infof("secware config length: %d", len(secwareCfg))

	if err != nil {
		return []SecwareConfig{}, err
	}

	secwareState, err := mgr.DockerRunnerIntf.ListAvailableSecware()
	if err != nil {
		return []SecwareConfig{}, err
	}

	secwareState = mgr.downDuplicateSecware(secwareState)

	secwareStateMap := make(map[string]*SecwareStatus)
	secwareConfigMap := make(map[string]SecwareConfig)

	for _, i := range secwareState {
		name := fmt.Sprintf("%d-%d", i.SecwareId, i.SecwareVersion)
		secwareStateMap[name] = i
	}

	for _, i := range secwareCfg {
		name := fmt.Sprintf("%d-%d", i.SecwareId, i.SecwareVersion)
		secwareConfigMap[name] = i
	}

	var availableSecwares []*SecwareStatus
	var errSecwares []SecwareConfig
	// 遍历所有需要使用的 secware，把还没启动的启动，已经启动的就不管
	for name, i := range secwareConfigMap {
		s, ok := secwareStateMap[name]
		isAvailable := false

		// 如果所需的 secware 已经存在而且运转正常，就跳过
		if ok && s.State == StateAvailable {
			isAvailable = true
		}

		if isAvailable {
			availableSecwares = append(availableSecwares, s)
			mgr.logger.Infof("secware %d-%d is available", i.SecwareId, i.SecwareVersion)
			continue
		}

		// 如果所需的 secware 已经存在但执行状态不正常，则先关闭。
		if s != nil {
			_, _ = mgr.DockerRunnerIntf.ComposeDown(s)
		}

		composeFilePath, err := mgr.downloadComposeFile(i.SecwareId, i.SecwareVersion, i.ComposeFileUrl)
		if err != nil {
			mgr.logger.Errorf("secware %d-%d download compose file failed. reason: %s", i.SecwareId, i.SecwareVersion, err.Error())
			continue
		}

		newState, err := mgr.DockerRunnerIntf.ComposeUp(i.SecwareId, i.SecwareVersion, composeFilePath)
		if err != nil {
			errSecwares = append(errSecwares, i)
			mgr.logger.Errorf("secware %d-%d failed to start. reason: %s", i.SecwareId, i.SecwareVersion, err.Error())
			continue
		}
		if newState.State == StateAvailable {
			availableSecwares = append(availableSecwares, newState)
		}
	}

	// 把已经废弃的 secware 关停
	for name, i := range secwareStateMap {
		if _, ok := secwareConfigMap[name]; !ok {
			_, _ = mgr.DockerRunnerIntf.ComposeDown(i)
		}
	}

	availableSecwaresMap := make(map[string]*SecwareStatus)
	for _, i := range availableSecwares {
		name := fmt.Sprintf("%d-%d", i.SecwareId, i.SecwareVersion)
		availableSecwaresMap[name] = i
	}

	mgr.rwLock.Lock()
	defer mgr.rwLock.Unlock()
	mgr.availableSecwaresMap = availableSecwaresMap
	mgr.availableSecwares = availableSecwares

	mgr.metricsIntf.IncSecwareSynced()
	return errSecwares, nil

}

// downDuplicateSecware 用于关闭重复的 secware
func (mgr *SecwareManager) downDuplicateSecware(statusList []*SecwareStatus) []*SecwareStatus {
	copyStatusList := make([]*SecwareStatus, len(statusList))
	copy(copyStatusList, statusList)

	sort.Slice(copyStatusList, func(i, j int) bool {
		return copyStatusList[i].Port < copyStatusList[j].Port
	})

	existsMap := make(map[int]bool)
	var newStatusList []*SecwareStatus
	for _, i := range copyStatusList {
		if _, ok := existsMap[i.SecwareId]; ok {
			_, _ = mgr.DockerRunnerIntf.ComposeDown(i)
			mgr.logger.Infof("secware %d-%d is duplicated, drop it", i.SecwareId, i.SecwareVersion)
			continue
		}
		existsMap[i.SecwareId] = true
		newStatusList = append(newStatusList, i)
	}

	return newStatusList
}

func (mgr *SecwareManager) GetSecwareState(secwareId int, secwareVersion int) (*SecwareStatus, error) {
	mgr.rwLock.RLock()
	defer mgr.rwLock.RUnlock()
	name := fmt.Sprintf("%d-%d", secwareId, secwareVersion)
	if s, ok := mgr.availableSecwaresMap[name]; ok {
		return s, nil
	}
	return nil, errors.New("secware not found")
}
