FROM alpine
WORKDIR /app
COPY go-todo-app /app/go-todo-app
COPY static /app/static
COPY templates /app/templates
EXPOSE 9033
CMD ["ls -l /app/ &&", "./go-todo-app"]
