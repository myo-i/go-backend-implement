-- name: CreateAccount :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- もともとはSET balance = balance + $2と記載していたが、複数のamountが足された値の為、名前がややこしい
-- だからカラム名はbalanceのまま、変換した後の名前を変えたくてsqlc.arg(amount)に変更した（たぶん）
-- account1.Balance - arg.Amountやaccount2.Balance + arg.Amountという記述は簡略化できるので下記に書き直した
-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;


-- name: DleteAccount :exec
DELETE FROM accounts 
WHERE id = $1;