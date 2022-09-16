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
## DB Transaction
### What is a db transaction?
* A single unit of work.
* Often made up of multiple db operations.
### Example
* Create a transfer record with amount = 10.
* Create an account entry for account1 with account = -10.
* Create an account entry for account2 with account = +10.
* Subtract 10 from the balance of account1.
* Add 10 to the balance of account2.
### Why do we need db transaction?
* To provide a reliable and consistent unit of work, even in case of system failure.
* To provide a isolation between programs that access the database concurrently.
### ACID Property
* Atomicity(A)  
Either all operations complete successfully or the transaction fails and the db is unchanged.
* Consistency(C)  
The db state must be valid after the transaction. All constraints must be satisfied.
* Isolation(I)  
Concurrent transactions must not affect each other.
* Durability(D)  
Data written by a successful transaction must be recorded in persistent storage.
### How to run SQL tx?
```
begin;
...
commit;

begin;
...
rollback;
```