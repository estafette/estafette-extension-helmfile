package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

type Client interface {
	Init(ctx context.Context) (err error)
}

func NewClient(ctx context.Context, credentialsPath, serviceAccountKeyfilePath string) (Client, error) {
	if credentialsPath == "" {
		return nil, fmt.Errorf("credentialsPath, this is now allowed")
	}
	if serviceAccountKeyfilePath == "" {
		return nil, fmt.Errorf("serviceAccountKeyfilePath, this is now allowed")
	}

	return &client{
		credentialsPath:           credentialsPath,
		serviceAccountKeyfilePath: serviceAccountKeyfilePath,
	}, nil
}

type client struct {
	credentialsPath           string
	serviceAccountKeyfilePath string
}

func (c *client) Init(ctx context.Context) (err error) {

	log.Info().Msg("Initializing credentials...")

	// read injected gke-update credentials
	serviceAccountKeyfile, err := c.getServiceAccountKeyfile(ctx, c.credentialsPath)
	if err != nil {
		return err
	}

	err = c.storeServiceAccountKeyfile(ctx, serviceAccountKeyfile, c.serviceAccountKeyfilePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getServiceAccountKeyfile(ctx context.Context, credentialsPath string) (string, error) {
	log.Debug().Msg("Unmarshalling injected gcp-infra credentials...")

	// use mounted credential file if present instead of relying on an envvar
	if runtime.GOOS == "windows" {
		credentialsPath = "C:" + credentialsPath
	}
	var gcpInfraCredentials []GCPInfraCredentials
	if foundation.FileExists(credentialsPath) {
		log.Info().Msgf("Reading credentials from file at path %v...", credentialsPath)
		credentialsFileContent, err := ioutil.ReadFile(credentialsPath)
		if err != nil {
			return "", err
		}
		var gcpInfraCredentials []GCPInfraCredentials
		if err := json.Unmarshal([]byte(credentialsFileContent), &gcpInfraCredentials); err != nil {
			return "", err
		}
		if len(gcpInfraCredentials) == 0 {
			return "", fmt.Errorf("No gcp-infra credentials injected")
		}
	} else {
		return "", fmt.Errorf("Credentials of type gcp-infra are not injected; configure this extension as trusted and inject credentials of type gcp-infra")
	}

	return gcpInfraCredentials[0].AdditionalProperties.ServiceAccountKeyfile, nil
}

func (c *client) storeServiceAccountKeyfile(ctx context.Context, serviceAccountKeyfile, path string) error {
	log.Debug().Msg("Storing gcp-infra credential service account keyfile on disk...")

	pathDir := filepath.Dir(path)
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		err = os.MkdirAll(pathDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(path, []byte(serviceAccountKeyfile), 0600); err != nil {
		return err
	}

	var keyFileMap map[string]interface{}
	err := json.Unmarshal([]byte(serviceAccountKeyfile), &keyFileMap)
	if err == nil {
		if clientEmail, ok := keyFileMap["client_email"]; ok {
			if clientEmailString, castOK := clientEmail.(string); !castOK {
				log.Debug().Msgf("Stored keyfile for service account %v on disk", clientEmailString)
				return nil
			}
		}
	}

	return nil
}
