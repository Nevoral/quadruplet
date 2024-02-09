# Define a default target
.PHONY: default
default: benzo

# Define the target to run 'templ generate'
.PHONY: robot
benzo:
	templ generate
	go run cmd/QuadrupletWeb/main.go
leg:
	templ generate
	go run cmd/LegWeb/main.go