package pflag

import (
	"testing"

	logger "github.com/apsdehal/go-logger"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

func TestLogLevel(t *testing.T) {
	t.Parallel()

	require.Implements(t, (*pflag.Value)(nil), new(LogLevel))

	var lvl LogLevel

	for _, lvll := range LogLevelLitterals {
		err := lvl.Set(lvll)
		require.NoError(t, err)

		require.Equal(t, lvll, lvl.String())
	}

	err := lvl.Set("hopeitdoesnotexists")
	require.Error(t, err)
}

func TestLogLevelType(t *testing.T) {
	t.Parallel()

	lvl := LogLevel(logger.InfoLevel)
	require.Equal(t, "loglevel", lvl.Type())
}
