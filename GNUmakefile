default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	rm tests/key tests/key.pub
	ssh-keygen -f tests/key -N ""
	docker build -t remotefile-host tests
	docker rm -f remotefile-host
	docker run --rm -d -p 8022:22 --name remotefile-host remotefile-host
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
	docker rm -f remotefile-host
	rm tests/key tests/key.pub
