darwin: ## Build for macOS.
	env GOOS=darwin GOARCH=amd64 go build -o bin/moul-darwin -i github.com/sophearak/moul

linux: ## Build for Linux.
	env GOOS=linux GOARCH=amd64 go build -o bin/moul-darwin -i github.com/sophearak/moul
