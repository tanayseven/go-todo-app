Go TODO App
===========

A simple TODO app using Go, HTMX, SQLITE
----------------------------------------

Setup
-----

1. Install all the dependencies `go mod tidy`
2. Copy the initial database `cp main.init.db main.db`
3. Run the project `go run .`
4. Open the browser and go to `http://localhost:8080`
5. Generate a build for production `go build .`


TODO
----
- [X] Setup Admin
- [X] Setup Templating
- [X] Setup Logging
- [X] Add GORM support to store and retrieve the data
- [X] Implement Add
- [X] Implement Delete
- [X] Implement Update
- [ ] Implement Mark Done
- [ ] Implement Mark Undone
- [ ] Write the frontend for the TODO List
- [ ] Write simple tests
- [ ] Setup user auth
- [ ] Add HTMX for better UX
