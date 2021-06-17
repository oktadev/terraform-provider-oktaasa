#!/usr/bin/env bash

# Check for missing error handlers
echo "==> Checking for unchecked errors..."

if ! which errcheck > /dev/null; then
    echo "==> Installing errcheck..."
    go get -u github.com/kisielk/errcheck
fi

err_files=$(errcheck ./...)

if [[ -n ${err_files} ]]; then
    echo 'Unchecked errors found in the following places:'
    echo "${err_files}"
    echo "Please handle returned errors. You can check directly with \`make errcheck\`"
    exit 1
fi

exit 0
