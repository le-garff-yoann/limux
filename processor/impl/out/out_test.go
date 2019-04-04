package out

import (
	"limux/processor"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutConfigure(t *testing.T) {
	t.Parallel()

	src, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	dst, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	defer func() {
		os.RemoveAll(src)
		os.RemoveAll(dst)
	}()

	var (
		dur = time.Second * 1

		exec []string
	)

	p := &Out{
		Interval:            &dur,
		Src:                 src,
		Dst:                 dst,
		ArchiveBasename:     "foo",
		ArchiveInnerDirname: "in",
		Exec:                &exec,
	}

	require.NoError(t, p.Configure())

	// FIXME: this is quick and dirty.
	p.Src = "./hopeitdoesnotexists"
	require.NoError(t, p.Configure())

	// FIXME: this is quick and dirty.
	p.Dst = "./hopeitdoesnotexists"
	require.Error(t, p.Configure())

	p.Dst = dst
	p.ArchiveBasename = ""
	require.Error(t, p.Configure())

	p.ArchiveBasename = "foo_{{ .Now.UnixNano }"
	require.Error(t, p.Configure())

	p.ArchiveBasename = "foo"
	exec = append(exec, "echo {{ .ArchiveFullPath }}")
	require.NoError(t, p.Configure())

	exec = append(exec, "echo {{ .ArchiveFullPath }")
	require.Error(t, p.Configure())
}

func TestOut(t *testing.T) {
	t.Parallel()

	src, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	dst, _ := ioutil.TempDir(os.TempDir(), "t_limux")
	defer func() {
		os.RemoveAll(src)
		os.RemoveAll(dst)
	}()

	var (
		dur = time.Second * 1

		exec []string
	)

	out := &Out{
		Interval:        &dur,
		Src:             src,
		Dst:             dst,
		ArchiveBasename: "x",
		Exec:            &exec,
	}

	events := make(chan processor.Event)

	out.Start(events)

	file, err := ioutil.TempFile(src, "t_limux")
	defer file.Close()
	require.NoError(t, err)
	file.Close()

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
