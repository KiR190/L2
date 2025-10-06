/*

Программа выводит числа от 1 до 8 в произвольном порядке,
потому что значения приходят из двух горутин со случайной задержкой,
а select выбирает один из каналов

С помощью select можно конкурентно получать значения сразу из нескольких каналов.
В данном случае select ждёт данные либо из a, либо из b
Если данные готовы в обоих каналах — select выбирает случайный case
Если один из каналов закрыт, заменяем его переменную на nil, и он больше никогда не будет выбран.
Канал nil в Go никогда не выдаст значения, и выборка из него в select невозможна, так что фактически мы отключаем его обработку.
Таким образом, цикл продолжается, пока хотя бы один канал активен

Использование select позволяет объединить два канала в один и корректно обработать закрытие каждого из них

*/
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func asChan(vs ...int) <-chan int {
    c := make(chan int)
    go func() {
        for _, v := range vs {
            c <- v
            time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
        }
        close(c)
    }()
    return c
}

func merge(a, b <-chan int) <-chan int {
    c := make(chan int)
    go func() {
        for {
            select {
            case v, ok := <-a:
                if ok {
                    c <- v
                } else {
                    a = nil
                }
            case v, ok := <-b:
                if ok {
                    c <- v
                } else {
                    b = nil
                }
            }
            
            if a == nil && b == nil {
                close(c)
                return
            }
        }
    }()
    return c
}

func main() {
    rand.Seed(time.Now().Unix())
    a := asChan(1, 3, 5, 7)
    b := asChan(2, 4, 6, 8)
    c := merge(a, b)
    for v := range c {
        fmt.Print(v)
    }
}