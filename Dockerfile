FROM centurylink/ca-certs
WORKDIR /app
COPY ./kanban-go /app
COPY ./uploads /app

CMD ["/app/kanban-go"]