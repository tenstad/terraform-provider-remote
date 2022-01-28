mkdir ~/.ssh
cp tests/key ~/.ssh
chmod 600 ~/.ssh/key

eval $(ssh-agent)
ssh-add ~/.ssh/key

export SKIP_TEST_MOVE_FILE_BY_MODIFYING_PROVIDER=1

go test ./... -v $TESTARGS -timeout 120m
