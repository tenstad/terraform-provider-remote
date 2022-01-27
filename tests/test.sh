mkdir ~/.ssh
cp tests/key ~/.ssh
chmod 600 ~/.ssh/key

eval $(ssh-agent)
ssh-add ~/.ssh/key

go test ./... -v $TESTARGS -timeout 120m
