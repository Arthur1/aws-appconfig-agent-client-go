package appconfigagentv2

import (
	"context"
	"path/filepath"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const appConfigAgentImage = "public.ecr.aws/aws-appconfig/aws-appconfig-agent:2.x"

func setupAppConfigAgentTestcontainers(ctx context.Context) (testcontainers.Container, error) {
	dataDir, err := filepath.Abs(filepath.Join(".", "_testdata"))
	if err != nil {
		return nil, err
	}
	req := testcontainers.ContainerRequest{
		Image: appConfigAgentImage,
		Files: []testcontainers.ContainerFile{{
			HostFilePath:      dataDir,
			ContainerFilePath: "/",
			FileMode:          0o700,
		}},
		Env: map[string]string{
			"LOCAL_DEVELOPMENT_DIRECTORY": "/_testdata/",
		},
		ExposedPorts: []string{"2772/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("2772/tcp"),
			wait.ForLog("serving on localhost:2772"),
		),
	}
	ctr, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	return ctr, err
}
