default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	docker rm -f remotehost
	docker network rm remotefile || true
	docker network create remotefile
	docker build -t remotehost tests
	docker run --rm -d --net remotefile --name remotehost remotehost
	docker run --rm --net remotefile -v $(PWD):/app --workdir /app -e "TF_ACC=1" -e "TF_ACC_TERRAFORM_VERSION=0.13.4" golang:1.15 go test ./... -v $(TESTARGS) -timeout 120m
	docker rm -f remotehost
	docker network rm remotefile
