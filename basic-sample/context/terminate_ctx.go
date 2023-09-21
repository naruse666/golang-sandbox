package main

import (
  "fmt"
  "context"
  "time"
)

func main() {
  ctx := context.Background()

  do(ctx)
}

func do(ctx context.Context) {
  ctx, ctxCancel := context.WithCancel(ctx)

  resultCh := make(chan string)

  go terminate(ctx, resultCh)

  resultCh <- "value1"
  resultCh <- "value2"

  fmt.Println("pre cancel")

  ctxCancel()

  time.Sleep(100 * time.Millisecond)

  fmt.Println("post cancel")
}

func terminate(ctx context.Context, ch <-chan string) {
  for {
    select {
    case <-ctx.Done():
      fmt.Println("terminate")
      return
    case result := <-ch:
      fmt.Println(result)
    }
  }
}
