
mongo 集群配置 [mongo文档](https://www.mongodb.com/docs/manual/core/transactions/#arbiters)
```
mongodb 6.0.x  分片集群方式


注意:
代码中mongodb 事务操作:

mongodb 集群配置server必须符合
1. 各个 shard 每个副本集中不能包含仲裁节点,见上面文档
2. writeConcernMajorityJournalDefault set to true,写关注打开
db.adminCommand({
    "setDefaultRWConcern" : 1,
    "defaultWriteConcern" : {
        "w" : "majority"
    }
})
3. 读关注打开,mongodb 6 以上版本 默认打开
mongodb Client配置
1. 客户端连接配置  写关注  {w: "majority", j: true} 
```

mongo shard主要参数
```
replication:
  replSetName: shard1
  enableMajorityReadConcern: true
...
setParameter:
  maxTransactionLockRequestTimeoutMillis: 36000000
  transactionLifetimeLimitSeconds: 3600
  wiredTigerConcurrentReadTransactions: 512
  wiredTigerConcurrentWriteTransactions: 512
```