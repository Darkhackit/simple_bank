-- name: CreateUser :one
INSERT INTO users (
     email, username, active, name, password
) VALUES ($1 ,$2 , $3 ,$4,$5) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE username = $1 LIMIT 1;

-- name: UpdateUser :one
UPDATE users SET
                 password = COALESCE(sqlc.narg(password),password),
                 name = COALESCE(sqlc.narg(name),name),
                 email = COALESCE(sqlc.narg(email),email)
                 WHERE
                 username = sqlc.arg(username)
             RETURNING *;
