
# Kubectx

Command line utility to manage the current environment consisting of kubeconfig, context and namespace. The environment is stored in `~/.kube/kubectx.json` file. The environment can be stored per terminal, if the environment variable `TERMINAL_ID` exists. This variable can be created e.g. in your `~/.zshrc` with:
````
export TERMINAL_ID=$(cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 10 | head -n 1) 
````
