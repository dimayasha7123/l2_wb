package main

import (
	"fmt"
	"time"
)

/*
=== Or channel ===

Реализовать функцию, которая будет объединять один или более done каналов в single канал если один из его составляющих
каналов закроется.
Одним из вариантов было бы очевидно написать выражение при помощи select, которое бы реализовывало эту связь,
однако иногда неизестно общее число done каналов, с которыми вы работаете в рантайме.
В этом случае удобнее использовать вызов единственной функции, которая, приняв на вход один или более or каналов,
реализовывала весь функционал.

Определение функции:
var or func(channels ...<- chan interface{}) <- chan interface{}

Пример использования функции:
sig := func(after time.Duration) <- chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
}()
return c
}

start := time.Now()
<-or (
	sig(2*time.Hour),
	sig(5*time.Minute),
	sig(1*time.Second),
	sig(1*time.Hour),
	sig(1*time.Minute),
)

fmt.Printf(“fone after %v”, time.Since(start))
*/

type orChannel func(channels ...<-chan interface{}) <-chan interface{}

func or(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		ret := make(chan interface{})
		go func(ret chan interface{}) {
			ret <- 1
		}(ret)
		return ret
	}
	if len(channels) == 1 {
		return channels[0]
	}

	ret := make(chan interface{})
	collector := make(chan interface{})
	for _, ch := range channels {
		go func(collector chan<- interface{}, ch <-chan interface{}) {
			collector <- <-ch //считаем, что ch это Done канал, т.е. любое отправленное значение или закрытие канала рассматривается как завершение работы канала
			fmt.Println("read from one channel...")
		}(collector, ch)
	}
	go func(collector <-chan interface{}, out chan interface{}) {
		count := 0
		for count < len(channels) {
			select {
			case value := <-collector:
				if count == 0 {
					out <- value
				}
				count++
			}
		}
		fmt.Println("out of collect func")
	}(collector, ret)

	return ret
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	var orFunc orChannel
	orFunc = or

	start := time.Now()
	<-orFunc(
		sig(5*time.Second),
		sig(4*time.Second),
		sig(3*time.Second),
		sig(2*time.Second),
		sig(2*time.Second),
	)

	fmt.Printf("done after %v\n", time.Since(start))

	time.Sleep(6 * time.Second)
}

// Вывод:
//
//	read from one channel...
//	done after 1.000196254s
//	read from one channel...
//	read from one channel...
//	read from one channel...
//	read from one channel...
//	out of collect func
//
//	Process finished with the exit code 0

// а это значит, что мы считали данные из всех каналов и таким образом не заблокировали функции, которые в них писали
