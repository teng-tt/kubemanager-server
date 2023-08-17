# 亲和调度

亲和调度是在生产实践环节用的最为广泛的调度方式。

## node亲和调度

亲和调度方式看起来比我们前面学习的nodeName和nodeSelector较为复杂，但是仔细梳理，其实发现非常有规律，也非常便于理解。

### node硬亲和
硬亲和，"硬" 即强硬的意思，是一种强制选择策略，下面是Pod的写法。

```yaml
 affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: app
            operator: In
            values:
            - test
```

requiredDuringSchedulingIgnoredDuringExecution说明：
- 前半段requiredDuringScheduling表示下面定义的规则必须强制满足（require）才会调度Pod到节点上
- 后半段IgnoredDuringExecution表示已经在节点上运行的Pod不需要满足下面定义的规则，即去除节点上的某个标签，那些需要节点包含该标签的Pod不会被重新调度

operator可选值：
- NotIn：标签的值不在某个列表中
- Exists：某个标签存在
- DoesNotExist：某个标签不存在
- Gt：标签的值大于某个值（字符串比较）
- Lt：标签的值小于某个值（字符串比较）


### node软亲和
node软亲和，相比于上一步的硬亲和，是一种优先选择策略，权重越高，则越先选中。
```yaml
affinity:
    nodeAffinity:
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 80 
        preference: 
          matchExpressions: 
          - key: app
            operator: In 
            values: 
            - test
      - weight: 20 
        preference: 
          matchExpressions: 
          - key: app
            operator: In 
            values: 
            - test2
```
preferredDuringScheduling

IgnoredDuringExecution，表示会根据规则优先选择哪些节点。

## Pod 亲和调度

podA - > node01

podC -> node02

podB - > PodA

前面学习的Node亲和，让我们可以指定Pod将要调度的节点。这里要学习的Pod亲和，则关注的是Pod本身。

### Pod硬亲和

调度Pod时，期望和满足某些标签的Pod部署在一起，例如咱们部署前端服务的时候，希望和后端服务部署在一起。

```yaml
affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - topologyKey: kubernetes.io/hostname
        labelSelector:
          matchExpressions: 
          - key: app
            operator: In 
            values: 
            - backend
```
### Pod软亲和

```yaml
affinity:
  podAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - weight: 100
      podAffinityTerm:
        topologyKey: kubernetes.io/hostname
        labelSelector:
          matchExpressions:
          - key: security
            operator: In
            values:
            - S2
```

topologyKey:表示作用域的意思，满足该标签的节点在同一个作用域。
## Pod 反亲和调度

和Pod亲和刚好相反，不要和满足Pod的调度到一起。

podA -> node1

podB -> anynode(出podA所在的节点)

### Pod反亲和+硬亲和

```yaml
affinity:
  podAntiAffinity:
    requiredDuringSchedulingIgnoredDuringExecution:
    - topologyKey: kubernetes.io/hostname
      labelSelector:
        matchExpressions:
        - key: app
          operator: In
          values:
          - frontend
```
### Pod 反亲和+软亲和

```yaml
podAntiAffinity:
  preferredDuringSchedulingIgnoredDuringExecution:
  - weight: 100
    podAffinityTerm:
      labelSelector:
        matchExpressions:
        - key: security
          operator: In
          values:
          - S2
      topologyKey: kubernetes.io/hostname
```
