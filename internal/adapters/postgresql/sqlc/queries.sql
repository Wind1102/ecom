-- name: ListProducts :many
SELECT * FROM products;

-- name: FindProductById :one
SELECT * from products where id = $1;


