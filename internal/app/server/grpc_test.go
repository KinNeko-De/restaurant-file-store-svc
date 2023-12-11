package server

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartGrpcServer_ProcessAlreadyListenToPort_AppCrash(t *testing.T) {
	if os.Getenv("EXECUTE") == "1" {
		StartGrpcServer(make(chan struct{}), make(chan struct{}), ":3110")
		return
	}

	blockingcmd := exec.Command(os.Args[0], "-test.run=TestStartGrpcServer_ProcessAlreadyListenToPort_AppCrash")
	blockingcmd.Env = append(os.Environ(), "EXECUTE=1")
	blockingErr := blockingcmd.Start()

	require.Nil(t, blockingErr)
	defer blockingcmd.Process.Kill()

	time.Sleep(1 * time.Second) // give the service some time to start

	cmd := exec.Command(os.Args[0], "-test.run=TestStartGrpcServer_ProcessAlreadyListenToPort_AppCrash")
	cmd.Env = append(os.Environ(), "EXECUTE=1")
	err := cmd.Run()
	require.NotNil(t, err)
	exitCode := err.(*exec.ExitError).ExitCode()
	assert.Equal(t, 50, exitCode)
}

func TestStartGrpcServer_PortMalformed(t *testing.T) {
	if os.Getenv("EXECUTE") == "1" {
		StartGrpcServer(make(chan struct{}), make(chan struct{}), "malformedPort")
		return
	}

	runningApp := exec.Command(os.Args[0], "-test.run=TestStartGrpcServer_PortMalformed")
	runningApp.Env = append(os.Environ(), "EXECUTE=1")
	err := runningApp.Run()
	require.NotNil(t, err)
	exitCode := err.(*exec.ExitError).ExitCode()
	assert.Equal(t, 50, exitCode)
}
