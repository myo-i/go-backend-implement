# go-backend-implement
## Dockerで立てたpostgresに接続できない...
原因は
1. ローカルポート5432が既に使われていたこと
2. コンテナ内のポート5432以外を設定していたこと
https://stackoverflow.com/questions/69462794/docker-and-postgres-server-closed-the-connection-unexpectedly-error-when-using
