package main

import (
	"context"
	"runtime"

	"github.com/alecthomas/kingpin"
	"github.com/estafette/estafette-extension-helmfile/clients/credentials"
	"github.com/estafette/estafette-extension-helmfile/clients/helmfile"
	"github.com/estafette/estafette-extension-helmfile/clients/kind"
	"github.com/estafette/estafette-extension-helmfile/services/extension"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

var (
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	action                    = kingpin.Flag("action", ".").Envar("ESTAFETTE_EXTENSION_ACTION").Enum(extension.AllowedActions()...)
	infraCredentialsJSON      = kingpin.Flag("gcp-infra-credentials", "GCP infra credentials configured at service level, passed in to this trusted extension.").Envar("ESTAFETTE_CREDENTIALS_GCP_INFRA").Required().String()
	serviceAccountKeyfilePath = kingpin.Flag("service-account-keyfile-path", "Path to store the service account keyfile.").Envar("GOOGLE_APPLICATION_CREDENTIALS").Required().String()
	kindHost                  = kingpin.Flag("kind-host", "Hostname of kind container.").Default("kubernetes").OverrideDefaultFromEnvar("ESTAFETTE_EXTENSION_KIND_HOST").String()
	file                      = kingpin.Flag("file", "Yaml file to be used by helmfile.").Default("helmfile.yaml").OverrideDefaultFromEnvar("ESTAFETTE_EXTENSION_FILE").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(foundation.NewApplicationInfo(appgroup, app, version, branch, revision, buildDate))

	// create context to cancel commands on sigterm
	ctx := foundation.InitCancellationContext(context.Background())

	credentialsClient, err := credentials.NewClient(ctx, *infraCredentialsJSON, *serviceAccountKeyfilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating credentials.Client")
	}

	kindClient, err := kind.NewClient(ctx, *kindHost)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating kind.Client")
	}

	helmfileClient, err := helmfile.NewClient(ctx, *file)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating helmfile.Client")
	}

	extensionService, err := extension.NewService(ctx, credentialsClient, kindClient, helmfileClient)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating extension.Service")
	}

	// do the actual work
	err = extensionService.ExecuteAction(ctx, extension.Action(*action))
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed executing action %v", *action)
	}
}
