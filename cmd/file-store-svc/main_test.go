//go:build component

package main

import (
	"context"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/kinneko-de/restaurant-file-store-svc/test/testing/envvar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthV1 "google.golang.org/grpc/health/grpc_health_v1"
)

func TestMain_MetricConfigIsMissing(t *testing.T) {
	if os.Getenv("EXECUTE") == "1" {
		main()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestMain_MetricConfigIsMissing")
	cmd.Env = append(os.Environ(), "EXECUTE=1")
	err := cmd.Run()
	require.NotNil(t, err)
	exitCode := err.(*exec.ExitError).ExitCode()
	assert.Equal(t, 40, exitCode)
}

// test does not run on windows
// In case you broke something, the test will run forever
// In the pipeline you will see:
// panic: test timed out after 5m0s
// running tests:
// TestMain_ApplicationListenToInterrupt_GracefullShutdown (5m0s)
func TestMain_ApplicationListenToSIGTERM_AndGracefullyShutdown(t *testing.T) {
	if os.Getenv("EXECUTE") == "1" {
		main()
		return
	}

	envvar.SetAllNeceassaryEnvironemntVariables(t)
	cmd := exec.Command(os.Args[0], "-test.run=TestMain_ApplicationListenToSIGTERM_AndGracefullyShutdown")
	cmd.Env = append(os.Environ(), "EXECUTE=1")
	err := cmd.Start()
	require.Nil(t, err)
	time.Sleep(1 * time.Second)
	cmd.Process.Signal(syscall.SIGTERM)
	err = cmd.Wait()
	require.Nil(t, err)
	exitCode := cmd.ProcessState.ExitCode()
	assert.Equal(t, 0, exitCode)
}

func TestMain_HealthCheckIsServing_Liveness(t *testing.T) {
	serviceToCheck := "liveness"

	if os.Getenv("EXECUTE") == "1" {
		main()
		return
	}

	envvar.SetAllNeceassaryEnvironemntVariables(t)
	runningApp := exec.Command(os.Args[0], "-test.run=TestMain_HealthCheckIsServing_Liveness")
	runningApp.Env = append(os.Environ(), "EXECUTE=1")
	blockingErr := runningApp.Start()
	require.Nil(t, blockingErr)
	defer runningApp.Process.Kill()

	expectedStatus := healthV1.HealthCheckResponse_SERVING
	healthResponse, err := waitForStatus(t, serviceToCheck, expectedStatus, 100*time.Millisecond, 100)

	require.Nil(t, err)
	require.NotNil(t, healthResponse)
	assert.Equal(t, expectedStatus, healthResponse.Status)
}

func TestMain_HealthCheckIsServing_Readiness(t *testing.T) {
	serviceToCheck := "readiness"

	if os.Getenv("EXECUTE") == "1" {
		main()
		return
	}

	envvar.SetAllNeceassaryEnvironemntVariables(t)
	runningApp := exec.Command(os.Args[0], "-test.run=TestMain_HealthCheckIsServing_Readiness")
	runningApp.Env = append(os.Environ(), "EXECUTE=1")
	blockingErr := runningApp.Start()
	require.Nil(t, blockingErr)
	defer runningApp.Process.Kill()

	expectedStatus := healthV1.HealthCheckResponse_SERVING
	healthResponse, err := waitForStatus(t, serviceToCheck, expectedStatus, 100*time.Millisecond, 500)

	require.Nil(t, err)
	require.NotNil(t, healthResponse)
	assert.Equal(t, expectedStatus, healthResponse.Status)
}

func waitForStatus(t *testing.T, serviceToCheck string, expectedStatus healthV1.HealthCheckResponse_ServingStatus, interval time.Duration, iterations int) (*healthV1.HealthCheckResponse, error) {
	conn, dialErr := grpc.Dial("localhost:3110", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.Nil(t, dialErr)
	defer conn.Close()

	client := healthV1.NewHealthClient(conn)
	count := 0
	var healthResponse *healthV1.HealthCheckResponse
	var err error
	for count < iterations {
		healthResponse, err = client.Check(context.Background(), &healthV1.HealthCheckRequest{Service: serviceToCheck})
		if healthResponse != nil && healthResponse.Status == expectedStatus {
			t.Logf("health check succeeded after %v iterations", count)
			break
		} else {
			t.Logf("health check failed after %v iterations", count)
		}
		time.Sleep(interval)
		count++
	}
	return healthResponse, err
}
