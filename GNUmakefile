default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	docker rm -f remotehost
	docker build -t remotehost tests
	docker run --rm -d -p 8022:22 --name remotehost remotehost
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
	docker rm -f remotehost
