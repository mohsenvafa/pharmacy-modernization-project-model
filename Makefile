dev:
	templ generate -watch \
		-proxyport=7332 \
		-proxy="http://localhost:8080" \
		-cmd="go run ./cmd/server" \
		-open-browser=false

mock-iris:
	go run ./cmd/iris_mock
