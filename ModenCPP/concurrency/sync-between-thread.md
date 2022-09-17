# 2.4 并发操作的同步

> 不同线程之间有时候需要满足一定的执行顺序，那么就要采取一定的措施。

## 等待事件或等待其他条件
**引入**
假如我们乘坐一趟长途火车，我们需要在目标站点下车，那么为了可以确保准时下车，要么彻夜不眠，一直等着车到达最终的目的地，但是这样会疲惫不堪，另外就是定一个闹钟，每过一段时间就查看一下，这个看起来还不错，但是假如闹钟坏了呢... 其实最好的办法就是乘务员或者车厢广播在即将到站时进行通知。
---

换到我们的并发编程上，假如有甲乙两个线程，甲线程需要在乙线程完成任务后才可以执行，那么有一下几种方式

* 方式一
> 在共享数据内部维护一个受互斥保护的标志，线程乙完成后就设置为标志成立

但是以上方式存在问题，首先就是线程甲需要不断查验标志，其次就是甲在查验标志的时候要锁住互斥，那么如果线程乙恰好同时完成任务，也欲设置标志成立，则线程乙无法锁住互斥。
* 方式二
> 在每次查验之间短暂休眠，使用  `std::this_thread::sleep_for()` 函数实现

大概的代码如下所示

```cpp
bool flag = false;
std::mutex mu;

void wait_for_flag() {
  std::unique_lock<std::mutex> uq(mu);
  while (!flag) {
    uq.lock();
    std::this_thread::sleep_for(std::chrono::seconds(2));
    uq.unlock();
  }
}
```

以上方式看似解决了我们的问题，但是也随之引入了新的问题，那就是，休眠时间长短带来的问题，如果休眠时间过长,就会过度休眠，导致线程乙已经完成了任务，但是甲还在休眠。休眠时间过短就会让休眠没有效果

* 方式三
> 使用 `C++` 标准库的工具来等待时间发生

对于以上线程模型。我们完全可以使用 `std::condition_variable` 

### 凭借条件变量等待条件发生

```cpp
#include <atomic>
#include <chrono>
#include <iostream>
#include <mutex>
#include <queue>
#include <thread>

bool flag = false;
std::mutex mu;
std::condition_variable condition;
std ::queue<int> data_queue;
void write() {
  std::this_thread::sleep_for(std::chrono::seconds(1));
  {
    std::lock_guard<std::mutex> lk(mu);
    data_queue.push(7);
  }

  condition.notify_one();
}

void read() {
  std::unique_lock<std::mutex> lk(mu);

  condition.wait(lk, []() { return !data_queue.empty(); });

  int value = data_queue.front();
  data_queue.pop();

  lk.unlock();

  std::cout << value << std::endl;
}

int main() {
  std::thread t1(write);
  std::thread t2(read);

  if (t1.joinable()) {
    t1.join();
  }

  if (t2.joinable()) {
    t2.join();
  }

  return 0;
}
```

对于以上代码，详细解释一下关于条件变量的用法，首先在 `write()` 函数内，当队列中已经推入数据的时候，就会调用 `std::condition_varaable` 的成员函数 `notify_one()` 通知正在等待的线程甲。

对于线程甲来说，在 `std::condition_variable` 的成员函数 `wait()` 内部调用  `lambda` 表达式，如果表达式结果为真，那么就 `wait()` 就返回，执行下面的语句，如果表达式结果为  `false` 就 解锁互斥，继续进入等待状态

为什么使用 `std::unique_lock` ?
> 上面已经说了，如果 `wait()` 里面的 `lambda` 结果为 `false`， 那么这个时候就要解锁互斥啦,注意， `std::unique_lock` 的构造函数如果不传入第二个参数的话，还是会默认锁住互斥的，`std::lock_guard<>` 不提供对互斥的灵活操作，仅仅是构造的时候锁住互斥，析构的时候解锁互斥

## 使用 `std::future` 等待一次性事件

**引入**

假如我们在机场候机，等待登机，那么登机事件仅仅会发生一次。


### 从后台任务返回值

> 在 `C++` 内部，使用 `std::future` 来模拟这样的一次性事件的发生。在  `C++` 内部最基本的就是置于后台的计算任务完成。得出结果

```cpp
#include <chrono>
#include <future>
#include <iostream>
#include <mutex>
#include <thread>

int add(int a, int b) {
  std::this_thread::sleep_for(std::chrono::seconds(2));
  return a + b;
}
int main() {
  std::future<int> f = std::async(add, 8, 9);

  std::this_thread::sleep_for(std::chrono::seconds(2));

  std::cout << f.get() << std::endl;

  return 0;
}
```
上面出现了一个新的 `std::async()` 这是一个函数模板，也位于 `<future>` 头文件中，可以表示为按照异步的方式启动任务。同时从 `std::async()` 处获得 `std::future` 对象运行的函数一旦完成，其结果就有该 `std::future` 对象所有；如果要使用这个值，那就再这个对象上调用 `get()` 即可。

`std::async()` 的使用方法和 `std::thread` 差不多。包括传参方式
但是也有不同的地方，就是 `std::async` 可以指定一个参数指定任务的运行方式，参数类型是  `std::launch` ，`std::launch::deferred` 表示在当前线程上延后调用任务函数，等到在 `std::future` 对象上调用了 `wait()` 或者 `get()` 函数后任务才会运行；`std::launch::async` 指定必须另外开启专属的线程，在其上运行任务；该函数的值还可以是  `std::launch::async|std::launch::deferred` ，表示由 `std::async()` 自行选择运行方式。

### 关联 `std::future` 实例和任务

