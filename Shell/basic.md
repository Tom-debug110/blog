# Linux Shell 经典脚本

## 猜数游戏

```bash
#!/bin/bash

# 脚本生成一个 100 以内的随机数,提示用户猜数字,根据用户的输入,提示用户猜对了,
# 猜小了或猜大了,直至用户猜对脚本结束。

# RANDOM 为系统自带的系统变量,值为 0‐32767的随机数
# 使用取余算法将随机数变为 1‐100 的随机数
num=$[RANDOM%100+1]
echo "$num"

# 使用 read 提示用户猜数字
# 使用 if 判断用户猜数字的大小关系:‐eq(等于),‐ne(不等于),‐gt(大于),‐ge(大于等于),
# ‐lt(小于),‐le(小于等于)
while :
do 
 read -p "计算机生成了一个 1‐100 的随机数,你猜: " cai  
    if [ $cai -eq $num ]   
    then     
        echo "恭喜,猜对了"     
        exit  
     elif [ $cai -gt $num ]  
     then       
            echo "Oops,猜大了"    
       else      
            echo "Oops,猜小了" 
  fi
done
```

## 查看有多少远程IP在连接本机

```bash
#!/bin/bash
# 查看有多少远程的 IP 在连接本机(不管是通过 ssh 还是 web 还是 ftp 都统计) 

# 使用 netstat ‐atn 可以查看本机所有连接的状态,‐a 查看所有,
# -t仅显示 tcp 连接的信息,‐n 数字格式显示
# Local Address(第四列是本机的 IP 和端口信息)
# Foreign Address(第五列是远程主机的 IP 和端口信息)
# 使用 awk 命令仅显示第 5 列数据,再显示第 1 列 IP 地址的信息
# sort 可以按数字大小排序,最后使用 uniq 将多余重复的删除,并统计重复的次数
netstat -atn  |  awk  '{print $5}'  | awk  '{print $1}' | sort -nr  |  uniq -c
```

## 打印 `Hello World` 

```bash
#!/bin/bash


# 注意函数名之后要空一个空格给打括号
function example {
    echo "Hello World"
}

example
```


## 打印某一个进程的pid

```bash
#!/bin/sh`

v1="Hello"
v2="world"
v3=${v1}${v2}
echo $v3


# 查看当前进程 显示vscode 进程，打印第二个参数
pidlist=`ps -ef|grep code|grep -v "grep"|awk '{print $2}'`
echo $pidlist
echo "tomcat Id list :$pidlist"  #显示pid

```

## 石头剪刀布


```bash
# !/bin/bash

# 

game=(石头 剪刀 布)

num=$[RANDOM%3]

computer=${game[$num]}

echo "请出拳"
echo "1. 石头" 
echo "2. 剪刀"
echo "3. 布"

read -p "请选择" person
case $person in
1)
  if [ $num -eq 0 ]
  then 
    echo "平局"
    elif [ $num -eq 1 ]
    then
      echo "你赢"
    else 
      echo "计算机赢"
fi;;
2)
if [ $num -eq 0 ]
then
    echo "no"
elif [ $num -eq 1 ]  
then
    echo "pingju"
else
    echo "Yes"
fi;;
3)
if [ $num -eq 0]
then 
    echo "Yes"
elif [ $num -eq 1 ]
then 
    echo "no"
else 
    echo "pingju"
fi;;
esac

```

## 九九乘法表

```bash

#!/bin/bash
for i in `seq 9`
do 
    for j in `seq $i`
    do 
    echo  "$j*$i= $[i*j]"
    done
    echo
done

```

## 判断当前用户是否是 root

```bash

#!/bin/bash
if [ $USER == "root" ]
then 
    echo "hello root"
else
    echo "no root"

fi

```
