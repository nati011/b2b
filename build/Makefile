# Run the application
.PHONY: run
run:
	clear
	go run ../cmd/main.go

# Re-Run the application
.PHONY: rer
rer:
	sudo kill -9 $$(sudo lsof -t -i :8000)
	make run

# Run all tests
.PHONY: test
test:
	go test -p 1 ../cmd
