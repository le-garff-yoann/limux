package cmd

import (
	"filemux/config"
	"filemux/processor/broadcaster"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/posener/wstest"
	"github.com/stretchr/testify/require"
)

func TestSubcommandRunInitConfigExistingFile(t *testing.T) {
	t.Parallel()

	tfile, _ := ioutil.TempFile(os.TempDir(), "t_filemux")
	defer os.RemoveAll(tfile.Name())

	conf, err := initConfig(tfile.Name())
	require.NotNil(t, conf)
	require.NoError(t, err)
}

func TestSubcommandRunInitConfigNotExistingFile(t *testing.T) {
	t.Parallel()

	// FIXME: this is quick and dirty.
	conf, err := initConfig("./hopeitdoesnotexists")
	require.Nil(t, conf)
	require.Error(t, err)
}

func TestSubcommandRunRouterEvents(t *testing.T) {
	t.Parallel()

	dir, _ := ioutil.TempDir(os.TempDir(), "t_filemux")
	defer os.RemoveAll(dir)

	conf, _ := config.New([]byte(fmt.Sprintf(`out:
  - interval: 10
    src: %s
    dst: %s
    archive_basename: foo_{{ .Now.UnixNano }}
    archive_inner_dirname: in
    exec: []

in:
  - src: %s
    dst: %s`, dir, dir, dir, dir)))

	var (
		br = broadcaster.New()

		r = router(conf, br)
	)

	h := http.Handler(r.Get("events").GetHandler())

	d := wstest.NewDialer(h)
	c, resp, err := d.Dial("ws://x/x", nil)
	require.NoError(t, err)

	defer c.Close()

	require.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	var msg interface{}

	require.NoError(t, c.ReadJSON(&msg))
	require.NoError(t, c.ReadJSON(&msg))
}
