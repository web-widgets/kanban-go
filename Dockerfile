FROM centurylink/ca-certs
WORKDIR /app
COPY ./kanban-go /app

CMD ["/app/kanban-go"]