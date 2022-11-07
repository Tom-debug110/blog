# Linux 下 Socket 编程

Socket 编程算是最基本的网络编程模型了，即使不是相关方面的专家，也应该对 `Socket` 编程相关的知识非常的熟悉

## 一、`Socket` 编程基本步骤
### 1.1 服务端
* 创建 `socket` （套接字）
* 绑定套接字返回的文件描述符和对应的服务端口
* 监听端口
* 接受客户端的连接
* 处理客户端发送的信息

### 1.2客户端
* 创建套接字
* 连接套节字返回的文件描述符和服务端口
* 发送信息

> 以上就是 `socket` 编程的大致步骤，相比之下，客户端的步骤显然要简单的多，下面会结合代码详细说明

## 二、服务端代码

### 2.1 创建套接字

套接字算是一种抽象概念
> 所谓套接字(Socket)，就是对网络中不同主机上的应用进程之间进行双向通信的端点的抽象。一个套接字就是网络上进程通信的一端，提供了应用层进程利用网络协议交换数据的机制。从所处的地位来讲，套接字上联应用进程，下联网络协议栈，是应用程序通过网络协议进行通信的接口，是应用程序与网络协议栈进行交互的接口    
>
>-- 王雷,TCP/IP网络编程基础教程,北京理工大学出版社,2017.02,第4页 

个人觉得上面的定义挺好的，他就是一个抽象概念，并不是指某一个具体的东西,同时套接字常用的话是有三种类型的，下面会看到
```c
  socket(int __domain, int __type, int __protocol)
```
通过 `socket()` 函数创建一个对应的套接字,同时返回和这个套接字对应的文件描述符,以下是函数内注释
> Create a new socket of type TYPE in domain DOMAIN, using protocol PROTOCOL. If PROTOCOL is zero, one is chosen automatically. Returns a file descriptor for the new socket, or -1 for errors.

* `__domain` 表示套接字支持的IP协议是 `IPv4` 还是 `IPv6` 
1. `AF_INET` 网络层使用`IPv4`协议簇
2. `AF_INET6` 网络层使用 `IPv6` 协议簇

    ... 
* `__type` 表示套接字类型
1. `SOCK_STREAM` 表示流式套接字，,传送的数据是有保证的，也被称为TCP套接字
2. `SOCK_DGRAM` 数据报套接字,不保障数据可以被对方接收到
3. `SOCK_RAW` 原始套接字，一般用于对较底层协议直接访问(IP ICMP等)，常用于网络协议分析，检验新的网络协议实现

* `__protocol` 指定网络协议，这个网络协议的指定适合前面的套接字类型紧密相关联的，如果前面使用了 `SOCK_STREAM` 就不能使用`IPPROTO_TCP` 组合啦
1. `IPPROTO_TCP` 采用 `TCP` 协议
2. `IPPROTO_UDP` 采用 `UDP` 协议
3. `IPPORT_FTP` 采用 `FTP` 协议

    ...

> 以上的只是冰山一角，如果查看 `/usr/include/x86_64-linux-gnu/bits/socket.h` 可以看到全部的参数类型, `socket()` 函数的申明位于 `/usr/include/x86_64-linux-gnu/sys/socket.h` 中

还有一个事情就是，上面的参数介绍完了，就是这个 `socket()` 函数的返回类型啦，`Linux` 遵循一切皆文件的概念，创建的 `socket` 如何进行数据交换呢？那就是使用创建这个 `socket` 时返回的 文件描述符来进行读写数据
### 2.2 绑定套节字文件描述符和套接字地址
> 标题很长，主要是希望表达清楚意思，这个时候上面创建套接字时返回的文件描述符也就派上用途了。

```cpp

# define __CONST_SOCKADDR_ARG	const struct sockaddr *
/* Give the socket FD the local address ADDR (which is LEN bytes long).  */
extern int bind (int __fd, __CONST_SOCKADDR_ARG __addr, socklen_t __len)
     __THROW;
```


* `struct sockaddr_in`
> Structure describing an Internet socket address. 
这个结构体位于 `/usr/include/netinet/in.h` 中，正如注释所描述的，该结构主要是用来描述套接字地址的

```cpp
/* Structure describing an Internet socket address.  */
struct sockaddr_in
  {
    __SOCKADDR_COMMON (sin_);
    in_port_t sin_port;			/* Port number.  */
    struct in_addr sin_addr;		/* Internet address.  */

    /* Pad to size of `struct sockaddr'.  */
    unsigned char sin_zero[sizeof (struct sockaddr)
			   - __SOCKADDR_COMMON_SIZE
			   - sizeof (in_port_t)
			   - sizeof (struct in_addr)];
  };
  ```

  ```cpp
  memset(&servaddr, 0, sizeof(servaddr));
  servaddr.sin_family = AF_INET;
  servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
  servaddr.sin_port = htons(6666);
