go test -v $(go list -f '{{.Dir}}/...' -m | xargs)
