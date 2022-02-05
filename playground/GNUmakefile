default: install

# Start host containers
hosts: clean
	cd .. && $(MAKE) hosts

# Stop containers
clean:
	rm -f *tfstate*
	cd .. && $(MAKE) clean

# Install provider
install:
	cd .. && $(MAKE) install
