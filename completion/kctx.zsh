#compdef kctx

local KUBECONFIG=$(kcfg)
_arguments "1: :($(kubectl --kubeconfig $KUBECONFIG config get-contexts --output='name'))"
