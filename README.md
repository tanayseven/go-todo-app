Go TODO App
===========

A simple TODO app using Go, HTMX, SQLITE, Tailwind
----------------------------------------

Setup
-----

1. Install Tailwind cli from here: https://tailwindcss.com/blog/standalone-cli
2. Install all the dependencies `go mod tidy`
3. Copy the initial database `cp main.init.db main.db`
4. Run the project `go run .`
5. Open the browser and go to `http://localhost:8080`
6. Generate a build for production `go build .`


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
