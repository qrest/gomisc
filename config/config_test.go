package config

import (
	"flag"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestReadAndWriteConfig(t *testing.T) {
	require.Error(t, WriteConfig("", nil))

	type TestConfig struct {
		Host     string `yaml:"host"`
		Port     uint   `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}

	// create temp file and close it immediately
	file, err := os.CreateTemp("", "go_test_logfile")
	require.NoError(t, err)
	fName := file.Name()
	require.NoError(t, file.Close())

	var emptyConfig TestConfig
	require.NoError(t, WriteConfig(fName, emptyConfig))
	var otherConfig TestConfig
	require.NoError(t, ReadConfig(fName, &otherConfig))
	require.Equal(t, emptyConfig, otherConfig)

	notEmptyConfig := TestConfig{
		Host:     "Host1",
		Port:     123456,
		User:     "some_user",
		Password: "some_password",
	}
	require.NoError(t, WriteConfig(fName, notEmptyConfig))
	var otherConfig2 TestConfig
	require.NoError(t, ReadConfig(fName, &otherConfig2))
	require.Equal(t, notEmptyConfig, otherConfig2)

	require.Panics(t, func() { _ = ReadConfig(fName, nil) })
	require.Error(t, ReadConfig("file_which_does_not_exist", nil))
	// passed argument with no yaml config
	require.Error(t, ReadConfig(fName, true))
}

func TestSetConfigFlags(t *testing.T) {
	require.Panics(t, func() {
		SetConfigFlags("asdf", nil, nil)
	})

	require.NotPanics(t, func() {
		var createFile bool
		var filePath string
		SetConfigFlags("asdf", &filePath, &createFile)
		flag.Parse()
	})
}

func TestGetLogfile(t *testing.T) {
	// this should not work
	logfile, err := GetLogfile("")
	require.Error(t, err)
	require.Nil(t, logfile)

	// this should work

	// create temp file and close it immediately
	file, err := os.CreateTemp("", "go_test_logfile")
	require.NoError(t, err)
	fName := file.Name()
	require.NoError(t, file.Close())

	logfile, err = GetLogfile(fName)
	require.NoError(t, err)
	require.NotNil(t, logfile)
}
