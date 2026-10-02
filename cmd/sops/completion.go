package main

import "fmt"

// Based on https://github.com/urfave/cli/blob/v1-maint/autocomplete/zsh_autocomplete,
// which only completes files when sops prints no suggestions at all. sops
// almost always prints flags or subcommands, so files are offered alongside
// subcommands instead.
var Zshcompletion = `#compdef %s

_cli_zsh_autocomplete() {

  local -a opts
  local cur prev ret=1
  cur=${words[CURRENT]}
  prev=${words[CURRENT-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[1,CURRENT-1]} ${cur} --generate-bash-completion)}")
    [[ "${opts[1]}" != "" ]] && _describe 'values' opts && ret=0
    return $ret
  fi

  # After a flag, sops suggests more flags rather than flag values or files,
  # so only ask it for subcommands when the previous word is not a flag.
  if [[ "$prev" != "-"* ]]; then
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[1,CURRENT-1]} --generate-bash-completion)}")
    [[ "${opts[1]}" != "" ]] && _describe 'values' opts && ret=0
  fi

  _files && ret=0
  return $ret
}

compdef _cli_zsh_autocomplete %s
`

// https://github.com/urfave/cli/blob/v1-maint/autocomplete/bash_autocomplete
var Bashcompletion = `#! /bin/bash

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}
complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete %s
`

func GenBashCompletion(name string) string {
	return fmt.Sprintf(Bashcompletion, name)
}

func GenZshCompletion(name string) string {
	return fmt.Sprintf(Zshcompletion, name, name)
}
