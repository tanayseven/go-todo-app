FROM debian:bookworm-slim
WORKDIR /app
COPY ./go-todo-app /app/go-todo-app
COPY static /app/static
COPY templates /app/templates
EXPOSE 9033
RUN ls -l /app
CMD ["/app/go-todo-app"]
