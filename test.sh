#!/bin/bash
set -x
set -e
#clear

verbose=""

echo $@

if [[ "$@" == "-v" ]]
then
  verbose="-v"
fi


go build -o tsoda ./soda

function test {
  echo "Testing $1"
  export SODA_DIALECT=$1
  echo ./tsoda -v
  ! ./tsoda drop -e $SODA_DIALECT -c ./database.yml
  ! ./tsoda create -d -e $SODA_DIALECT -c ./database.yml
  ./tsoda migrate -d -e $SODA_DIALECT -c ./database.yml -d
  go test ./... $verbose
}

#test "postgres"
#test "sqlite"
#test "mysql"
test 'sqlserver'

rm tsoda
