# rabbitmqExample

l1 Message publisher and consumeer without ack and durable    简单示例  

l2 Message publisher and consumer example with message ack and durable  消息持久化

## 持久化

RabbitMQ 持久化包含3个部分

- exchange 持久化，在声明时指定 durable 为 true
- queue 持久化，在声明时指定 durable 为 true
- message 持久化，在投递时指定 delivery_mode=2（1是非持久化）

queue 的持久化能保证本身的元数据不会因异常而丢失，但是不能保证内部的 message 不会丢失。要确保 message 不丢失，还需要将 message 也持久化

如果 exchange 和 queue 都是持久化的，那么它们之间的 binding 也是持久化的。

如果 exchange 和 queue 两者之间有一个持久化，一个非持久化，就不允许建立绑定。

l3 Example of using fanout exchange

l4 Example of using direct exchange

l5 Example of using topic exchange

​	direct类型差不多，但direct类型要求routingkey完全相等，这里的routingkey可以有通配符：'*','#'.

l6 Example of RPC,流量控制QOS

client 发送一个消息，server 消费后（RPC得到结果），server在发送结果消息，client 消费结果消息