Backend for Kanban
===========================

### How to start

- create config.yml with DB access config

```yaml
db:
  path: db.sqlite
  resetonstart: true
server:
  port: ":3000"
```

- start the backend

```shell script
go build
./kanba-go
```
