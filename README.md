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
## Migration
```
migrate create -ext sql -dir db/migration -seq init_schema
# up
migrate -path db/migration -database "postgresql://root:csq2400306@localhost:5432/bank?sslmode=disable" -verbose up
# down
migrate -path db/migration -database "postgresql://root:csq2400306@localhost:5432/bank?sslmode=disable" -verbose down
```
