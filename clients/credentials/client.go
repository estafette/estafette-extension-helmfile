package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

type Client interface {
	Init(ctx context.Context) (err error)
}

func NewClient(ctx context.Context, infraCredentialsJSON, serviceAccountKeyfilePath string) (Client, error) {
	if infraCredentialsJSON == "" {
		return nil, fmt.Errorf("infraCredentialsJSON, this is now allowed")
	}
	if serviceAccountKeyfilePath == "" {
		return nil, fmt.Errorf("serviceAccountKeyfilePath, this is now allowed")
	}

	return &client{
		infraCredentialsJSON:      infraCredentialsJSON,
		serviceAccountKeyfilePath: serviceAccountKeyfilePath,
	}, nil
}

type client struct {
	infraCredentialsJSON      string
	serviceAccountKeyfilePath string
}

func (c *client) Init(ctx context.Context) (err error) {

	log.Info().Msg("Initializing credentials...")

	// read injected gke-update credentials
	serviceAccountKeyfile, err := c.getServiceAccountKeyfile(ctx, c.infraCredentialsJSON)
	if err != nil {
		return err
	}

	err = c.storeServiceAccountKeyfile(ctx, serviceAccountKeyfile, c.serviceAccountKeyfilePath)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) getServiceAccountKeyfile(ctx context.Context, infraCredentialsJSON string) (string, error) {
	log.Debug().Msg("Unmarshalling injected gcp-infra credentials...")

	var gcpInfraCredentials []GCPInfraCredentials
	if err := json.Unmarshal([]byte(infraCredentialsJSON), &gcpInfraCredentials); err != nil {
		return "", err
	}
	if len(gcpInfraCredentials) == 0 {
		return "", fmt.Errorf("No gcp-infra credentials injected")
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
