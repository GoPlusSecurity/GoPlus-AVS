package tests

import (
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	mgr "goplus/avs/secwaremanager"
	"goplus/avs/secwaremanager/mocks"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

func newTestRunner() *mgr.DockerRunnerImpl {
	logger, _ := logging.NewZapLogger("development")

	runner := mgr.DockerRunnerImpl{
		Logger:              logger,
		ProjectNamePrefix:   "testsecware",
		PortProviderIntf:    &mocks.PortProvider{},
		CommandExecutorIntf: &mocks.CommandExecutor{},
		SecwareAccessorIntf: &mocks.SecwareAccessor{},
	}
	return &runner
}

func mockExecCommand(result string, exitCode int) *exec.Cmd {
	cmd := exec.Command(os.Args[0], "-test.run=TestExecCommandHelper")
	cmd.Env = append(cmd.Env, "EXEC_RESULT="+result)
	cmd.Env = append(cmd.Env, fmt.Sprintf("EXEC_CODE=%d", exitCode))
	return cmd
}

func TestExecCommandHelper(t *testing.T) {
	result := os.Getenv("EXEC_RESULT")
	execCodeStr := os.Getenv("EXEC_CODE")
	if execCodeStr == "" {
		os.Exit(0)
	}

	execCode, _ := strconv.Atoi(execCodeStr)
	if execCode != 0 {
		_, _ = fmt.Fprintf(os.Stderr, result)
		os.Exit(execCode)
	} else {
		_, _ = fmt.Fprintf(os.Stdout, result)
		os.Exit(0)
	}
}

func TestRunnerListing(t *testing.T) {
	runner := newTestRunner()

	mockCommandExecutor := runner.CommandExecutorIntf.(*mocks.CommandExecutor)
	mockSecwareAccessor := runner.SecwareAccessorIntf.(*mocks.SecwareAccessor)
	mockCommandExecutor.On("ExecCommand", "docker", "compose", "ls", "--format", "json").Return(mockExecCommand("[]", 0)).Once()

	res, err := runner.ListAvailableSecware()
	if err != nil {
		t.Errorf("ListAvailableSecware() error = %v", err)
		return
	}
	if len(res) != 0 {
		t.Errorf("ListAvailableSecware() = %v, want %v", res, 0)
	}

	mockState := mgr.SecwareStatus{
		SecwareId:          111,
		SecwareVersion:     222,
		Port:               7777,
		State:              "Running",
		ComposeProjectName: "testsecware-111-222-7777",
	}

	mockCommandExecutor.On("ExecCommand", "docker", "compose", "ls", "--format", "json").Return(mockExecCommand(`[{"Name":"testsecware-111-222-7777"}]`, 0)).Once()
	mockSecwareAccessor.On("GetSecwareMeta", &mockState).Return(mgr.SecwareMeta{SecwareId: 111, SecwareVersion: 222}, nil).Once()
	mockSecwareAccessor.On("GetSecwareHealth", &mockState).Return(mgr.SecwareHealth{Health: true}, nil).Once()
	res, err = runner.ListAvailableSecware()
	if err != nil {
		t.Errorf("ListAvailableSecware() error = %v", err)
		return
	}
	if len(res) != 1 {
		t.Errorf("ListAvailableSecware() = %v, want %v", res, 1)
	}

	if res[0].State != "Available" {
		t.Errorf("ListAvailableSecware() = %v, want %v", res[0].State, "Available")
	}
}

func TestComposeUp(t *testing.T) {
	runner := newTestRunner()

	mockPortProvider := runner.PortProviderIntf.(*mocks.PortProvider)
	mockCommandExecutor := runner.CommandExecutorIntf.(*mocks.CommandExecutor)
	mockSecwareAccessor := runner.SecwareAccessorIntf.(*mocks.SecwareAccessor)

	mockPort := 6789
	mockComposeFile := "testsecware-111-222.yml"
	mockState := mgr.SecwareStatus{
		SecwareId:          111,
		SecwareVersion:     222,
		Port:               mockPort,
		State:              "Running",
		ComposeProjectName: "testsecware-111-222-6789",
	}
	mockPortProvider.On("GetAvailablePort").Return(mockPort, nil).Once()

	mockCommandExecutor.On("ExecCommand", "docker", "compose", "-f", mockComposeFile, "pull").Return(mockExecCommand("", 0)).Once()
	mockCommandExecutor.On("ExecCommand", "docker", "compose", "-f", mockComposeFile, "up", "-d").Return(mockExecCommand("", 0)).Once()

	mockSecwareAccessor.On("GetSecwareMeta", &mockState).Return(mgr.SecwareMeta{SecwareId: 111, SecwareVersion: 222}, nil).Once()
	mockSecwareAccessor.On("GetSecwareHealth", &mockState).Return(mgr.SecwareHealth{Health: true}, nil).Once()

	state, err := runner.ComposeUp(111, 222, mockComposeFile)
	if err != nil {
		t.Errorf("ComposeUp() error = %v", err)
		return
	}
	if state.State != "Available" {
		t.Errorf("ComposeUp() = %v, want %v", state.State, "Available")
	}

	if state.Port != mockPort {
		t.Errorf("ComposeUp() = %v, want %v", state.Port, mockPort)
	}
}

func TestComposeDown(t *testing.T) {
	runner := newTestRunner()

	mockCommandExecutor := runner.CommandExecutorIntf.(*mocks.CommandExecutor)

	mockState := mgr.SecwareStatus{
		SecwareId:          111,
		SecwareVersion:     222,
		Port:               7777,
		State:              "Running",
		ComposeProjectName: "testsecware-111-222-7777",
	}

	mockCommandExecutor.On("ExecCommand", "docker", "compose", "-p", "testsecware-111-222-7777", "down").Return(mockExecCommand("", 0)).Once()

	state, err := runner.ComposeDown(&mockState)
	if err != nil {
		t.Errorf("ComposeDown() error = %v", err)
		return
	}
	if state.State != "Down" {
		t.Errorf("ComposeDown() = %v, want %v", state.State, "Down")
	}
}
