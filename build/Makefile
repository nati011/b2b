# Run the application
.PHONY: run
run:
	clear
	go run ../cmd

# Re-Run the application
.PHONY: rer
rer:
	sudo kill -9 $$(sudo lsof -t -i :4000)
	make run