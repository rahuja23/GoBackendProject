- Go 1.22
- Docker
- Postgres running on Docker
- Swagger for docs
- Golang migrate for migrations 
- chi - Routing and mux
- airverse-air- For Hot reloading
- direnv - for loading env vars

- migrations for creating and altering table schemas

```aiignore
# To Create a migration file
make migration <Name of the migration file>
# To apply all the migrations
make migrate-up
# To take down the migrations
make migrate-down
```
- makefile to execute migration commands
- CRUD (CREATE READ UPDATE DELETE API HANDLERS)
- Adding go-playground/validator for API request validation