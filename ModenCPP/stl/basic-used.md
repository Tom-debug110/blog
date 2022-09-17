# STL 基本用法回顾

## `vector<>`
> 动态变长数组,可以自行扩充内存。
### 常用成员函数

```cpp
namespace fmt {
void println(std::vector<int> &v) {
  for (int i : v) {
	std::cout << i << " ";
  }
  std::cout << std::endl;
}
}
```

#### 1. `assign()` 
> 给调用对象赋值使用，但是会擦除调用对象原来的值，标准库提供了三种重载
* 赋指定个数的值

`constexpr void assign( size_type count, const T& value );`
> 比如需要一个初始化全为某一个值的 `std::vector<>`

* 把迭代器指定范围的值赋给调用对象

`constexpr void assign( InputIt first, InputIt last );`
> 根据现有的同类型容器创建

* 由初始化列表进行赋值
`constexpr void assign( std::initializer_list<T> ilist );`
> 这样的话就和直接根据初始化列表创建容器没有什么区别啦

**代码示例**
```cpp
#include<vector>
#include<iostream>
namespace fmt {
void println(std::vector<int> &v) {
  for (int i : v) {
	std::cout << i << " ";
  }
  std::cout << std::endl;
}
}
int main() {
  std::vector<int> v{1, 2, 3, 4, 5, 6, 45};
  std::vector<int> v1{10, 8, 7, 9};

  // 根据迭代器创建
  v1.assign(v.begin(), v.begin() + 3);
  fmt::println(v1);

  // 根据初始化列表
  v1.assign({7, 8, 9, 7, 45, 56, 12, 35});
  fmt::println(v);

  // 创建指定数量的类型
  v1.assign(5, 45);
  fmt::println(v1);

  return 0;
}
```
#### 2. `get_allocator()` 
> 这个其实并不常用，就是返回 `vector<>` 底层的内存分配器

#### 3. `at()` 
> 获取指定下标的元素，相较与 `[]` 来说，多了下标溢出检查

#### 4. `operator[]`
> 访问指定下标元素

#### 5. `front()` 
> 返回第一个元素的引用

---

一个空的 `std::vector` 在调用以上几个成员函数时，均会出现异常
```cpp

  std::vector<int> v1;
  std::cout << v1.front() << std::endl;
  std::cout << v1.at(0) << std::endl;
  std::cout << v1[0] << std::endl;
```
#### 6. `back()`
> 返回最后一个元素的引用

#### 7. `data()`
> 获取底层存储数据的数组的指针，主要应该是为了兼容 `C` 而特意留下的
```cpp
#include<vector>
#include<iostream>
namespace fmt {
void println(std::vector<int> &v) {
  for (int i : v) {
	std::cout << i << " ";
  }
  std::cout << std::endl;
}
void println(const int *p, int size) {
  for (int i = 0; i < size; i++) {
	std::cout << p[i] << " ";
  }
  std::cout << "\n";
}
}
int main() {
  std::vector<int> v{1, 2, 3, 4, 5, 6, 45};

  fmt::println(v.data(), v.size());

  return 0;
}
```
#### 8. `begin()`和 `cbegin()`
> 返回的是容器开端的迭代器， `c` 表示是一个常量
#### 9. `end()` 和 `cend()`
> 返回容器尾部的迭代器，同理，`c` 表示常量

#### 10. `empty()`
> 容器为空返回 `true` 否则返回 `false`
#### 11. `size()`
> 返回一个容器实际的元素数量

#### 12. `max_size()` 
> 返回可能的最大的容器大小，大小也受容器所含元素类型的影响

#### 13. `reserve()` 扩充容器的 `capacity`
> 只有在新的参数大于等于 当前的 `capacity` 时，这个函数才会进行扩容，实际上也就是内存的分配，否则，该函数什么也不做

#### 14. `clear()` 
> 清除容器内的内容，但是并不会回收内存
---
```cpp
int main() {
  std::vector<int> v{1, 2, 3, 4, 5, 6, 45};

  std::vector<int> v1;

  v1.assign(v.cbegin(), v.cend());

  std::cout << "capacity is :" << v1.capacity() << std::endl;
  std::cout << "size is :" << v1.size() << std::endl;
  std::cout << "max_size is :" << v1.max_size() << std::endl;
  std::cout << "vector if or not empty:" << std::boolalpha << v1.empty() << std::endl;

  v1.reserve(v1.capacity() + 10);
  std::cout << "after reserve, capacity is :" << v1.capacity() << std::endl;

  v1.clear();
  std::cout << "after clear,the size is:" << v1.size() << std::endl;
  fmt::println(v1);
  
  return 0;
}
```
#### 15. `insert()`
> 在指定的位置**之前**插入元素，有多种重载
* 在指定位置**之前**插入一个元素,其中又包含两个重载 ，对应左值和右值
`constexpr iterator insert( const_iterator pos, const T& value );`
* 在指定位置之前插入指定数量的元素
`constexpr iterator insert( const_iterator pos, size_type count, const T& value );`

* 在指定位置之前插入由迭代器指定的序列
`constexpr iterator insert( const_iterator pos, InputIt first, InputIt last );`
* 在指定位置之前插入由初始化列表指定的序列
`constexpr iterator insert( const_iterator pos,std::initializer_list<T>ilist )`

> 对于一个迭代器序列来说，The behavior is undefined if first and last are iterators into *this.
所以，对于这种情况要慎重

