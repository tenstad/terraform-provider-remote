ssh-keygen -f "$HOME/.ssh/known_hosts" -R "[localhost]:8022"
docker build -t terraform-provider-remotefile-test .
docker run --rm -d -p 8022:22 --name terraform-provider-remotefile-test terraform-provider-remotefile-test
