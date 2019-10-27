# fuzz_debug_platform


## 接口及负责人

### graph


Get:   `/graph`


用于告知前端fuzz测试使用的yy文件的图结构

负责人: 杜沁园  @DQinYuan

[view/graph.go](view/graph.go)


示例输出：

以下面的yy文件为例：

```
query:
    select | create

select:
    SELECT * FROM t
	
create:
    CREATE TABLE m AS select
```

他所产生的json是：


```
[
  {
    number: 0,
    head: "query",
    alter:  [{
       content: "select"，
       fanout: [1]
     },
    {
      content: "create",
      fanout: [2],
    }]
    },
   {
    number: 1
    head: "select",
    alter: [{
         content: "SELECT * FROM t",
         fanout:[],
       }
     ] 
   },
   {
     number: 2,
     head: "create",
     alter:[{
         head: "CREATE TABLE m AS select",
         fanout: [1]
       }
     ]
   }
]
```

### heat

Get:  `/heat`

负责人：杜沁园 @DQinYuan

[view/heat.go](view/heat.go)

用于告知前端yy文件产生的图的每个节点的热力以及一些附带信息

示例输出：

```

[
  {
     number: 0,   // 对应节点的编号
     alter: 0,      // 对应分支的编号(即之前传的alter数组的下标)
     heat: 30,    // 热力值
     sql:  "xxxx",   //示例sql
   },
   {
     number: 0,
     alter: 1,
     heat: 13,
     sql: "yyy"
   },
   {
     number: 1,
     alter: 0,
     heat: 3,
     sql: "zzz"
   },
  {
    number:2,
    alter: 0,
    heat: 23,
    sql: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  },
]
```

### codepos

Get: `/codepos`

负责人： 满俊朋  @mahjonp

[view/codepos.go](view/codepos.go)

示例输出：

```
```

## 模块及负责人

 - sqldebug: 杜沁园 @DQinYuan
 - sqlfuzz: 满俊朋   @mahjonp