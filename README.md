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
## Transaction Isolation Level（事务隔离等级）
所有数据库都必须满足ACID性质。`Isolation(I)`就是其中的一条性质。
### Read Phenomena（读现象）
* Dirty Read（脏读）  
A transaction reads data written by other concurrent uncommitted transaction.
* Non-repeatable Read（不可重复读）  
A transaction reads the same row twice and sees different value because it has been modified by other committed transaction.
* Phantom Read（幻读）  
A transaction re-executes q query to find rows that satisfy a condition and sees a different set of rows, due to changes by other committed transaction.
* Serialization Anomaly（序列化异常）  
The result of a group of concurrent committed transactions is impossible to achieve if we try to run them sequentially in any order without overlapping.
### 4 Standard Isolation Levels
* Read Uncommitted  
Can see data written by uncommitted transaction.
* Read Committed  
Only see data written by committed transaction.
* Repeatable Read  
Same read query always returns same result.
* Serializable  
Can achieve same result if execute transactions serially in some order instead of concurrency.
### Isolation Levels in MySQL
|                       | Read Uncommitted | Read Committed | Repeatable Read | Serializable |
|-----------------------|------------------|----------------|-----------------|--------------|
| Dirty Read            | ️✅               | ❌              | ❌               | ❌            |
| Non-repeatable Read   | ️✅               | ✅              | ❌               | ❌            |
| Phantom Read          | ️✅               | ✅              | ❌               | ❌            |
| Serialization Anomaly | ️✅               | ✅              | ️✅              | ❌            |
### Isolation Levels in PostgresSQL
|                       | Read Uncommitted | Read Committed | Repeatable Read | Serializable |
|-----------------------|------------------|----------------|-----------------|--------------|
| Dirty Read            | ️❌               | ❌              | ❌               | ❌            |
| Non-repeatable Read   | ️✅               | ✅              | ❌               | ❌            |
| Phantom Read          | ️✅               | ✅              | ❌               | ❌            |
| Serialization Anomaly | ️✅               | ✅              | ️✅              | ❌            |

## Git Action
### Workflow
![avatar](./doc/images/workflow.png)
* Is an automated procedure.
* Made up of 1+ jobs.
* Triggered by events, scheduled, or manually.
* Add `.yml` file to repository.
#### Example
```yaml
name: build-and-test # workflow名字

on: # 触发器
  push:
    # 向main push时触发
    branches: [main] 
  schedule:
    # 每15分钟触发
    - cron: '*/15 * * * *' 

# 工作
jobs:
  build:
    runs-on: ubuntu-latest
```
### Runner
![avatar](./doc/images/runner.png)
* Is a server to run the jobs.
* Run 1 job at a time.
* GitHub hosted or self-hosted.
* Report progress, logs & result to github.
### Job
![avatar](./doc/images/job.png)
* Is a set of steps execute on the same runner.
* Normal jobs run in parallel.
* Dependent jobs run serially.
#### Example
```yaml
jobs:
  # build job
  build: 
    runs-on: ubuntu-latest
    steps:
      - name: Check out code
        uses: actions/checkout@v2
      - name: Build server
        run: ./build_server.sh
  # test job
  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - run: ./test_server.sh
```
### Step And Action
![avatar](./doc/images/step-action.png)
#### Step
* Is an individual task.
* Run serially within a job.
* Contain 1+ actions.  
#### Action
* Is a standalone command.
* Run serially within a step.
* Can be reused.
#### Example
```yaml
jobs:
  build: 
    runs-on: ubuntu-latest
    steps:
      # action1
      - name: Check out code
        uses: actions/checkout@v2
      # action2
      - name: Build server
        run: ./build_server.sh
```
### Summary
![avatar](./doc/images/summary.png)
* Ways of triggering a workflow: event, scheduled or manually.
* A workflow consists of one or multiple jobs. A job is composed of multiple steps. Each step can have 1 or more actions.
* All jobs inside a workflow normally run in parallel, unless they depend on each other. Then in that case, they run serially.
* Each job will be run separately by a specific runner. The runner will report progress, logs, and results of the jobs back to GitHub.

## Load Env Vars
### Why file?
For development, easily specify default configuration for local development and testing.
### Why env vars?
For deployment, easily override the default configurations when deploy with docker containers.
### Why viper?
* Find, load, unmarshal config file.  
`json, toml, yaml, env, ini`
* Read config from environment variables or flags.  
Override existing values, set default values.
* Read config from remote system.  
Etcd, Consul.
* Live watching and writing config file.  
Reread changed file, save any modifications.
### Notice
如果环境变量和配置文件都配置两，viper会使用环境变量的值（即：环境变量的值会覆盖从配置文件中读到的值）。
