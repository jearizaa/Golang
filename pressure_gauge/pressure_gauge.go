package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type PressureGauge struct {
	ch chan struct{}
}

type message struct {
	s     string
	resCh chan string
	errCh chan error
}

func New(limit int) *PressureGauge {
	ch := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		ch <- struct{}{}
	}
	return &PressureGauge{ch}
}

func (pg *PressureGauge) Process(f func(message), m message) error {
	select {
	case <-pg.ch:
		f(m)
		return nil
	default:
		return errors.New("no more capacity")
	}
}

func main() {
	start := time.Now()
	length := 10000
	numbers := make([]string, length)
	for i := range numbers {
		numbers[i] = strconv.Itoa(i)
	}

	responses := make([]string, 0, length)
	resErrors := make([]error, 0)

	pg := New(1000)
	var wg sync.WaitGroup
	wg.Add(length)
	done := make(chan struct{})
	resCh := make(chan string, length)
	errCh := make(chan error)

	go func() {
		for {
			select {
			case <-done:
				fmt.Println("salio")
				return
			case res := <-resCh:
				responses = append(responses, res)
				pg.ch <- struct{}{}
				wg.Done()
			case err := <-errCh:
				resErrors = append(resErrors, err)
				pg.ch <- struct{}{}
				wg.Done()
			}
		}
	}()

	for i := 0; i < len(numbers); i++ {
		err := pg.Process(function, message{
			s:     numbers[i],
			resCh: resCh,
			errCh: errCh,
		})
		if err != nil {
			i--
		}
	}

	wg.Wait()
	close(done)
	close(resCh)
	close(errCh)

	fmt.Printf("responses %d, errors %d\n", len(responses), len(resErrors))
	fmt.Printf("Took %s", time.Since(start))
}

func function(m message) {
	go func() {
		res, err := method(m.s)
		if err != nil {
			m.errCh <- err
			return
		}
		m.resCh <- res
	}()
}

func method(s string) (string, error) {
	time.Sleep(15 * time.Millisecond)
	i, err := strconv.Atoi(s)
	if err != nil {
		return "", err
	}
	if i%5 == 0 {
		return "", fmt.Errorf("error with %d", i)
	}
	return s + "-Go", nil
}
