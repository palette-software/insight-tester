#!/bin/bash

# Fail if there are errors
set -e

SANITY_CHECK_INSTALL_DIR=/opt/insight-sanity-check
LOCKFILE=/tmp/sanity_check.flock

flock -n ${LOCKFILE} \
	${SANITY_CHECK_INSTALL_DIR}/dbcheck ${SANITY_CHECK_INSTALL_DIR}/tests/sanity_checks.yml ${SANITY_CHECK_INSTALL_DIR}/Config.yml > /dev/null
