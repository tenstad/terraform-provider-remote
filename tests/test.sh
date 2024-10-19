#!/bin/sh

eval "$(ssh-agent)"
chmod 0600 tests/key
ssh-add tests/key

ssh-keygen -f ~/.ssh/known_hosts -R remotehost
ssh-keygen -f ~/.ssh/known_hosts -R remotehost2

export SKIP_TEST_MOVE_FILE_BY_MODIFYING_PROVIDER=1

# shellcheck disable=SC2086
go test ./... -v $TESTARGS -timeout 120m
