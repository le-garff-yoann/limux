package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewValidConfig(t *testing.T) {
	t.Parallel()

	dir, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	defer os.RemoveAll(dir)

	conf, err := New([]byte(fmt.Sprintf(`out:
  - interval: 10
    src: %s
    dst: %s
    archive_basename: foo_{{ .Now.UnixNano }}
    archive_inner_dirname: in
    exec: ["cmd.exe", "/c", "echo 1"]

in:
  - src: %s
    dst: %s`, dir, dir, dir, dir)))
	require.NoError(t, err)
	require.Len(t, conf.Processors(), 2)

	require.Equal(t, time.Duration(10), *conf.Out[0].Interval)
	require.Len(t, *conf.Out[0].Exec, 3)
}

func TestNewInvalidConfig(t *testing.T) {
	t.Parallel()

	_, err := New([]byte("foo"))
	require.Error(t, err)
}
