# Backend-in-go
## Docker
```
# 拉取镜像
docker pull postgres:14.5
# 启动
docker run --name bank-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=csq2400306 -d postgres:14.5
# 进入命令行
docker exec -it bank-postgres psql -U root
# 查看日志
docker logs bank-postgres
# 停止
docker stop bank-postgres
# 启动
docker start bank-postgres
```
## Migration（数据库迁移）
```
migrate create -ext sql -dir db/migration -seq init_schema
# up
migrate -path db/migration -database "postgresql://root:csq2400306@localhost:5432/bank?sslmode=disable" -verbose up
# down
migrate -path db/migration -database "postgresql://root:csq2400306@localhost:5432/bank?sslmode=disable" -verbose down
```
## CRUD
* Create: insert new records to the database.
* Read: select or search for records in the database.
* Update: change some fields of the record in the database.
* Delete: remove records from the database.
## CRUD的几种方式
### Database/SQL
* Very fast & straightforward.
* Manual mapping SQL fields to variables
* Easy to make mistakes, not caught until runtime.
### GORM
* CRUD functions already implemented, vary short production code.
* Must learn to write queries using gorm's function.
* Run slowly on high load.
### SQLX
* Quite fast & east to use.
* Fields mapping via query text & struct tags.
* Failure won't occur until runtime.
### SQLC
* Very fast & easy to use.
* Automatic code generation.
* Catch SQL query errors before generating codes.

