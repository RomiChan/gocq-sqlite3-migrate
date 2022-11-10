# gocq-sqlite3-migrate
go-cqhttp leveldb v3 到 sqlite3 迁移工具

## 安装

你可以下载 [已经编译好的二进制文件](https://github.com/RomiChan/gocq-sqlite3-migrate/releases).

从源码安装:
```bash
$ go install github.com/RomiChan/gocq-sqlite3-migrate@latest
```

## 使用方法

```bash
./gocq-sqlite3-migrate -from xxx -to yyy
```
默认值：
 * from: `data/leveldb-v3`
 * to: `data/sqlite3/msg.db`
