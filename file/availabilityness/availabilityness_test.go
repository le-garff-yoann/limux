package availabilityness

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAvailabilityness(t *testing.T) {
	tfile, _ := ioutil.TempFile(os.TempDir(), "t_filemux")
	defer os.RemoveAll(tfile.Name())

	fanessFiles, fanessLogs := Watch(tfile.Name())
	require.NotNil(t, fanessFiles)
	require.NotNil(t, fanessLogs)

	// FIXME: this is quick and dirty.
	timeout := time.NewTimer(CheckInterval * CheckAttempts * 2)

	var logCnt int

	for {
		select {
		case f, _ := <-fanessFiles:
			require.Equal(t, tfile.Name(), f.Fname)
			require.NotNil(t, f.File)
			require.NotZero(t, logCnt)

			return
		case e := <-fanessLogs:
			logCnt++

			t.Log(e)
		case <-timeout.C:
			t.Error("Abnormally blocked")

			return
		}
	}
}

func TestAvailabilitynessNonExistingFile(t *testing.T) {
	// FIXME: this is a quick and dirty.
	nonExistingFile := "./hopeitdoesnotexists"

	fanessFiles, fanessLogs := Watch(nonExistingFile)
	require.NotNil(t, fanessFiles)
	require.NotNil(t, fanessLogs)

	// FIXME: this is quick and dirty.
	timeout := time.NewTimer(CheckInterval * CheckAttempts * 2)

	var logCnt int

	for {
		select {
		case f := <-fanessFiles:
			require.Equal(t, nonExistingFile, f.Fname)
			require.Nil(t, f.File)

			return
		case e := <-fanessLogs:
			logCnt++

			t.Log(e)
		case <-timeout.C:
			t.Error("Abnormally blocked")

			return
		}
	}
}
