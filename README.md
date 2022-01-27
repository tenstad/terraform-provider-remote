# Terraform Provider Remote

**Documentation**: https://registry.terraform.io/providers/tenstad/remote/latest/docs/resources/file

## Issues

Please create an issue if you find typos, have any questions, or problems with
the provider. Feature requests are welcome, but you may have to get your hands
dirty and take part in the implementation to get the feature merged.

## Pull Requests

Pull requests are welcome! 😊 You may create an issue to discuss the changes
first, but it is not a strict requirement.

Before creating a PR, please ensure that all documentation is correctly
generated (`go generate`), and that go modules are tidy (`go mod tidy`).

## Developing the Provider

### Requirements

If you wish to work on the provider, you'll first need Go installed on your
machine. You might also want Terraform and Docker, depending on the work.

- [Go](https://golang.org/doc/install) >= 1.16
- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Docker](https://www.docker.com/get-started) (for test purposes)

### Workflow

- `go install` to compile the provider. The provider binary will be built in
  `$GOPATH/bin`.
- `make test` to acceptance test the provider. NOTE: Docker is required, tests
  are performed between container.
- `go generate` to generate documentation.
