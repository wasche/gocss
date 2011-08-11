#!/bin/bash

function rpad {
  word="$1"
  while [ ${#word} -lt $2 ]; do
    word="$word$3"
  done
  echo -n "$word"
}

white='\033[1;37m'
red='\033[0;31m'
green='\033[0;32m'
reset='\033[0m'

cecho()
{
  local default_msg="No message passed."
  message=${1:-$default_msg}
  color=${2:-$white}
  echo -en "$color"
  echo -n "$message"
  echo -e "$reset"
  return
}

declare -i tests
declare -i failures
tests=0
failures=0

for file in tests/*.css; do
  tests=$tests+1
  rpad $file 40 "."
  ./gocss < $file | diff -q $file.min - >/dev/null
  if [[ $? == 0 ]]; then
    cecho "PASS" $green
  else
    cecho "FAIL" $red
    failures=$failures+1
  fi
done

echo "$tests tests, $failures failures"

