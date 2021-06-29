#!/bin/bash

#
# recursively render the chart for all fixture value.yaml files beneath the path fiven as an argument
#

FIXTURE_PATH=${1:-"."}
ALL_FIXTURES=`find $FIXTURE_PATH -name values.yaml | xargs -n 1 dirname`

echo "I am running in " `pwd`

for fixture in $ALL_FIXTURES
    do
	echo "Fixture:" $fixture
	mkdir -p $fixture/rejects $fixture/local-render
	helm3 template ./build/package/helm/monoskope -n the-namespace --output-dir $fixture/local-render --values $fixture/values.yaml
    done
