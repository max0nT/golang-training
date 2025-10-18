# Work with golang swagger utility
swag *args:
    docker run --rm -v ./src:/code ghcr.io/swaggo/swag:latest {{args}}

# Update openApi spec
update-spec:
    docker run --rm -v ./src:/code ghcr.io/swaggo/swag:latest init --parseInternal

# Run application
run:
    go run src/main.go
