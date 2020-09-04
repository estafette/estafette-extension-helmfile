package main

import (
	"context"
	"runtime"

	"github.com/alecthomas/kingpin"
	"github.com/estafette/estafette-extension-helmfile/clients/credentials"
	"github.com/estafette/estafette-extension-helmfile/services/helmfile"
	"github.com/estafette/estafette-extension-helmfile/services/kind"
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
	action                    = kingpin.Flag("action", ".").Envar("ESTAFETTE_EXTENSION_ACTION").Enum(string(ActionLint), string(ActionApply), string(ActionDiff))
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

	// extract service account keyfile from injected credentials
	credentialsClient, err := credentials.NewClient(ctx, *infraCredentialsJSON, *serviceAccountKeyfilePath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating credentials.Client")
	}

	kindService, err := kind.NewService(ctx, *kindHost)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating kind.Service")
	}

	helmfileService, err := helmfile.NewService(ctx, credentialsClient, kindService, *file)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed creating helmfile.Service")
	}

	switch *action {
	case string(ActionLint):
		err = helmfileService.Lint(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed linting helmfile")
		}

	case string(ActionDiff):
		err = helmfileService.Diff(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed diffing helmfile")
		}

	case string(ActionApply):
		err = helmfileService.Apply(ctx)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed applying helmfile")
		}

	default:
		log.Fatal().Msgf("action %v is not supported", *action)
	}
}
