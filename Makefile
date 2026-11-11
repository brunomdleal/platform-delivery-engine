.PHONY: install-cli

install-cli:
	cd cli && go build -o bin/pde main.go && ./bin/pde install