#!/usr/bin/env bash

git clone https://github.com/palette-software/deploy-scripts.git

if [ "${CURRENT_ENV_DIR}" == "linux_amd64" ]; then
    DEPLOY_FILE="${TRAVIS_BUILD_DIR}/dbcheck/_build/*" deploy-scripts/rpm/deploy-to-rpm.sh
fi

# The existence of GITHUB_TOKEN variable is the switch for Github release
if [[ -n $GITHUB_TOKEN ]]; then

	pushd deploy-scripts/github
	# This script is expected to be executed from its folder
	./release-to-github.sh
	popd
fi
