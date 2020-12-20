package config_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/fajardm/gobackend-server/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad_PanicError(t *testing.T) {
	res, err := config.Load("xxx")

	assert.Nil(t, res)
	assert.Error(t, err)
}

func TestLoad_ReadInConfigError(t *testing.T) {
	mockJSON := `{"app": {"name": "app", "version": "1.0",}`
	mockPath := "./config.test.json"

	err := ioutil.WriteFile(mockPath, []byte(mockJSON), 0644)
	assert.NoError(t, err)

	res, err := config.Load(mockPath)

	assert.Nil(t, res)
	assert.Error(t, err)

	err = os.Remove(mockPath)
	assert.NoError(t, err)
}

func TestLoad_UnmarshalError(t *testing.T) {
	mockJSON := `{"app": {"name": "app", "version": {"xxx":"xxx"}}}`
	mockPath := "./config.test.json"

	err := ioutil.WriteFile(mockPath, []byte(mockJSON), 0644)
	assert.NoError(t, err)

	res, err := config.Load(mockPath)

	assert.Nil(t, res)
	assert.Error(t, err)

	err = os.Remove(mockPath)
	assert.NoError(t, err)
}

func TestLoad_NoError(t *testing.T) {
	mockJSON := `{"app": {"name": "app", "version": "1.0"}}`
	mockPath := "./config.test.json"

	err := ioutil.WriteFile(mockPath, []byte(mockJSON), 0644)
	assert.NoError(t, err)

	res, err := config.Load(mockPath)

	assert.Equal(t, "app", res.App.Name)
	assert.Equal(t, "1.0", res.App.Version)
	assert.NoError(t, err)

	err = os.Remove(mockPath)
	assert.NoError(t, err)
}
