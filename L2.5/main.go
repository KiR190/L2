/* 
Программа выведет:
error

Я языке Go интерфейсная переменная хранит тип и значение
Интерфейс считается nil только когда и тип, и значение равны nil
Если тип не nil, но само значение nil, то интерфейс не равен nil

error — это встроенный интерфейсный тип в Go, который имеет метод Error(). 
Когда мы присваиваем результат функции test() переменной err,
тип интерфейса становится *customError (тип=*customError, значение=nil),

Таким образом, хотя значение равно nil, тип интерфейса определён и не является nil, 
поэтому условие err != nil выполняется.

*/

package main

type customError struct {
    msg string
}

func (e *customError) Error() string {
    return e.msg
}

func test() *customError {
    // ... do something
    return nil
}

func main() {
    var err error
    err = test()
    if err != nil {
        println("error")
        return
    }
    println("ok")
}