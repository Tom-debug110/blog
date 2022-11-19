# 跳表
## 引言
跳表是由链表实现的，但是在查询和插入的性能确毫不逊色于二叉树，所以，在大名鼎鼎的Redis 中，跳表就在其中被应用


具体的可以参考照片博文，介绍的比较详细啦。我当时很不理解的地方就是对于插入的时候，还有就是这个插入时选取的概率，为什么是 2.25 ，而不是其他的，可以参考下面的这个博客
[跳表](https://leetcode.cn/problems/design-skiplist/solutions/1696545/she-ji-tiao-biao-by-leetcode-solution-e8yh/)

[跳表](https://www.jianshu.com/p/9d8296562806)