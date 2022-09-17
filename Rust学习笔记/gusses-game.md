# 猜数游戏

```rust
use std::cmp::Ordering;
use std::io;
use rand::Rng;

fn main() {
    println!("Guess the number!");

    let secret_num = rand::thread_rng().gen_range(10..101);
    println!("the secret number is: {}", secret_num);
    println!("Please input your number:");

    loop {
        let mut guess = String::new();

        io::stdin().read_line(&mut guess).expect("Failed to read line");

        let guess: i32 = match guess.trim().parse() {
            Ok(num) => num,
            Err(_) => continue,
        };


        match guess.cmp(&secret_num) {
            Ordering::Less => {
                println!("Too Small!");
            }
            Ordering::Equal => {
                println!("You Win!");
                break;
            }
            Ordering::Greater => {
                println!("Too Big");
            }
        }
    }
}
```

注意其中生成随机数的哪一行，不少书籍上使用的是 `0.3.x` 版本的写法，其实在目前的 `0.8.x` 版本中已经更换了，对应的参数也不太一样