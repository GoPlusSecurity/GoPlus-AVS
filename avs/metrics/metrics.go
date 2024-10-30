package metrics

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"goplus/avs/config"
)

type AvsMetricsInterface interface {
	IncTaskHandled()
	SetTaskDuration(float64)
	IncTaskFailed()
	IncHealthReported()
	IncSecwareSynced()
	SetSecwareNum(int)
}

const (
	PromNamespace = "avs_operator"
)

type AvsMetrics struct {
	addressOperator     common.Address
	addressOperatorStr  string
	numTaskHandled      *prometheus.CounterVec
	taskHandledDuration *prometheus.SummaryVec
	numTaskFailed       *prometheus.CounterVec
	numHealthReported   *prometheus.CounterVec
	numSecwareSynced    *prometheus.CounterVec
	secwareNum          *prometheus.GaugeVec
}

func NewAvsMetrics(cfg config.Config) (*AvsMetrics, error) {
	return &AvsMetrics{
		addressOperator:    cfg.AddressOperator,
		addressOperatorStr: cfg.AddressOperator.String(),
		numTaskHandled: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: PromNamespace,
				Name:      "num_task_handled",
				Help:      "The number of tasks handled by the avs operator",
			}, []string{"operator_pk"}),
		taskHandledDuration: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace:  PromNamespace,
				Name:       "task_handled_duration",
				Help:       "The duration of tasks handled by the avs operator",
				Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.95: 0.01, 0.99: 0.001},
			}, []string{"operator_pk"}),
		numTaskFailed: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: PromNamespace,
				Name:      "num_task_failed",
				Help:      "The number of tasks failed by the avs operator",
			}, []string{"operator_pk"}),
		numHealthReported: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: PromNamespace,
				Name:      "num_health_reported",
				Help:      "The number of health reports reported by the avs operator",
			}, []string{"operator_pk"}),
		numSecwareSynced: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: PromNamespace,
				Name:      "num_secware_synced",
				Help:      "The number of secwares synced by the avs operator",
			}, []string{"operator_pk"}),
		secwareNum: promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: PromNamespace,
				Name:      "secware_num",
				Help:      "The number of secwares managed by the avs operator",
			}, []string{"operator_pk"}),
	}, nil
}

func (m *AvsMetrics) IncTaskHandled() {
	m.numTaskHandled.WithLabelValues(m.addressOperatorStr).Inc()
}

func (m *AvsMetrics) SetTaskDuration(duration float64) {
	m.taskHandledDuration.WithLabelValues(m.addressOperatorStr).Observe(duration)
}

func (m *AvsMetrics) IncTaskFailed() {
	m.numTaskFailed.WithLabelValues(m.addressOperatorStr).Inc()
}

func (m *AvsMetrics) IncHealthReported() {
	m.numHealthReported.WithLabelValues(m.addressOperatorStr).Inc()
}

func (m *AvsMetrics) IncSecwareSynced() {
	m.numSecwareSynced.WithLabelValues(m.addressOperatorStr).Inc()
}

func (m *AvsMetrics) SetSecwareNum(num int) {
	m.secwareNum.WithLabelValues(m.addressOperatorStr).Set(float64(num))
}

var _ AvsMetricsInterface = (*AvsMetrics)(nil)
