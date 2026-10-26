# Protocol configuration

Only `http` runs today. An `https://` URL is the HTTPS API case, and the CLI already accepts it on extract and load.

The other sections are the configuration to save when those protocols are added. They are not written by the CLI yet. `Source` and `Destination` currently store `name`, `protocol`, `url`, and `method` (`contentType` on a load). Fields such as `database` and `collection` are the intended shape, not fields the menu can save.

Examples use placeholders. Do not commit real passwords.

## http

Reads or writes an HTTP API, including HTTPS.

Extract fields: `url`, `method` (`GET`).

Load fields: `url`, `method` (`POST`), `contentType`.

Later package: the standard library `net/http`. This one is already in use.

```json
{
  "name": "docs",
  "protocol": "http",
  "url": "https://example.com/data.json",
  "method": "GET"
}
```

```json
{
  "name": "hook",
  "protocol": "http",
  "url": "https://example.com/hook",
  "method": "POST",
  "contentType": "text/plain"
}
```

## mongodb

Reads or writes one collection.

Extract and load fields: `uri`, `database`, `collection`. A load also needs the document or update to write. That payload is the body from the previous step, not a field on the destination.

Later package: `go.mongodb.org/mongo-driver`.

```json
{
  "name": "users",
  "protocol": "mongodb",
  "uri": "mongodb://user:password@localhost:27017",
  "database": "app",
  "collection": "users"
}
```

## postgres

Reads or writes rows with SQL.

Extract fields: `dsn`, or `host`, `port`, `user`, `password`, `database`, and `sslmode`, plus the `query`.

Load fields: the same connection fields, plus the `statement`. Use a placeholder for the body (`$1`) instead of pasting the body into the SQL text.

`sslmode` is commonly `require` or `disable`. Prefer `require` outside a local machine.

Later package: `github.com/jackc/pgx/v5`.

```json
{
  "name": "orders",
  "protocol": "postgres",
  "dsn": "postgres://user:password@localhost:5432/app?sslmode=disable",
  "query": "select id, total from orders"
}
```

## mysql

Reads or writes rows with SQL. The DSN form is `user:password@tcp(host:3306)/database`.

Extract fields: `dsn` and `query`.

Load fields: `dsn` and `statement`. Use `?` as the placeholder for the body.

Later package: `github.com/go-sql-driver/mysql` with `database/sql`.

```json
{
  "name": "orders",
  "protocol": "mysql",
  "dsn": "user:password@tcp(localhost:3306)/app",
  "query": "select id, total from orders"
}
```

## sqlite

Reads or writes a database file on disk. There is no host.

Extract fields: `path` and `query`.

Load fields: `path` and `statement`. Use `?` as the placeholder for the body.

Later package: `modernc.org/sqlite` with `database/sql`. That driver is pure Go.

```json
{
  "name": "local",
  "protocol": "sqlite",
  "path": "data/app.db",
  "query": "select id, total from orders"
}
```

## redis

Reads or writes a key.

Extract fields: `addr`, optional `password`, optional `db`, and `key`.

Load fields: the same connection fields and `key`. The value is the body from the previous step.

Later package: `github.com/redis/go-redis/v9`.

```json
{
  "name": "cache",
  "protocol": "redis",
  "addr": "localhost:6379",
  "password": "",
  "db": "0",
  "key": "last-body"
}
```

## s3

Reads or writes one object.

Extract and load fields: `bucket`, `region`, and `key`. Take credentials from the environment (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, and optionally `AWS_SESSION_TOKEN`). Leave them out of the JSON file.

Later package: `github.com/aws/aws-sdk-go-v2`.

```json
{
  "name": "archive",
  "protocol": "s3",
  "bucket": "example-bucket",
  "region": "us-east-1",
  "key": "exports/last-body.txt"
}
```
