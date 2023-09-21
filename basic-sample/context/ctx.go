package main

import (
  "fmt"
  "context"
)

func doSomething(ctx context.Context) {
  fmt.Printf("key1 is %s\n", ctx.Value("key1"))

  newCtx := context.WithValue(ctx, "key1", "value2")

  anotherFunc(newCtx)

  fmt.Printf("key1 is %s\n", ctx.Value("key1"))
}

func anotherFunc(ctx context.Context) {
  fmt.Printf("key1 is %s\n", ctx.Value("key1"))
}

func main() {
  ctx := context.Background()

  ctx = context.WithValue(ctx, "key1", "value1")

  doSomething(ctx)
}
