package appconfigagentv2

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const appConfigAgentImage = "public.ecr.aws/aws-appconfig/aws-appconfig-agent:2.x"

func setupAppConfigAgentTestcontainers(t *testing.T) string {
	t.Helper()
	dataDir, err := filepath.Abs(filepath.Join(".", "_testdata"))
	require.NoError(t, err)
	ctx := context.Background()
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
	require.NoError(t, err)

	testcontainers.CleanupContainer(t, ctr)
	baseURL, err := ctr.Endpoint(ctx, "http")
	require.NoError(t, err)
	return baseURL
}
