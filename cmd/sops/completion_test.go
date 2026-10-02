package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Stubs for the zsh completion system and the sops binary, so the generated
// completion function can run outside an interactive shell. Like urfave/cli,
// the sops stub suggests flags whenever the argument before the
// --generate-bash-completion marker starts with "-", and subcommands otherwise.
const zshCompletionStubs = `
compdef() { :; }
_describe() { print -r -- "$1: ${(P)2}"; }
_files() { print -r -- files; }
sops() {
  if [[ "${@[-2]}" == -* ]]; then
    print -l -- --input-type --indent
  else
    print -l -- 'groups:Modify the groups on a SOPS file'
  fi
}
`

const zshCompletionRun = `
words=("$@")
CURRENT=$#
_cli_zsh_autocomplete
`

func TestZshCompletion(t *testing.T) {
	zsh, err := exec.LookPath("zsh")
	if err != nil {
		t.Skip("zsh is not installed")
	}

	tests := []struct {
		name  string
		words []string
		want  []string
	}{
		{
			name:  "file after flag",
			words: []string{"sops", "-d", "-i", "./dev"},
			want:  []string{"files"},
		},
		{
			name:  "empty word after flag",
			words: []string{"sops", "-d", "-i", ""},
			want:  []string{"files"},
		},
		{
			name:  "flag",
			words: []string{"sops", "-d", "-i", "--in"},
			want:  []string{"values: --input-type --indent"},
		},
		{
			name:  "subcommand or file",
			words: []string{"sops", ""},
			want:  []string{"values: groups:Modify the groups on a SOPS file", "files"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := zshCompletionStubs + GenZshCompletion("sops") + zshCompletionRun
			args := append([]string{"-f", "-c", script, "zsh"}, tt.words...)
			out, err := exec.Command(zsh, args...).CombinedOutput()
			require.NoError(t, err, string(out))
			assert.Equal(t, tt.want, strings.Split(strings.TrimSpace(string(out)), "\n"))
		})
	}
}
