.PHONY: make
build:
	@go mod tidy
	@go generate && go build .
run:
	@go mod tidy
	@go generate && go run main.go
docker:
	@docker compose down
	@docker compose up -d

text ?= defalut git commit message
git:
	@echo "Adding changes to the staging area..."
	@git add .
	@echo "Committing changes with message: $(text)"
	@git commit -m "$(text)"
	@echo "Pushing changes to remote repository..."
	@git push