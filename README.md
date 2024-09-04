# The Recipe Book

Start the database and cache:
`docker-compose up -d`

Exec into postgres DB:
`docker-compose exec postgres psql -U admin the_recipe_book`

## Migrations

1. `brew install golang-migrate`
2. `migrate create -ext sql -dir migrations -seq ${migration_name}`