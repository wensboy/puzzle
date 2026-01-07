---
title: "pkg section"
date: "2026/01/05 pm"
author: "wendisx"
---

# pkg section

- [database](#db)
  - [sqlite](#sqlite)
  - [mysql](#mysql)
  - [pgsql](#pgsql)
- [log](#log)
  - [non-structed](#nsl)
  - [structed](#sl)
- [router](#router)

---

some simple design ideas.

```go
/*
    索性直接用golang code block来写思路. 实际上, puzzle我有一段时间没有编写了. 最早是为了自己能够在闲暇之余能够找点除了video之外的乐趣, 比如写点通用的代码, 做点通用的工程. 作为一个主要做backend的人, 当然需要写好一个差不多的server框架才行了. (这里并不是直接写web framework, 那个没时间研究同时没技术完成.) 这次不写那么多废话了, 毕竟是写给自己使用的, 怎么方便怎么来. 也是考完试了, 写点golang放松一下, 总不能每天写算法. (虽然算法水平一般.)
*/
// 防止因为过多的随意取名导致function调用时来回切换, 简单介绍一些原则.
func InitX(args...) // 初始化默认的对象
func SetupX(args...) // 配置默认对象: 采用function配置形式
func ResetX() // 重置对象属性
func FreeX() // 释放对象
```

## <a id="db">database</a>

database 如何设计api. 其实这就是调用习惯的问题: init() -> setup() -> run() -> free(). database在后端项目中的唯一作用就是维护一个local database buffer pool用于直接和database进行交互. 为了软件工程上的开闭原则, 最小化的对象操作在database层面上无非就是crud四个. 因此设计时只需要尝试针对一个类型的database设计出对指定对象的最小可抽象操作即可.

```go
/*
    记录来自database package中的接口声明和细节
    1. golang并不支持method泛型, 这直接导致在多db下没有一个较好的封装实现, 只能尝试借助通用方法.
    2. 原则上将place和name参数分离, 防止一个function的功能过大.
    3. 请求实现分为one,more,page查询
    4. model设计区分固定字段和属性字段
*/
package database
// loc(package): pkg/db
// loc(test): pkg/db/database_test.go [passed]
```

# <a id="log">log</a>

log 设计主要包含几个要点: 显示提示, 内容输出, 结构化和非结构化. 在 container 环境中, 显然最核心的就是非结构化日志, 主要用于人可读的形式呈现, 而同时又需要保留结构化接口, 用于日志分析. 因此直接封装 slog 可以合理地处理这些. 

```go
/*
    记录来自 colorable log -- clog 的接口声明和细节
*/
package clog
// loc(package): pkg/clog
// loc(test): pkg/clog
```