`std::packaged_task<>` 模板类连结了 `future` 对象和函数(可调用对象),一个很有用的场景就是在线程之间传递任务或者传递函数(可调用对象)

先说一下 `std::packaged_task<>` 的基本用法吧

其模板参数是函数签名，和 `std::function<>` 很像,比如 `std::packaged_task<void()> pt`; 表示没有返回值也没有参数。

`std::packaged_task<int(std::string &s,int num)> pt`; 表示返回类型为 `int` 参数为 `std::string` 和 `int` .

同时，此类模板还具有成员函数 `get_future()` 返回一个 `std::future<>` 实例，其类型取决与  `std::packaged_task<>`

```cpp
#include <deque>
#include <future>
#include <iostream>
#include <mutex>

std::mutex mu;
std::deque<std::packaged_task<void()>> tasks;
void render() {
  while (!tasks.empty()) {
    std::packaged_task<void()> currentTask;
    {
      std::lock_guard<std::mutex> lk(mu);

      currentTask = std::move(tasks.front());
      tasks.pop_front();
    }

    currentTask(); // 调用 传递进来的任务
  }
}
template <typename Func>

std::future<void> postMessage(Func f) {
  std::packaged_task<void()> task(f);

  std::future<void> res = task.get_future();
  std::lock_guard<std::mutex> lk(mu);

  tasks.push_back(task);

  return res;
}
```
### 创建 `std::promise` 

```cpp
#include <future>
#include <iostream>
#include <thread>

void fn(std::future<std::string> *future) {
  // future会一直阻塞，直到有值到来
  std::cout << "wait the result\n";
  std::cout << future->get() << std::endl;
}

int main() {
  // promise 相当于生产者
  std::promise<std::string> promise;
  // future 相当于消费者, 右值构造
  std::future<std::string> future1 = promise.get_future();
  // 另一线程中通过future来读取promise的值
  std::thread t(fn, &future1);
  // 让read等一会儿:)
  std::this_thread::sleep_for(std::chrono::seconds(3));
  //
  promise.set_value("hello future");
  // 等待线程执行完成
  t.join();

  return 0;
}
```

## 等待时间期限
### 时长类
`std::chrono::duration<>` 接受两个模板参数，第一个是参数表示使用何种类型计数计时单元是数量，可以采用 `int` `long` `double` 等类型。最后一个参数表示一个计时单元代表多少秒，因为这些计时单位都是按照秒来计算的，比如一个小时的时长可以表示为
`std::chrono::duration<int,std::ratio<3600,1>>` 其中 `int` 类型也可以使用 `long long` 或者 `double` 类型来表示。不仅仅可以用来表示小时，还可以表示微妙和毫秒这些
如
```cpp
  std::chrono::duration<int, std::ratio<1000, 1>> milli;
  std::chrono::duration<int, std::ratio<1000000, 1>> nona;
```

实际上，标准库里面已经给出了一组 `tpedef` 预设的定义:

|定义|含义|
|:-:|:-:|
|`nanoseconds`| 纳秒|
|`mircroseconds`| 微秒| 
|`milliseconds`| 毫秒 |
|`seconds`| 秒 |
|`minutes`| 分钟 |
|`hours`| 小时|


> 以上的定义都包含在 `std::chrono`  `namespace` 中

同时由于各种  `duration` 不同，标准库还提供了一个 `duration_cast<>`



```cpp
std::chrono::duration<int> seconds(3600); //3600 秒，第二个参数默认是 1
  std::chrono::duration<int, std::ratio<3600, 1>>// 一小时是 3600 秒
  hours = std::chrono::duration_cast<std::chrono::duration<int, std::ratio<3600, 1>>>(seconds);

  std::cout << hours << std::endl;

```
> 这样转换是OK 的

最后输出 `1h` 

下面这样也是 OK 的

```cpp
std::chrono::duration<int> seconds(3600); //3600 秒，第二个参数默认是 1
  std::chrono::duration<int, std::ratio<3600, 1>>// 一小时是 3600 秒
  hours = std::chrono::duration_cast<std::chrono::hours>(seconds);
```

> 邪了门了，当时我自己在实验的时候，是不成功的，`3600s` 会变成 `3600h`,本来以为是因为  `gcc12` 和 `MSVC` 不同造成的差异，都又重新试了试之后，发现又好了，这就很奇怪啦


最后补充两点就是:

1. 计时单元都有一个  `count()` 成员函数，返回计时单元的数量，其实也就是，比如 `3h` 就是  `3` ，`45minutes` 就是 `45` 
2. 在 `C++14 以上版本`，可以使用 `std::chrono_literals`,最终导致就是可以直接使用 `auto m=45min` 这样的表达式来表示不同时长。

---


### 时间点类

`std::chrono::time_point<>`

```cpp
#include <iostream>
#include <ratio>
#include <chrono>
#include <thread>
using namespace std::chrono_literals;
int main() {
  std::chrono::time_point<std::chrono::steady_clock> start = std::chrono::steady_clock::now();

  std::this_thread::sleep_for(std::chrono::seconds(2));

  std::chrono::time_point<std::chrono::steady_clock> stop = std::chrono::steady_clock::now();

  std::cout << std::chrono::duration_cast<std::chrono::seconds>(stop - start) << std::endl;

  // 以毫秒形式输出
  std::cout << std::chrono::duration_cast<std::chrono::milliseconds>(stop - start) << std::endl;

  return 0;
}
```
以上代码可以算是一个简单的统计成员运行时间的demo，可以根据要求的时间精度，改变 `duration_cast<>` 里面的参数

## 运用同步操作简化代码
