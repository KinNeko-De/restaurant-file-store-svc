//go:build unit

package server

import (
	"context"
	"os"
	"os/exec"
	"runtime"
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

func TestCreateFileRepository_PersistentVolume_ConfigMissing_Path(t *testing.T) {
	t.Setenv(StorageTypeEnv, "1")
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), PersistentVolumePathEnv)
}

func TestCreateFileRepository_PersistentVolume_ConfiguredPathDoesNotExists(t *testing.T) {
	t.Setenv(StorageTypeEnv, "1")
	t.Setenv(PersistentVolumePathEnv, "i-am-not-there-path")
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "for persistent volume was not found. Please check the configuration of the mounted volume")
}

func TestCreateFileRepository_PersistentVolume_ConfiguredPathIsNotAccessable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows does not support chmod on directories")
	}

	pathToNotAccessableDirectory := t.TempDir() + "/readonly"

	os.Mkdir(pathToNotAccessableDirectory, 0000)
	defer os.RemoveAll(pathToNotAccessableDirectory)

	t.Setenv(StorageTypeEnv, "1")
	t.Setenv(PersistentVolumePathEnv, pathToNotAccessableDirectory)
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "for persistent volume was not found. Please check the configuration of the mounted volume")
}

func TestCreateFileRepository_PersistentVolume_ConfiguredPathIsAccessable(t *testing.T) {
	pathToAccessableDirectory := t.TempDir() + "/accessable"
	os.Mkdir(pathToAccessableDirectory, 0777)

	t.Setenv(StorageTypeEnv, "1")
	t.Setenv(PersistentVolumePathEnv, pathToAccessableDirectory)
	_, err := createFileRepository(context.Background(), make(chan struct{}), make(chan struct{}))

	assert.Nil(t, err)
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
