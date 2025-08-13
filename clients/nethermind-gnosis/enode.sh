#!/bin/bash

# Script to retrieve the enode
#
# This is copied into the validator container by Hive
# and used to provide a client-specific enode id retriever
#

# Immediately abort the script on any error encountered


set -e

TARGET_ENODE=$(
  sed -n -e 's/^.*This node.*: //p' /log.txt \
  | LC_ALL=C sed -E 's/\x1b\[[0-9;]*[a-zA-Z]//g' \
  | iconv -c -t UTF-8
)
echo "${TARGET_ENODE/|/}"
