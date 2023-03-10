# Steps

1.  Create the migration files for initializing scheme.

        migrate create -ext sql -dir db/migration -seq init_schema

    `/db/migration/000001_init_schema.up.sql` and `/db/migration/000001_init_schema.down.sql` files should have been created.

1.  Implements the files. You might use db diagram.io

## Notes

How to implement query in sqlc for partial update:

```SQL
-- name: UpdateMerchant :one
UPDATE merchants
SET balance = CASE
      WHEN @set_balance::boolean = TRUE THEN @balance
      ELSE balance
    END,
    profession = CASE
      WHEN @set_profession::boolean = TRUE THEN @profession
      ELSE profession
    END,
    title = CASE
      WHEN @set_title::boolean = TRUE THEN @title
      ELSE title
    END,
    about = CASE
      WHEN @set_about::boolean = TRUE THEN @about
      ELSE about
    END,
    image_url = CASE
      WHEN @set_image_url::boolean = TRUE THEN @image_url
      ELSE image_url
    END,
    rating = CASE
      WHEN @set_rating::boolean = TRUE THEN @rating
      ELSE rating
    END
WHERE id = @id
RETURNING *;
```

Other way to do it:

```SQL
-- name: UpdateMerchant :one
UPDATE merchants
SET balance = COALESCE(sqlc.narg(balance), balance),
    profession = COALESCE(sqlc.narg(profession), profession),
    title = COALESCE(sqlc.narg(title), title),
    about = COALESCE(sqlc.narg(about), about),
    image_url = COALESCE(sqlc.narg(image_url), image_url),
    rating = COALESCE(sqlc.narg(rating), rating)
WHERE id = sqlc.arg(id)
RETURNING *;
```

For be sure safe transaction you can use something like this:

```sql
-- name: GetMerchantForUpdate :one
SELECT * FROM merchants
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;
```
