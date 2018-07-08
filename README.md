
# Kubectx

Command line utility to manage the current environment consisting of kubeconfig, context and namespace. The environment is stored in `~/.kube/kubectx.json` file. The environment can be stored per terminal, if the environment variable `TERMINAL_ID` exists. This variable can be created e.g. in your `~/.zshrc` with:
````
export TERMINAL_ID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1) 
````

# Installation


## Binary

Go to the [releases](https://github.com/sbueringer/kubectx/releases) page and download the Linux or Windows version. Put the binary to somewhere you want (on UNIX-ish systems, /usr/local/bin or the like). Make sure it has execution bit set. If you want, you can shortcuts to `kcfg`, `kctx` and `kns`, e.g.:

````
function kcfg() { kubectx config    "$@" }
function kctx() { kubectx context   "$@" }
function kns()  { kubectx namespace "$@" }
````
## Auto completion

Auto completion scripts are provided for [zsh](./completion). They can be added or symlinked, e.g. to `/usr/share/zsh/site-functions`.


# Basic usage

kubectx is build with Cobra so the CLI is build in a familiar way (Cobra is also used in Docker and Kubernetes).

To print a description what kubectx can do, just execute:

````
$ kubectx
kubectx manages kubectl context incl. kubeconfig, context and namespace

Usage:
  kubectx [flags]
  kubectx [command]

Available Commands:
  completion  Generates bash completion scripts
  config      Gets and sets the current config
  context     Gets and sets the current context
  help        Help about any command
  namespace   Gets and sets the current namespace

Flags:
  -h, --help   help for kubectx

Use "kubectx [command] --help" for more information about a command.
````
