# Steps

1. Create the migration files for initializing scheme.

        migrate create -ext sql -dir db/migration -seq init_schema

    `/db/migration/000001_init_schema.up.sql` and `/db/migration/000001_init_schema.down.sql` files should have been created.

1. Implements the files. You might use db diagram.io
