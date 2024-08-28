FROM alpine
WORKDIR /app
COPY go-todo-app /app/go-todo-app
COPY static /app/static
COPY templates /app/templates
EXPOSE 9033
RUN ls -l
RUN ls -l /app/
CMD ["./go-todo-app"]
