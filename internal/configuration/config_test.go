// +build integration

package configuration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testScenario struct {
	configFile         *os.File
	credsFile          *os.File
	defaultProfileFile *os.File
}

func (s *testScenario) teardown() {
	os.Remove(s.configFile.Name())
	os.Remove(s.credsFile.Name())
	os.Remove(s.defaultProfileFile.Name())
}

func setupTestScenario(t *testing.T) testScenario {
	configFile, err := ioutil.TempFile("", "config*.json")
	require.NoError(t, err)

	configJson := `
{
	"*": {
		"loglevel": "info",
		"plugindir": "/tmp",
		"prereleasefeatures": "NOT_ASKED",
		"sendusagedata": "NOT_ASKED"
	}
}
`
	_, err = configFile.Write([]byte(configJson))
	require.NoError(t, err)

	credsFile, err := ioutil.TempFile("", "credentials*.json")
	require.NoError(t, err)

	credsJson := `
{
	"default": {
		"apiKey": "testApiKey",
		"region": "US",
		"accountID": 12345,
		"licenseKey": "testLicenseKey"
	}
}
`
	_, err = credsFile.Write(([]byte(credsJson)))
	require.NoError(t, err)

	defaultProfileFile, err := ioutil.TempFile("", "default-profile*.json")
	require.NoError(t, err)

	defaultProfileJson := `"default"`
	_, err = defaultProfileFile.Write(([]byte(defaultProfileJson)))
	require.NoError(t, err)

	// package-level vars
	configFileName = filepath.Base(configFile.Name())
	credsFileName = filepath.Base(credsFile.Name())
	defaultProfileFileName = filepath.Base(defaultProfileFile.Name())
	configDir = filepath.Dir(configFile.Name())

	testScenario := testScenario{
		configFile:         configFile,
		credsFile:          credsFile,
		defaultProfileFile: defaultProfileFile,
	}

	return testScenario
}

func TestLoad(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()

	err := load()
	require.NoError(t, err)

	require.Equal(t, "info", GetConfigValue("logLevel"))
	require.Equal(t, "testApiKey", GetCredentialValue("apiKey"))
	require.Equal(t, "default", defaultProfileValue)
}

func TestSetCredentialValues(t *testing.T) {
	// Must be called first
	testScenario := setupTestScenario(t)
	defer testScenario.teardown()

	// Must load the config prior to tests
	err := load()
	require.NoError(t, err)

	err = SetAPIKey("default", "NRAK-abc123")
	require.NoError(t, err)
	require.Equal(t, "NRAK-abc123", GetCredentialValue(apiKeyKey))

	err = SetRegion("default", "US")
	require.NoError(t, err)
	require.Equal(t, "US", GetCredentialValue(regionKey))

	err = SetAccountID("default", "123456789")
	require.NoError(t, err)
	require.Equal(t, "123456789", GetCredentialValue(accountIDKey))
}

// Create config files if they don't already exist.
func TestCreate(t *testing.T) {
	require.True(t, true)
}