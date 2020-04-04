package docker

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

const image = "nginx:1.17-alpine"

func TestInspect(t *testing.T) {
	t.Parallel()

	// append timestamp to container name to allow running tests in parallel
	name := "inspect-test-" + strconv.FormatInt(time.Now().UnixNano(), 10)

	// running the container detached to allow inspection while it is running
	options := &RunOptions{
		Detach: true,
		Name:   name,
	}

	id := Run(t, image, options)
	defer removeContainer(t, id)

	c := Inspect(t, id)

	require.Equal(t, id, c.ID)
	require.Equal(t, name, c.Name)
	require.IsType(t, time.Time{}, c.Created)
	require.Equal(t, true, c.Running)
}

func TestInspectWithExposedPort(t *testing.T) {
	t.Parallel()

	// choosing an unique high port to avoid conflict on test machines
	port := 13031

	options := &RunOptions{
		Detach: true,
		OtherOptions:         []string{fmt.Sprintf("-p=%d:80", port)},
	}

	id := Run(t, image, options)
	defer removeContainer(t, id)

	c := Inspect(t, id)

	require.NotEmptyf(t, c.Ports, "Container's exposed ports should not be empty")
	require.EqualValues(t, 80, c.Ports[0].ContainerPort)
	require.EqualValues(t, port, c.Ports[0].HostPort)
}

func TestInspectWithMappedVolumes(t *testing.T) {
	t.Parallel()

	options := &RunOptions{
		Detach: true,
		Volumes: []string{"/tmp:/foo/bar"},
	}

	id := Run(t, image, options)
	defer removeContainer(t, id)

	c := Inspect(t, id)

	require.NotEmptyf(t, c.Binds, "Container's mapped volumes should not be empty")
	require.Equal(t, "/tmp", c.Binds[0].Source)
	require.Equal(t, "/foo/bar", c.Binds[0].Destination)
}

func removeContainer(t *testing.T, id string) {
	cmd := shell.Command{
		Command: "docker",
		Args:    []string{"container", "rm", "-f", id},
	}

	shell.RunCommand(t, cmd)
}