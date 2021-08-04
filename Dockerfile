FROM debian:10-slim
WORKDIR /app
ADD ./kanban-go /app
ADD ./uploads /app/uploads

CMD ["/app/kanban-go"]