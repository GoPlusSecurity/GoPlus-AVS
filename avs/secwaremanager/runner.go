package secwaremanager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"goplus/avs/config"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	StateUnknown   = "Unknown"
	StateRunning   = "Running"
	StateDown      = "Down"
	StateMetaError = "MetaError"
	StateAvailable = "Available"
	StateUnhealthy = "Unhealthy"
)

type PortProviderInterface interface {
	GetAvailablePort() (int, error)
}

type PortProviderImpl struct{}

func (p *PortProviderImpl) GetAvailablePort() (int, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	return addr.Port, nil
}

type CommandExecutorInterface interface {
	ExecCommand(name string, arg ...string) *exec.Cmd
}

type CommandExecutorImpl struct{}

func (c *CommandExecutorImpl) ExecCommand(name string, arg ...string) *exec.Cmd {
	return exec.Command(name, arg...)
}

type SecwareStatus struct {
	SecwareId          int
	SecwareVersion     int
	Port               int
	State              string
	ComposeProjectName string
}

type DockerRunnerInterface interface {
	CheckDockerCompose() error
	ListAvailableSecware() ([]*SecwareStatus, error)
	ComposeUp(id int, version int, composeFilePath string) (*SecwareStatus, error)
	ComposeDown(status *SecwareStatus) (*SecwareStatus, error)
}

// DockerRunnerImpl 处理各个 secware 对应的 docker compose 的启停
type DockerRunnerImpl struct {
	Logger            logging.Logger
	ProjectNamePrefix string

	PortProviderIntf    PortProviderInterface
	CommandExecutorIntf CommandExecutorInterface
	SecwareAccessorIntf SecwareAccessorInterface
}

func NewDockerRunnerImpl(cfg config.Config) (*DockerRunnerImpl, error) {
	return &DockerRunnerImpl{
		Logger:              cfg.Logger,
		ProjectNamePrefix:   "secware",
		PortProviderIntf:    &PortProviderImpl{},
		CommandExecutorIntf: &CommandExecutorImpl{},
		SecwareAccessorIntf: &SecwareAccessorImpl{},
	}, nil
}

func (d *DockerRunnerImpl) CheckDockerCompose() error {
	cmd := exec.Command("docker", "compose", "version")
	_, err := cmd.Output()
	if err != nil {
		return errors.New("docker compose not found")
	}

	return nil
}

// secware project 的名字有特定的规范，使得其不仅仅是标识，还可以被解析为单独的字段
// <project_name>-<id>_<version>_<http_port>
func (d *DockerRunnerImpl) getProjectName(id int, version int, port int) string {
	name := fmt.Sprintf("%s-%d-%d-%d", d.ProjectNamePrefix, id, version, port)
	return name
}

func (d *DockerRunnerImpl) parseProjectName(name string) (*SecwareStatus, error) {
	if !strings.HasPrefix(name, d.ProjectNamePrefix) {
		return nil, errors.New("invalid project name")
	}

	parts := strings.Split(name, "-")
	if len(parts) != 4 {
		return nil, errors.New("invalid project name")
	}

	checkerId, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	checkerVersion, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, err
	}

	state := SecwareStatus{
		SecwareId:          checkerId,
		SecwareVersion:     checkerVersion,
		Port:               port,
		State:              StateRunning,
		ComposeProjectName: name,
	}

	return &state, nil
}

// checkSecwareAvailable 通过调用 secware 的一个简单接口来检测其是否可用
func (d *DockerRunnerImpl) checkSecwareAvailable(status *SecwareStatus) string {
	meta, err := d.SecwareAccessorIntf.GetSecwareMeta(status)
	if err != nil {
		return StateUnknown
	}
	if meta.SecwareId != status.SecwareId || meta.SecwareVersion != status.SecwareVersion {
		return StateMetaError
	}
	health, err := d.SecwareAccessorIntf.GetSecwareHealth(status)
	if err != nil {
		return StateUnknown
	}
	if !health.IsHealthy() {
		return StateUnhealthy
	}
	return StateAvailable
}

// ListSecware 使用 docker compose 命令获取当前所有运行的Secware, 并检查其状态
func (d *DockerRunnerImpl) ListAvailableSecware() ([]*SecwareStatus, error) {

	cmd := d.CommandExecutorIntf.ExecCommand("docker", "compose", "ls", "--format", "json")
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	type ComposeStatus struct {
		Name string `json:"Name"`
	}

	var dockerComposeStatus []ComposeStatus
	err = json.Unmarshal(out, &dockerComposeStatus)
	if err != nil {
		return nil, err
	}

	var secwareStatus []*SecwareStatus
	for _, status := range dockerComposeStatus {
		state, err := d.parseProjectName(status.Name)
		if err != nil {
			continue
		}
		secwareStatus = append(secwareStatus, state)
	}

	for _, s := range secwareStatus {
		state := d.checkSecwareAvailable(s)
		if state != StateAvailable {
			_, _ = d.ComposeDown(s)
		} else {
			s.State = state
		}
	}

	return secwareStatus, nil
}

// ComposeUp 启动Secware, 并等待其可用
func (d *DockerRunnerImpl) ComposeUp(id int, version int, composeFilePath string) (*SecwareStatus, error) {
	cmd := d.CommandExecutorIntf.ExecCommand("docker", "compose", "-f", composeFilePath, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	port, err := d.PortProviderIntf.GetAvailablePort()
	if err != nil {
		return nil, err
	}
	projectName := d.getProjectName(id, version, port)

	cmd = d.CommandExecutorIntf.ExecCommand("docker", "compose", "-f", composeFilePath, "up", "-d")
	cmd.Env = append(cmd.Env, fmt.Sprintf("SECWARE_PORT=%s:%d", "127.0.0.1", port))
	cmd.Env = append(cmd.Env, fmt.Sprintf("COMPOSE_PROJECT_NAME=%s", projectName))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	state := &SecwareStatus{
		SecwareId:          id,
		SecwareVersion:     version,
		Port:               port,
		State:              StateRunning,
		ComposeProjectName: projectName,
	}

	// 等待服务完全启动，进入待命状态
	state.State = d.waitForStabled(state)
	d.Logger.Info(fmt.Sprintf("Secware %d-%d Up Port:%d", id, version, port))
	return state, nil
}

func (d *DockerRunnerImpl) waitForStabled(status *SecwareStatus) string {
	timeout := 10
	var s string
	for i := 0; i < timeout; i++ {
		s = d.checkSecwareAvailable(status)
		if s == StateAvailable {
			return s
		}
		time.Sleep(1 * time.Second)
	}
	return s
}

// ComposeDown 关闭Secware, 不检查其是否存在直接关闭
func (d *DockerRunnerImpl) ComposeDown(status *SecwareStatus) (*SecwareStatus, error) {
	cmd := d.CommandExecutorIntf.ExecCommand("docker", "compose", "-p", status.ComposeProjectName, "down")
	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	status.State = StateDown
	d.Logger.Info(fmt.Sprintf("Secware %d-%d Down Port:%d", status.SecwareId, status.SecwareVersion, status.Port))
	return status, nil
}
