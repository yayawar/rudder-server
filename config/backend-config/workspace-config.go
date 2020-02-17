package backendconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/rudderlabs/rudder-server/utils/logger"
)

type WorkspaceConfig struct {
}

func (workspaceConfig *WorkspaceConfig) SetUp() {
}

func (workspaceConfig *WorkspaceConfig) GetWorkspaceIDForWriteKey(writeKey string) string {
	return ""
}

//GetBackendConfig returns sources from the workspace
func (workspaceConfig *WorkspaceConfig) GetBackendConfig() (SourcesT, bool) {
	if configFromFile {
		return workspaceConfig.getBackendConfigFromFile()
	}
	return workspaceConfig.getBackendConfigFromAPI()
}

func (workspaceConfig *WorkspaceConfig) getBackendConfigFromAPI() (SourcesT, bool) {
	url := fmt.Sprintf("%s/workspaceConfig", configBackendURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error("Errored when sending request to the server", err)
		return SourcesT{}, false
	}

	req.SetBasicAuth(configBackendToken, "")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Errored when sending request to the server", err)
		return SourcesT{}, false
	}

	var respBody []byte
	if resp != nil && resp.Body != nil {
		respBody, _ = ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
	}

	var sourcesJSON SourcesT
	err = json.Unmarshal(respBody, &sourcesJSON)
	if err != nil {
		logger.Error("Errored while parsing request", err, string(respBody), resp.StatusCode)
		return SourcesT{}, false
	}
	return sourcesJSON, true
}

func (workspaceConfig *WorkspaceConfig) getBackendConfigFromFile() (SourcesT, bool) {
	logger.Info("Reading workspace config from JSON file")
	data, err := ioutil.ReadFile(configJSONPath)
	if err != nil {
		logger.Error("Unable to read backend config from file.")
		return SourcesT{}, false
	}
	var configJSON SourcesT
	error := json.Unmarshal(data, &configJSON)
	if error != nil {
		logger.Error("Unable to parse backend config from file.")
		return SourcesT{}, false
	}
	return configJSON, true
}
