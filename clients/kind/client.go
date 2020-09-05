package kind

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Client interface {
	WaitForReadiness(ctx context.Context) (err error)
	PrepareKubeConfig(ctx context.Context) (err error)
}

// NewClient returns a new kind.Client
func NewClient(ctx context.Context, kindHost string) (Client, error) {
	if kindHost == "" {
		return nil, fmt.Errorf("kindHost is empty, this is now allowed")
	}

	return &client{
		kindHost: kindHost,
	}, nil
}

type client struct {
	kindHost string
}

func (c *client) WaitForReadiness(ctx context.Context) (err error) {
	log.Info().Msg("Waiting for kind host to be ready...")
	httpClient := &http.Client{
		Timeout: time.Second * 1,
	}

	for true {
		_, err := httpClient.Get(fmt.Sprintf("http://%v:10080/kubernetes-ready", c.kindHost))
		if err == nil {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	return nil
}

func (c *client) PrepareKubeConfig(ctx context.Context) (err error) {

	log.Info().Msg("Preparing kind host for using Helm...")
	httpClient := &http.Client{
		Timeout: time.Second * 1,
	}
	response, err := httpClient.Get(fmt.Sprintf("http://%v:10080/config", c.kindHost))
	if err != nil {
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	kubeConfig := strings.ReplaceAll(string(body), "localhost", c.kindHost)

	usr, _ := user.Current()
	homeDir := usr.HomeDir
	err = ioutil.WriteFile(filepath.Join(homeDir, ".kube/config"), []byte(kubeConfig), 0644)
	if err != nil {
		return
	}

	return nil
}
