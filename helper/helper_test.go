package helper

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDirValid(t *testing.T) {
	t.Parallel()

	// FIXME: this is quick and dirty.
	isDir, err := IsDir(".")
	require.Nil(t, err)
	require.True(t, isDir)
}

func TestIsDirInvalid(t *testing.T) {
	t.Parallel()

	// FIXME: this is quick and dirty.
	isDir, err := IsDir("./helper.go")
	require.Nil(t, err)
	require.False(t, isDir)
}

func TestIsDirInexistant(t *testing.T) {
	t.Parallel()

	// FIXME: this is quick and dirty.
	_, err := IsDir("./hopeitdoesnotexists")
	require.Error(t, err)
}

func TestParseTemplateValid(t *testing.T) {
	t.Parallel()

	var tpl bytes.Buffer
	require.Nil(t, ParseTemplate(`test {{ .foo | replace "bar" "fux" }}`,
		&tpl, map[string]interface{}{
			"foo": "bar",
		},
	))

	require.Equal(t, "test fux", tpl.String())
}

func TestParseTemplateInvalid(t *testing.T) {
	t.Parallel()

	var tpl bytes.Buffer
	require.Error(t, ParseTemplate(`test {{ .foo }`, &tpl))
}
