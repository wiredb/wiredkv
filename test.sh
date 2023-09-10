#!/bin/bash

echo "VaseDB Testing Shell Script."

echo ""

echo "1: testing cmd package."
echo "2: testing conf package."
echo "3: testing server package."
echo "4: testing codecov coverage."

echo ""

case_num=$1

if [ -z "$case_num" ]; then
    echo "Please provide an option (1, 2, or 3)."
    exit 1
fi

rm -rf ./_temp
mkdir -p ./_temp

if [ "$case_num" -eq 1 ]; then
    cd cmd && go test -c && sudo -S ./cmd.test
elif [ "$case_num" -eq 2 ]; then
    echo "Testing conf package"
elif [ "$case_num" -eq 3 ]; then
    echo "Testing server package"
elif [ "$case_num" -eq 4 ]; then
    go test ./... -race -coverprofile=coverage.out -covermode=atomic -v
else
    echo "Invalid option. Please provide a valid option (1, 2, or 3)."
fi
