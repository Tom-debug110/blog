# 并发编程初探
本节内容主要是C++ 并发编程初步探讨
## 开启一个线程

首先来看并发编程的第一个例子
### 方式一
```cpp
#include <iostream>
#include <thread>

void concurrencyFunc(int a, int b) {
  std::cout << "a is " << std::endl;
  std::cout << "b is " << std::endl;
  std::cout << "a + b is :" << a + b << std::endl;
}
int main() {
  std::cout << "First concurrency program\n";
  std::thread t(concurrencyFunc, 9, 6);

  // 如果运行报错，添加以下语句
  // 详见下一节 [管控线程]
  //t.join();

  return 0;
}
```
因为 `C++11` 自带了 `std::thread` 所以完全可以舍弃 `C` 风格的那一套 (`pthread`)啦。
在以上代码中，我们定义了一个 函数(可执行单元的形式之一)，在 ` main()` 里面创建了一个 `std::thread` 对象 (`object`) 。同时传入一个 **可调用对象**(`callable object`),除此之外，在后面传递这个可调用对象的参数。那么系统就会发起一个线程去单独执行我们的  `concurrencyFunc()`；  

**这样和直接在 `main()` 函数里面调用**有什么区别？

区别就是 `concurrencyFunc()` 是单独一个线程去执行的，直接调用的话就是在 `main()` 执行线程里面和其他单元串行执行。

## 管控线程
