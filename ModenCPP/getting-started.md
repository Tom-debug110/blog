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

### 方式二
> 使用 `lambda` 表达式

`lambda` 表达式是 `C++11` 新的特性，具体的这里不过多介绍。但是在这里，一个 `lambda` 表达式就是一个可执行单元。
```cpp
#include <iostream>
#include <thread>

int main() {
  std::cout << "First concurrency program\n";
  std::thread t(
      [](int a, int b) {
        std::cout << "Hello lambda\n";
        std::cout << "a + b is:" << a + b << std::endl;
      },
      3, 11);

  // 如果运行报错，添加以下语句
  t.join();

  return 0;
}
```
以上就是使用  `lambda` 开启一个线程的方法，参数也是放到后面传递给实际执行的可调用单元

同时我们还可以在 `lambda` 里面再次调用一个函数或者其他的执行单元

```cpp
#include <cmath>
#include <iostream>
#include <thread>

int main() {
  std::cout << "First concurrency program\n";
  std::thread t(
      [](int a, int b) {
        std::cout << "Hello lambda\n";
        std::cout << "a ^ b is:" << pow(a, b) << std::endl;
      },
      3, 11);

  // 如果运行报错，添加以下语句
  t.join();

  return 0;
}
````



## 管控线程
> `std::thread` 只是一个类而已，他的作用也仅仅是发起一个线程，让传递进来的可执行单元在新的线程上运行


上面代码中提到了一个 `std::thread` 的成员函数  `join()` 他的意思可以理解为主线程等待子线程结束，在我们的代码里面可以理解为是 `main()` 所运行的线程等待可执行单元所在的线程，一直到这个可执行单元执行完毕。

如果不等待的话并且在主线程结束子线程还未完成时， `std::thread()` 就会调用 `std::terminate()` 

可以在一下代码中尝试执行 `join()` 操作和不执行的结果

```cpp
#include <chrono>
#include <cmath>
#include <iostream>
#include <thread>
void concurencyFunc() {
  std::cout << "concurrencyFunc:\n start\n";

  std::this_thread::sleep_for(std::chrono::milliseconds(2000));

  std::cout << "end\n";
}
int main() {
  std::cout << "First concurrency program\n";

  std::thread t(concurencyFunc);

  // t.join();

  std::cout << "main thread end\n";

  return 0;
}
```
