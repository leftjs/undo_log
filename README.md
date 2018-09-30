# undo_log

数据库 undo log 模拟实现

# 题目

有一个交易系统，存储用户的信息，用户之间可以进行转账。转账时直接对账户信息进行修改，如果转账有问题（例如 Tom 给 Jerry 转 10 块钱，但是 Tom 并没有 10 块钱），则对该交易进行回滚。同时该交易系统存在安全问题，可能被恶意攻击，需要对一批交易进行回滚。

# 要求

- 使用 undo 日志实现交易回滚的功能
- 自定义 undo 日志的格式
- 当交易存在问题时，立即回滚交易
- 可以回滚某个交易之后的所有交易
- 可能同时会有多次交易并行运行
- undo 日志不需要保留太久，历史的 undo 日志定期删除

# 提示

- 注意代码可读性，添加必要的注释
- 注意代码风格与规范，添加必要的单元测试和文档
- 注意异常处理，尝试优化性能

# 项目说明

## 目录结构说明

- src 为程序主目录，包括如下内容

  1. `file` :为文件读写模块，用于 db 以及 log 的相关操作
  2. `db` :数据存储模块
  3. `logs`:日志存储模块
  4. `config`:系统配置模块
  5. `transaction`:事务管理模块
  6. `util`: 存储一些用到的工具类

- log 为 undo log 日志存储目录，已在 ignore 中忽略，github 上看不到，程序执行时会生成，下同
- data 为 用户数据存放目录
- test 为测试数据存放目录

## 日志格式说明

```bash
<T1 START> # T1, T2, T....  事务 id 自增长，开启一个事务，下同
<T1,1,2,3,4> # 其中1表示 fromId, 2表示from user的原Cash,3表示 toId, 4表示to user的原Cash
<T1 COMMIT> # T1事务提交
<T2 UNDO> # T1事务回滚
```

## 执行顺序

## 功能说明

1. 用户添加功能，内存与文件系统同步，先写文件系统后写内存，可并发执行，具体查看测试用例和代码
2. undo log 功能，undo log 的日志的写入粒度拆分的很细，主要是为了用户宕机恢复，测试用例中没有做宕机模拟，可在启动时调用`triggerUndo` 做一次 undo check，将没有完成的事务进行回滚
3. undo log 写入的时候同样会内存与文件系统同步，先写文件后写内存，可并发执行
4. undo log 的在内存数据结构使用 hashmap+双向链表，保证日志顺序的同时，通过读取 Tail 节点可用于快速判断当前是否已经结束，是否可 put，或者在 undo check 时候检查是否需要 rollback
5. 启动时会 scan 出最新的 log file， 读取后构建内存 undo log
6. log 匹配，日志文件以`[timestamp]_[random_string].log`命名,用于实现时间窗口，用于 gc 时删除指定时间区间的日志，加入 random_string 主要是防止 gc 过快，产生重名 log 文件
7. 事务功能，事务 id 为从 1 开始的自增，考虑到事务并发过多，可能超过 int 的表示，可定时进行日志的 gc 操作来重置事务 id
8. 事务请求中如果存在 `cash < 0` 或者 transfer 之后导致用户金额会零，会直接进行回滚，回滚成功写`<T%d UNDO>`
9. 调用`db.RollbackAfter(int)`会先做一次 undo check，回滚为提交事务，然后在进行已经 commit 的事务的回滚操作
10. gc 时候会进行一个 undo check，检查出所有正在进行的事务，等待事务完成后在进行 gc 操作，
11. gc 操作会重新生成当前时间的 log file 用于 log 的写入，同时 gc 之后会重置事务 ID

## TODO

1. 用户数据定期与文件系统 sync
2. 定时滑窗清理日志
3. rollback 某个事务 id 之后时候后的 sync
4. 一转多，多转一类型的 转账事务的支持
5. ....

# 输出内容

```bash
go main.go # 输出内容如下

2018/09/30 16:41:54 ------before transaction------
2018/09/30 16:41:54 [ID: 1] - Tom has 10 money
2018/09/30 16:41:54 [ID: 2] - Jerry has 10 money
2018/09/30 16:41:54 [ID: 3] - Spike has 10 money
2018/09/30 16:41:54 ------before transaction------
2018/09/30 16:41:54
2018/09/30 16:41:54 do parallel transactions, if error occurs will rollback!!!!
2018/09/30 16:41:54 after transfer from user cash < 0
2018/09/30 16:41:54 transaction: T4 need to undo
2018/09/30 16:41:54
2018/09/30 16:41:54 ------after transaction------
2018/09/30 16:41:54 [ID: 1] - Tom has 10 money
2018/09/30 16:41:54 [ID: 2] - Jerry has 5 money
2018/09/30 16:41:54 [ID: 3] - Spike has 15 money
2018/09/30 16:41:54 ------after transaction------
```

由于 transaction 为并发调用，故不能保证 transaction 的执行顺序，因此可能结果会不同
