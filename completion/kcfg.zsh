#compdef kcfg

PREV=""
_arguments "1: :($(ls ${HOME}/.kube/config* | sort -n))"
