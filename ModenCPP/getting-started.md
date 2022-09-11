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
###  `join()`
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

  //延迟两秒,暂时忽略即可
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
### `detach()`
`std::thread` 中的成员函数 `detach()` 是显示的分离一个  `std::thread` 对象和一个线程，让其在后台执行。并且 `std::thread` 对象在销毁时，也就是 `main()` 运行结束时，要清理资源，并不会调用 `std::terminate()` 但是这个线程的归属权和控制权都转移给 `C++ Runtime`.

以上代码如果使用  `detach()` 的话，最终的运行结果和不适用 `join()` 应该是一样的，但是程序会显示正常结束，而不是被 `std::terminate()` 中断而异常退出

<svg t="1662886891580" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="1376" width="16" height="16"><path d="M1001.661867 796.544c48.896 84.906667 7.68 157.013333-87.552 157.013333H110.781867c-97.834667 0-139.050667-69.504-90.112-157.013333l401.664-666.88c48.896-87.552 128.725333-87.552 177.664 0l401.664 666.88zM479.165867 296.533333v341.333334a32 32 0 1 0 64 0v-341.333334a32 32 0 1 0-64 0z m0 469.333334v42.666666a32 32 0 1 0 64 0v-42.666666a32 32 0 1 0-64 0z" fill="#FAAD14" p-id="1377"></path></svg>
一旦调用 `detach()` 函数，那么一个  `std::thread` 对象就和 线程分离啦，如果再次调用会引发异常。所以对于一个 `std::thread ` 对象来说，只有在确保其是 `joinable()` 时，才可以调用 `detach()`.同理对于 `join()` 来说也是一样的

```cpp
if (t.joinable()) {
    t.detach();
  }

  if (t.joinable()) {
    t.join();
  }
  ```
  ## 向线程函数传递参数

> 向一个线程函数或者说在一个线程里面被调用的可执行单元传递参数，存在一个问题就是是值的拷贝还是引用,其实都不是的，严格来说是右值引用







