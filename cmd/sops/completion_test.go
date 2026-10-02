//go:build !windows

package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Stubs for the shell completion systems and the sops binary, so the generated
// completion functions can run outside an interactive shell. Like urfave/cli,
// the sops stubs suggest flags whenever the argument before the
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

const bashCompletionStubs = `
complete() { :; }
sops() {
  local args=("$@")
  if (( $# > 1 )) && [[ "${args[$# - 2]}" == -* ]]; then
    printf '%s\n' --input-type --indent
  else
    printf '%s\n' groups
  fi
}
`

const bashCompletionRun = `
COMP_WORDS=("$@")
COMP_CWORD=$(( $# - 1 ))
_cli_bash_autocomplete
echo "${COMPREPLY[*]}"
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

// An empty COMPREPLY makes bash fall back to file name completion, because the
// script registers the function with "complete -o default".
func TestBashCompletion(t *testing.T) {
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash is not installed")
	}

	tests := []struct {
		name  string
		words []string
		want  string
	}{
		{
			name:  "file after flag",
			words: []string{"sops", "-d", "-i", "./dev"},
			want:  "",
		},
		{
			name:  "empty word after flag",
			words: []string{"sops", "-d", "-i", ""},
			want:  "",
		},
		{
			name:  "flag",
			words: []string{"sops", "-d", "-i", "--in"},
			want:  "--input-type --indent",
		},
		{
			name:  "subcommand",
			words: []string{"sops", "gr"},
			want:  "groups",
		},
		{
			name:  "file",
			words: []string{"sops", "./dev"},
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := bashCompletionStubs + GenBashCompletion("sops") + bashCompletionRun
			args := append([]string{"--noprofile", "--norc", "-c", script, "bash"}, tt.words...)
			out, err := exec.Command(bash, args...).CombinedOutput()
			require.NoError(t, err, string(out))
			assert.Equal(t, tt.want, strings.TrimSpace(string(out)))
		})
	}
}