```
在使用之前需要先把这块内存初始化以下的，也就是清空内存。
之后我们要对这个结构体进行一些配置，比如，设置协议簇，绑定IP地址，这里最初我也有疑问，作为服务端，IP地址还需要绑定吗？不就是我自己吗，端口号倒是需要绑定，说明是在哪个端口进行监听的。

其实绑定IP 主要是因为，如果一个设备上有多个网卡的话，那要在那个网卡上监听的？同时我们也会发现，我们一般把 `localhost` 解释为 `127.0.0.1`,还有一个 `0.0.0.0` .他们之间的关系比较复杂抽空可以单独写个文章。


配置好上面的 `struct sockaddr_in` 之后

```cpp
if (bind(listenfd, (struct sockaddr *)&servaddr, sizeof(servaddr)) == -1) {
    printf("bind socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
```
### 2.3 监听端口
> 如果上一步创建完成了，那么服务端就要开始监听了。

```cpp
if (listen(listenfd, 10) == -1) {
    printf("listen socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
```

### 2.4 接受连接
调用 `accept()` 处理客户端的连接

```cpp
 if ((connfd = accept(listenfd, (struct sockaddr *)NULL, NULL)) == -1) {
      printf("accept socket error: %s(errno: %d)", strerror(errno), errno);
      continue;
    }
```

### 2.5 处理数据

```cpp
    n = recv(connfd, buff, MAXLINE, 0);
    buff[n] = '\0';
    printf("recv msg from client: %s\n", buff);
    close(connfd);
```

### 2.6 全部代码
```cpp
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>
#define MAXLINE 4096
int main(int argc, char **argv) {
  int listenfd, connfd;
  struct sockaddr_in servaddr;
  char buff[4096];
  int n;

  if ((listenfd = socket(AF_INET, SOCK_STREAM, IPPORT_FTP)) == -1) {
    printf("create socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
  memset(&servaddr, 0, sizeof(servaddr));
  servaddr.sin_family = AF_INET;
  servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
  servaddr.sin_port = htons(6666);
  if (bind(listenfd, (struct sockaddr *)&servaddr, sizeof(servaddr)) == -1) {
    printf("bind socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
  if (listen(listenfd, 10) == -1) {
    printf("listen socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
  printf("======waiting for client's request======\n");
  while (1) {
    if ((connfd = accept(listenfd, (struct sockaddr *)NULL, NULL)) == -1) {
      printf("accept socket error: %s(errno: %d)", strerror(errno), errno);
      continue;
    }
    n = recv(connfd, buff, MAXLINE, 0);
    buff[n] = '\0';
    printf("recv msg from client: %s\n", buff);
    close(connfd);
  }
  close(listenfd);
}

```
## 三、客户端代码

相对来说，客户端的步骤和代码都要简单不少。但是也是有一些共性的地方

### 3.1 创建套接字


```cpp
if ((sockfd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
    printf("create socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
```

### 3.2 连接服务端
直接调用 `connect()` 函数即可，但是在调用之前要先创建服务端的套接字地址，和服务端那边同样的数据类型

```cpp
servaddr.sin_family = AF_INET;
  servaddr.sin_port = htons(6666);
  if (inet_pton(AF_INET, argv[1], &servaddr.sin_addr) <= 0) {
    printf("inet_pton error for %s\n", argv[1]);
    exit(0);
  }
```
> 上面是根据传进来的命令行参数进行转换，再赋值给 `servaddr.sin_addr` 成员

### 3.3 发送数据

```cpp
if (send(sockfd, sendline, strlen(sendline), 0) < 0) {
    printf("send msg error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
```


### 3.4 全部代码

```cpp


#include <arpa/inet.h>
#include <errno.h>
#include <netinet/in.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#define MAXLINE 4096
int main(int argc, char **argv) {
  int sockfd, n;
  char recvline[4096], sendline[4096];
  struct sockaddr_in servaddr;
  if (argc != 2) {
    printf("usage: ./client <ipaddress>\n");
    exit(0);
  }
  if ((sockfd = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
    printf("create socket error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
  memset(&servaddr, 0, sizeof(servaddr));
  servaddr.sin_family = AF_INET;
  servaddr.sin_port = htons(6666);
  if (inet_pton(AF_INET, argv[1], &servaddr.sin_addr) <= 0) {
    printf("inet_pton error for %s\n", argv[1]);
    exit(0);
  }
  if (connect(sockfd, (struct sockaddr *)&servaddr, sizeof(servaddr)) < 0) {
    printf("connect error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
  printf("send msg to server: \n");
  fgets(sendline, 4096, stdin);
  if (send(sockfd, sendline, strlen(sendline), 0) < 0) {
    printf("send msg error: %s(errno: %d)\n", strerror(errno), errno);
    exit(0);
  }
  close(sockfd);
  exit(0);
}
```