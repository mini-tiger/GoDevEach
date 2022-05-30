policy是策略或者说是规则的定义。它定义了具体的规则。

request是对访问请求的抽象，它与e.Enforce()函数的参数是一一对应的

matcher匹配器会将请求与定义的每个policy一一匹配，生成多个匹配结果。

effect根据对请求运用匹配器得出的所有结果进行汇总，来决定该请求是允许还是拒绝。

PERM模型
Policy: 定义权限的规则
Effect: 定义组合了多个 Policy 之后的结果, allow/deny
Request: 访问请求, 也就是谁想操作什么
Matcher: 判断 Request 是否满足 Policy

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act
p2 = sub, act # p2 定义的是 sub 所有的资源都能执行 act

[matchers]
m = r.sub == p.sub && r.obj == p.obj && r.act == p.act

[policy_effect]
e = some(where (p.eft == allow))
```
[https://casbin.org/docs/en/syntax-for-models](https://casbin.org/docs/zh-CN/syntax-for-models)

上面模型文件规定了权限由sub,obj,act三要素组成，只有在策略列表中有和它完全相同的策略时，该请求才能通过。
匹配器的结果可以通过p.eft获取，some(where (p.eft == allow))表示只要有一条策略允许即可。

```
[role_definition]
g = _, _

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```
g = _,_定义了用户——角色，角色——角色的映射关系，前者是后者的成员，拥有后者的权限。然后在匹配器中，
我们不需要判断r.sub与p.sub完全相等，只需要使用g(r.sub, p.sub)来判断请求主体r.sub是否属于p.sub这个角色即可。