#compdef kns

local KUBECONFIG=$(kcfg)
_arguments "1: :($(kubectl --kubeconfig $KUBECONFIG get namespaces -o=jsonpath='{range .items[*].metadata.name}{@}{"\n"}{end}'))"
