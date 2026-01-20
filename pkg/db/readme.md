---
last updated: 2026-01-16 am
---

# Database Design Research

Some database design and related research will be recorded here.

- [Query Analyze](#query_analyze)
- [RBAC Simple](#rbac_simple)
- [Comment System](#comment_system)
- [Redis Design](#redis_design)

---

## <a id="query_analyze">Query Analyze</a>

This section attempts a simple analysis of SELECT queries in SQL for relational databases, 
as this is the most common and frequently used **DML** operation. Ultimately, we simply want to extract the final dataset from a large dataset using the operations provided by SELECT,
representing a data processing procedure. In reality, perfectly using SELECT to obtain the
result isn't particularly complicated. In principle, it requires deconstructing the actual
problem in reverse. Let's try to understand this by starting with a practical problem:

*eg_1: Try to retrieve the name of the student with the highest score in each subject from stu
(id, name) and sc(id, subject, score).*

Here, our final result needs to be "name", calculated from the "highest score in each subject". Here's an example implementation:

`SELECT name FROM stu LEFT JOIN sc sc1 ON stu.id = sc1.id WHERE sc1.score = (SELECT max(sc2.
score) FROM sc sc2 WHERE sc2.subject = sc1.subject)`

One of the most interesting aspects of SELECT is its execution order. The `from table_expr
where ... group by ... having ... select ... order by ... limit ...` clause is intuitively
important, considering both the scope of influence of table aliases and logical operations. 
Aggregate functions cannot appear in the `where` clause, but they can be introduced after 
grouping.

## <a id="rbac_simple">RBAC Simple</a>

Implement a simple RBAC system called **NBAC**, containing only basic operations while possessing a certain
degree of scalability. We abstract the actual entities as: `user, namespace, permission, and 
resource`. Each entity has a unique identifier, which is for the purpose of allowing for easy 
expansion of operations, similar to the permission design in Unix.

**NBAC** implements a concept similar to a file system in an operating system. In fact, it 
involves similar operations, but only the minimal fields are implemented.

**NBAC** is a simplified version of the **RBAC** system, mainly used for permission 
management in specific system environments. It mainly uses the concept of namespace to 
restrict the operations of different users on resources. This requires a set of rules to 
restrict to a certain extent. For users, we can always create in a general and natural way. 
Users have natural permissions to create some resources, which is obviously only related to 
the resources themselves. If you need to build your own resources in the real world, it is 
obviously a lack of construction materials, not any necessary constraint permissions. This 
self-created resource, whether it is a namespace or a resource, has no restrictions on the 
creator. The creator can set outsider access permissions on the namespace and resource, which 
includes the default permissions of the resource itself, which can also be directly called 
the maximum accessible permissions. At the same time, the necessary access permissions can be 
directly specified for each user who may need to access the resource, which greatly 
facilitates refinement and is used to control the granularity of permissions. A resource is 
an abstract association collection, and it does not represent any concrete resource. Its 
extern_id is enough to associate it with any possible actual resource, but this only depends 
on the actual resource creation, which is not discussed in the permission system. Permissions 
should be a predefined limited set. Obviously using uint64 can basically cover any 
permissions(of course it must be well designed). The association between permissions is 
loosely coupled, which represents any possible independent capabilities. These detailed 
permissions must be reasonable. Assume that non-2^n permissions are allowed to be defined in 
permissions, This means that `code` will become a value that is not a power of 2 under 
uint64. But no matter what, we just regard this as a collection of some special permissions, 
similar to role combinations of some permissions, which means that namespace or resources can 
continue to use these. But one problem is that code is more likely to be defined as an index 
in the db, which can directly reduce the number of permissions in the table, thereby further 
restricting permissions.

## <a id="comment_system">Comment System</a>

**ref**
- [bilibili](https://www.bilibili.com/opus/737531797122842865)

## <a id="redis_design">Redis Design</a>
