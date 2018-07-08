#compdef kubectx

_arguments \
  '1: :->level1' \
  '2: :_files'
case $state in
  level1)
    case $words[1] in
      kubectx)
        _arguments '1: :(completion config context help namespace)'
      ;;
      *)
        _arguments '*: :_files'
      ;;
    esac
  ;;
  *)
    _arguments '*: :_files'
  ;;
esac
