package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-time.After(2 * time.Second)
		log.Println("cancelling")
		cancel()
	}()
	Run(ctx, task1, task2)
	<-time.After(5 * time.Second)
}

func main2() {
	ids := []string{"x", "1", "y", "2"}
	str := stringifyUserIDs(ids)
	str2 := stringifyUserIDs_o(ids)
	log.Println(str)
	log.Println(str2)
}

func stringifyUserIDs_o(IDs []string) string {
	return "[\"" + strings.Join(IDs, "\",\"") + "\"]"
}

func stringifyUserIDs(IDs []string) string {
	sparator := "\",\""
	joined_ids := strings.Join(IDs, sparator)
	return fmt.Sprintf(`["%v"]`, joined_ids)
}

func task1(ctx context.Context, errCh chan<- error) {
	i := 0
	for {
		i++
		log.Println("task 1")
		<-time.After(1 * time.Second)
		if i > 30 {
			break
		}
		select {
		case <-ctx.Done():
			errCh <- errors.New("cancelled")
			return
		default:
		}
	}
	errCh <- nil
}

func task2(ctx context.Context, errCh chan<- error) {
	i := 0
	for {
		i++
		log.Println("task 2")
		<-time.After(500 * time.Millisecond)
		if i > 60 {
			break
		}
		select {
		case <-ctx.Done():
			errCh <- errors.New("cancelled")
			return
		default:
		}
	}
	errCh <- nil
}

func Run(ctx context.Context, tasks ...func(context.Context, chan<- error)) []error {
	errCh := make(chan error)
	errs := make([]error, 0)
	for _, t := range tasks {
		go t(ctx, errCh)
	}
	for range tasks {
		select {
		case err := <-errCh:
			if err != nil {
				errs = append(errs, err)
			}
		case <-ctx.Done():
			return errs
		}
	}
	close(errCh)
	if len(errs) == 0 {
		return nil
	}
	return errs
}
