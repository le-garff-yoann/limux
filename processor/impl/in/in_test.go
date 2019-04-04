package in

import (
	"limux/processor"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mholt/archiver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInConfigure(t *testing.T) {
	t.Parallel()

	src, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	dst, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	defer func() {
		os.RemoveAll(src)
		os.RemoveAll(dst)
	}()

	p := &In{
		Src: src,
		Dst: dst,
	}

	require.NoError(t, p.Configure())

	// FIXME: this is quick and dirty.
	p.Src = "./hopeitdoesnotexists"
	require.Error(t, p.Configure())

	// FIXME: this is quick and dirty.
	p.Dst = "./hopeitdoesnotexists"
	require.Error(t, p.Configure())

	p.Src = src
	require.Error(t, p.Configure())
}

func TestIn(t *testing.T) {
	t.Parallel()

	src, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	dst, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	defer func() {
		os.RemoveAll(src)
		os.RemoveAll(dst)
	}()

	p := &In{
		Src: src,
		Dst: dst,
	}

	events := make(chan processor.Event)

	require.NoError(t, p.Start(events))

	// FIXME: this is quick and dirty.
	require.NoError(t, archiver.Archive([]string{filepath.Join(".", "in.go")}, filepath.Join(src, "x.tar")))

	for {
		select {
		case e := <-events:
			t.Log(e)

			if e.Type == processor.Fin {
				files, err := ioutil.ReadDir(src)
				require.NoError(t, err)
				require.Empty(t, files)

				files, err = ioutil.ReadDir(dst)
				require.NoError(t, err)
				assert.NotEmpty(t, files)

				return
			}
		}
	}
}
