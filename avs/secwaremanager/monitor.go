// Package secwaremanager: SecwareMonitorImpl 用于监控各个 secware 的健康状态。
package secwaremanager

import (
	"context"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"goplus/avs/config"
	"goplus/avs/metrics"
	"time"
)

type SecwareMonitorInterface interface {
	Init() error
	Start(ctx context.Context) error
}

type SecwareMonitorImpl struct {
	logger      logging.Logger
	metricsIntf metrics.AvsMetricsInterface

	manager             *SecwareManager
	gatewayAccessorIntf GatewayAccessorInterface
	secwareAccessorIntf SecwareAccessorInterface
}

func NewSecwareMonitorImpl(cfg config.Config, manager *SecwareManager, gatewayAccessor *GatewayAccessorImpl, metricsIntf metrics.AvsMetricsInterface) (*SecwareMonitorImpl, error) {
	return &SecwareMonitorImpl{
		logger:              cfg.Logger,
		manager:             manager,
		gatewayAccessorIntf: gatewayAccessor,
		secwareAccessorIntf: &SecwareAccessorImpl{},
		metricsIntf:         metricsIntf,
	}, nil
}

func (m *SecwareMonitorImpl) Init() error {
	return nil
}

// Start 启动 monitor
func (m *SecwareMonitorImpl) Start(ctx context.Context) error {
	m.logger.Info("SecwareMonitorImpl Start")

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		healthyResult := m.checkRunningSecware()
		m.reportToGateway(healthyResult)
		select {
		case <-ctx.Done():
			m.logger.Info("SecwareMonitorImpl Exit")
			return nil
		case <-ticker.C:
		}
	}
}

type SecwareHealthResult struct {
	SecwareId      int  `json:"id"`
	SecwareVersion int  `json:"version"`
	Health         bool `json:"health"`
}

// checkRunningSecware 获取所有 secware 的健康情况，返回汇总的结果。
func (m *SecwareMonitorImpl) checkRunningSecware() []SecwareHealthResult {
	availableSecwares := m.manager.GetAvailableSecwareList()
	healthyResult := make([]SecwareHealthResult, len(availableSecwares))

	// 遍历每一个 secware，查询健康情况
	for idx, secware := range availableSecwares {
		health, err := m.secwareAccessorIntf.GetSecwareHealth(secware)

		healthy := false
		if err == nil && health.IsHealthy() {
			healthy = true
		}

		result := SecwareHealthResult{
			SecwareId:      secware.SecwareId,
			SecwareVersion: secware.SecwareVersion,
			Health:         healthy,
		}
		healthyResult[idx] = result
	}

	return healthyResult
}

// reportToGateway 把监控情况汇报给 gateway
func (m *SecwareMonitorImpl) reportToGateway(results []SecwareHealthResult) {
	for _, result := range results {
		m.logger.Infof("Secware %d-%d is healthy: %t", result.SecwareId, result.SecwareVersion, result.Health)
	}

	m.metricsIntf.SetSecwareNum(len(results))

	err := m.gatewayAccessorIntf.ReportHealth(results)
	if err != nil {
		m.logger.Error("Failed to report health to gateway", err)
	} else {
		m.metricsIntf.IncHealthReported()
		m.logger.Info("Reported health to gateway successfully")
	}
}
