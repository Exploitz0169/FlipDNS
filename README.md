
## Migrations locally
See: https://github.com/golang-migrate/migrate/tree/master

```
brew install golang-migrate
brew install sqlc
```

```
export POSTGRES_URL="postgresql://flip:postgres@localhost:5432/flipdns?sslmode=disable"

# Creating migration
migrate create -ext sql -dir migrations -seq <migration_name>

# Running migration
migrate -database ${POSTGRES_URL} -path migrations <up|down>

```

Then run sqlc to keep repository up-to-date:

```
sqlc generate
```
