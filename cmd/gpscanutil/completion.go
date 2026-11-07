package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// flagSpec describes an option for the generated completion scripts.
type flagSpec struct {
	name     string
	arg      string // value placeholder; "" for boolean options
	help     string
	dirArg   bool     // complete the value as a directory
	excludes []string // options that can't be combined with this one
}

var flagSpecs = []flagSpec{
	{name: "top", arg: "N", help: "number of entries (per folder with --tree) to show"},
	{name: "depth", arg: "N", help: "only descend N levels below the root"},
	{name: "min-size", arg: "SIZE", help: "hide smaller entries/changes, e.g. 500M or 1.5G"},
	{name: "root", arg: "PATH", help: "start at PATH instead of the volume root", dirArg: true},
	{name: "files", help: "only show files", excludes: []string{"dirs", "tree"}},
	{name: "dirs", help: "only show folders", excludes: []string{"files"}},
	{name: "tree", help: "print an indented tree instead of a flat list", excludes: []string{"files"}},
	{name: "bytes", help: "print sizes as exact byte counts"},
	{name: "q", help: "no progress messages on stderr"},
}

var commands = []struct{ name, help string }{
	{"introspect", "list the biggest files/folders of a scan"},
	{"diff", "list the biggest changes between two scans"},
	{"completion", "print a shell completion script"},
	{"version", "print the version"},
	{"help", "show usage"},
}

// dashed returns the option as typed on the command line.
func dashed(name string) string {
	if len(name) == 1 {
		return "-" + name
	}
	return "--" + name
}

func runCompletion(args []string) error {
	if len(args) != 1 {
		return usageError("completion expects one argument: zsh or bash")
	}
	var script string
	switch args[0] {
	case "zsh":
		script = zshCompletion()
	case "bash":
		script = bashCompletion()
	default:
		return usageError(fmt.Sprintf("unsupported shell %q (supported: zsh, bash)", args[0]))
	}
	_, err := os.Stdout.WriteString(script)
	return err
}

// zshQuote escapes s for use inside a single-quoted _arguments/_describe
// spec, where ':', '[' and ']' are special too.
func zshQuote(s string) string {
	return strings.NewReplacer(`'`, `'\''`, `:`, `\:`, `[`, `\[`, `]`, `\]`).Replace(s)
}

func zshCompletion() string {
	var opts, cmds []string
	for _, f := range flagSpecs {
		spec := ""
		if len(f.excludes) > 0 {
			ex := make([]string, len(f.excludes))
			for i, e := range f.excludes {
				ex[i] = dashed(e)
			}
			spec = "(" + strings.Join(ex, " ") + ")"
		}
		spec += dashed(f.name)
		if f.arg != "" {
			spec += "="
		}
		spec += "[" + zshQuote(f.help) + "]"
		if f.arg != "" {
			action := ""
			if f.dirArg {
				action = "_directories"
			}
			spec += ":" + zshQuote(f.arg) + ":" + action
		}
		opts = append(opts, "'"+spec+"'")
	}
	for _, c := range commands {
		cmds = append(cmds, "'"+zshQuote(c.name)+":"+zshQuote(c.help)+"'")
	}
	return render(zshTemplate, map[string]string{
		"Opts":     strings.Join(opts, "\n    "),
		"Commands": strings.Join(cmds, "\n        "),
	})
}

func bashCompletion() string {
	var all, valued, dirValued []string
	for _, f := range flagSpecs {
		all = append(all, dashed(f.name))
		switch {
		case f.dirArg:
			dirValued = append(dirValued, dashed(f.name))
		case f.arg != "":
			valued = append(valued, dashed(f.name))
		}
	}
	var cmds []string
	for _, c := range commands {
		cmds = append(cmds, c.name)
	}
	return render(bashTemplate, map[string]string{
		"Commands":  strings.Join(cmds, " "),
		"Flags":     strings.Join(all, " "),
		"Valued":    strings.Join(valued, "|"),
		"DirValued": strings.Join(dirValued, "|"),
	})
}

func render(tmpl string, data map[string]string) string {
	var b strings.Builder
	template.Must(template.New("").Parse(tmpl)).Execute(&b, data)
	return b.String()
}

// The zsh script works both when eval'ed/sourced and when installed as
// _gpscanutil in a directory on $fpath.
const zshTemplate = `#compdef gpscanutil
# zsh completion for gpscanutil. Load with:
#   eval "$(gpscanutil completion zsh)"     (in ~/.zshrc, after compinit)
# or install it on your $fpath:
#   gpscanutil completion zsh > "${fpath[1]}/_gpscanutil"

_gpscanutil() {
  local curcontext="$curcontext" state line
  typeset -A opt_args
  local -a opts=(
    {{.Opts}}
  )
  local scan='_files -g "*.gpscan(-.)"'

  _arguments -C '1:command:->command' '*::arg:->args' && return
  case $state in
    command)
      local -a cmds=(
        {{.Commands}}
      )
      _describe -t commands 'gpscanutil command' cmds
      ;;
    args)
      case $line[1] in
        introspect)
          _arguments -S $opts "1:scan file:$scan"
          ;;
        diff)
          _arguments -S $opts "1:old scan file:$scan" "2:new scan file:$scan"
          ;;
        completion)
          _arguments '1:shell:(zsh bash)'
          ;;
      esac
      ;;
  esac
}

if [[ "${funcstack[1]}" == _gpscanutil ]]; then
  _gpscanutil "$@"
else
  compdef _gpscanutil gpscanutil
fi
`

// The bash script avoids bash-completion helpers and bash 4 features, so it
// also works with macOS's bash 3.2.
const bashTemplate = `# bash completion for gpscanutil. Load with:
#   eval "$(gpscanutil completion bash)"     (in ~/.bashrc)

_gpscanutil() {
  local cur=${COMP_WORDS[COMP_CWORD]} prev=${COMP_WORDS[COMP_CWORD-1]}
  # bash splits "--opt=value" into "--opt" "=" "value".
  if [[ $cur == "=" ]]; then
    cur=""
  elif [[ $prev == "=" && $COMP_CWORD -ge 2 ]]; then
    prev=${COMP_WORDS[COMP_CWORD-2]}
  fi
  local IFS=$'\n'
  COMPREPLY=()

  if [[ $COMP_CWORD -eq 1 ]]; then
    COMPREPLY=($(IFS=' '; compgen -W "{{.Commands}}" -- "$cur"))
    return
  fi

  case ${COMP_WORDS[1]} in
    completion)
      [[ $COMP_CWORD -eq 2 ]] && COMPREPLY=($(IFS=' '; compgen -W "zsh bash" -- "$cur"))
      ;;
    introspect|diff)
      case $prev in
        {{.Valued}}) return ;;
        {{.DirValued}}) COMPREPLY=($(compgen -d -- "$cur")); return ;;
      esac
      if [[ $cur == -* ]]; then
        COMPREPLY=($(IFS=' '; compgen -W "{{.Flags}}" -- "$cur"))
      else
        COMPREPLY=($(compgen -d -- "$cur") $(compgen -f -X '!*.gpscan' -- "$cur"))
      fi
      ;;
  esac
}

complete -o filenames -F _gpscanutil gpscanutil
`
