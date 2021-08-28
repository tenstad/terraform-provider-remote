default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	docker rm -f remotehost
	docker network rm remote || true
	docker network create remote
	docker build -t remotehost tests
	docker run --rm -d --net remote --name remotehost remotehost
	docker run --rm --net remote -v $(PWD):/app --workdir /app -e "TF_ACC=1" -e "TF_ACC_TERRAFORM_VERSION=0.13.4" golang:1.15 go test ./... -v $(TESTARGS) -timeout 120m
	docker rm -f remotehost
	docker network rm remote
