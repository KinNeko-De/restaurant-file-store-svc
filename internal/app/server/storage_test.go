//go:build unit

package server

import (
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFileRepository_StorageType_InvalidValue(t *testing.T) {
	t.Setenv(StorageTypeEnv, "IMustBeAnInteger")
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "IMustBeAnInteger")
}

func TestCreateFileRepository_ConfigMissing_StorageTypeEnv(t *testing.T) {
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), StorageTypeEnv)
}

func TestCreateFileRepository_PersistenceVolume_ConfigMissing_Path(t *testing.T) {
	t.Setenv(StorageTypeEnv, "1")
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), PersistentVolumePathEnv)
}

func TestCreateFileRepository_GoogleStorage_ConfigMissing(t *testing.T) {
	t.Setenv(StorageTypeEnv, "2")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "i-am-not-there-credentials.json")
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	assert.NotNil(t, err)
}

func TestCreateFileRepository_GoogleStorage_DummyConfig(t *testing.T) {
	t.Setenv(StorageTypeEnv, "2")
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../../test/testing/googlecloud/dummycredentials.json")
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	assert.Nil(t, err)
}

func TestInitializeStorage_AnyError_AppCrash(t *testing.T) {
	if os.Getenv("EXECUTE") == "1" {
		// will crash because of missing environment variables
		InitializeStorage(context.Background(), make(chan struct{}), make(chan struct{}))
		return
	}

	runningApp := exec.Command(os.Args[0], "-test.run=TestInitializeStorage_AnyError_AppCrash")
	runningApp.Env = append(os.Environ(), "EXECUTE=1")
	err := runningApp.Run()
	require.NotNil(t, err)
	exitCode := err.(*exec.ExitError).ExitCode()
	assert.Equal(t, 55, exitCode)
}