```cpp
#include<vector>
#include<iostream>
namespace fmt {
void println(std::vector<int> &v) {
  for (int i : v) {
	std::cout << i << " ";
  }
  std::cout << std::endl;
}
void println(const int *p, int size) {
  for (int i = 0; i < size; i++) {
	std::cout << p[i] << " ";
  }
  std::cout << "\n";
}
}
int main() {
  std::vector<int> v{1, 2, 3, 4, 5, 6, 45};
  std::vector<int> v1{89, 100, 784, 230, 158};

  v1.insert(v1.begin(), 999);
  fmt::println(v1);

  v1.insert(v1.begin(), 10, 888);
  fmt::println(v1);

  v1.insert(v1.begin(), v.begin(), v.end());
  fmt::println(v1);

  v1.insert(v1.begin(), {777, 777, 777, 666, 66, 333});
  fmt::println(v1);

  return 0;
}
```
> 对于插入一个迭代器区间，即使这个迭代器是自己的,我自己在尝试的时候也没有异常发生，包括区间重叠，但是还是按照标准库的建议，尽量不要这样做

#### 16. `emplace()`
> 在指定的位置之前构造元素
这是 `C++11` 新增的成员函数，又被称为是原位构造
> Inserts a new element into the container directly before pos.
The element is constructed through std::allocator_traits::construct, which typically uses placement-new to construct the element in-place at a location provided by the container. However, if the required location has been occupied by an existing element, the inserted element is constructed at another location at first, and then move assigned into the required location.
The arguments args... are forwarded to the constructor as std::forward<Args>(args).... args... may directly or indirectly refer to a value in the container.
If the new size() is greater than capacity(), all iterators and references are invalidated. Otherwise, only the iterators and references before the insertion point remain valid. The past-the-end iterator is also invalidated.

可以阅读一下原文

#### 17. `erase()`
清除指定位置和指定区间的元素，并不是单单的使用零值代替

#### 18. `push_back()`
向容器末尾添加元素
> 1) The new element is initialized as a copy of value.
> 2) value is moved into the new element.
#### 19. `emplace_back()`
同样也是向容器末尾添加元素，但是底层机制有很大差异
> Appends a new element to the end of the container. The element is constructed through std::allocator_traits::construct, which typically uses placement-new to construct the element in-place at the location provided by the container
#### 20. `pop_back()`
弹出最后一个元素
#### 21. `resize()`
和 `reserve()` 还是存在很大差别的，如果第一个参数值大于 `size()` 那么就在容器末尾添加默认的零值或者使用调用者提供的值
如果第一个参数小于当前的 `size()`，那么容器就缩减，但是容量不变

#### 22. `swap()`
交换两个容器的内容，没记错的话，是指针直接交换，减少拷贝和复制的开销，但实际上我们经常用来释放全部或者部分容器的内存


---
## `std::string`
[std::string](https://en.cppreference.com/w/cpp/string/basic_string/assign)
#### 1. `assign()`
> 和 `std::vector<>` 异曲同工，只不过了一个被赋值的参数，也是存在多种重载

* 直接从一个已经存在的 `std::string` 中进行赋值操作
`constexpr basic_string& assign( const basic_string& str );`
* 指定被赋值字符串的开始位置个从开始位置需要赋值多少个字符
`basic_string& assign( const basic_string& str,size_type pos, size_type count = npos);`

* 支持传入 `C` 风格字符串
`constexpr basic_string& assign( const CharT* s );`
* 使用迭代器
`template< class InputIt >constexpr basic_string& assign( InputIt first, InputIt last );`
* 支持使用初始化列表
`constexpr basic_string& assign( std::initializer_list<CharT> ilist );`
#### 2. `get_allocator()`
#### 3. `at()`
#### 4. `operator[]`
#### 5. `front()`
#### 6. `back()`
#### 7. `data()`
#### 8. `c_str()`
#### 9. `begin()` 和 `cbegin()`
#### 10. `end()`和 `cend()`
#### 11 `empty()`
#### 12. `size()` 和 `length()`
#### 13. `max_size()`
#### 14. `reserve()`
#### 15. `capacity()`
#### 16. `clear()`
#### 17. `insert()`
> 主要是插入相关的操作,对于字符串处理来说还是比较经常使用的操作

* 在指定位置(索引)上插入若干个给定字符
`constexpr basic_string& insert( size_type index, size_type count, CharT ch );`
* 在指定位置插入一个`\0` 结尾的字符
`constexpr basic_string& insert( size_type index, const CharT* s );`
* 在指定位置插入给定字符串的子串
`constexpr basic_string& insert( size_type index,const CharT* s, size_type count );`
也就是插入 `s[0..count)` 字串
* 在指定位置插入一个字串
`constexpr basic_string& insert( size_type index, const basic_string& str );`

* 在指定位置插入一个字符，该位置使用迭代器指定
`constexpr iterator insert( const_iterator pos, CharT ch );`

* 在指定位置插入指定数量的字符
`constexpr iterator insert( const_iterator pos, size_type count, CharT ch );`

* 插入迭代器范围指定的序列
`template< class InputIt >constexpr iterator insert( const_iterator pos, InputIt first, InputIt last );`

* 插入初始化列表指定的序列

`constexpr iterator insert( const_iterator pos,std::initializer_list<CharT> ilist );`

> 不管是使用迭代器还是初始化列表，其插入位置都是使用迭代器来指定的

> 关于 `std::string_view<>` 暂时不讨论
 

