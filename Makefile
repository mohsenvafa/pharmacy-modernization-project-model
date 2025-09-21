dev:
	templ generate -watch 		-proxyport=7332 		-proxy="http://localhost:8080" 		-cmd="go run ./cmd/server" 		-open-browser=false & 	npx tailwindcss -i ./web/styles/input.css -o ./web/public/app.css --watch

css:
	npx tailwindcss -i ./web/styles/input.css -o ./web/public/app.css

css\:watch:
	npx tailwindcss -i ./web/styles/input.css -o ./web/public/app.css --watch
