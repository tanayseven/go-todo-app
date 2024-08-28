Go TODO App
===========

A simple TODO app using Go, HTMX, SQLITE, Tailwind
----------------------------------------

Setup
-----

1. Install Tailwind cli from here: https://tailwindcss.com/blog/standalone-cli
2. Install all the dependencies `go mod tidy`
3. go install github.com/air-verse/air@latest
4. Run this command while editing styles `./tailwindcss -i main.css -o ./static/main.css --watch`
5. Run the project `air --build.cmd "go build -o bin/api ." --build.bin "./bin/api"`
6. Open the browser and go to `http://localhost:9033`
7. Generate a build for production `go build .`


TODO
----
- [X] Setup Admin
- [X] Setup Templating
- [X] Setup Logging
- [X] Add GORM support to store and retrieve the data
- [X] Implement Add
- [X] Implement Delete
- [X] Implement Update
- [X] Implement Mark Done
- [X] Implement Mark Undone
- [X] Write the frontend for the TODO List
- [X] Add HTMX for better UX
- [ ] Write simple tests
- [ ] Setup user auth
