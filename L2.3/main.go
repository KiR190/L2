/* 

Интерфейс считается nil только если оба поля — и тип, и значение — равны nil.
Если тип присутствует, а значение — nil, — интерфейс не равен nil.

err := Foo() — err получает интерфейсное значение, у которого type = *os.PathError, value = nil.
fmt.Println(err) — fmt видит, что внутри интерфейса лежит конкретное значение nil (nil-указатель),
и печатает <nil>.
fmt.Println(err == nil) — сравнение с nil для интерфейса проверяет, равны ли оба поля (type и value).
У нас type != nil (есть *os.PathError), поэтому выражение даёт false.

Именно поэтому вывод:
<nil>
false

*/
package main

import (
    "fmt"
    "os"
)

func Foo() error {
    var err *os.PathError = nil
    return err
}

func main() {
    err := Foo()
    fmt.Println(err)        // <nil>
    fmt.Println(err == nil) // false
}