package giida

import (
  "fmt"
  "testing"
  // "time"
)

type doubleOnce struct {
  In <-chan int
  Out chan<- int
}

func (proc *doubleOnce) Execute() {
  i := <-proc.In
  // time.Sleep(1 * time.Second)
  proc.Out <- 2 * i
}

// Test a simple task with multi inputs
func TestSimpleMultiTask(t *testing.T) {
  tests := []struct {
    in int
    expected int
  }{
    {12, 24},
    {5, 10},
    {0, 0},
  }

  for i, test := range tests {
    task := new(doubleOnce)

    name := fmt.Sprintf("doubler%d", i)
    proc := &Process{
      Name: name,
      task: task,
      Ports: map[string]chan int {
        "In": make(chan int),
        "Out": make(chan int),
      },
    }
    task.In = proc.Ports["In"]
    task.Out = proc.Ports["Out"]
    wait := proc.Run()
    proc.Ports["In"] <- test.in
    if got := <-proc.Ports["Out"]; got != test.expected {
      t.Errorf("%d * 2 != %d", test.in, got)
    }
    <-wait
  }
}

func TestTaskWithTwoInputs(t *testing.T) {
  tests := []struct {
    x int
    y int
    sum int
  }{
    {3, 38, 41},
    {3, 4, 7},
    {92, 4, 96},
    {-1, 1, 0},
  }

  for _, test := range tests {
    x := make(chan int)
    y := make(chan int)
    sum := make(chan int)

    task := new(adder)
    task.X = x
    task.Y = y
    task.Sum = sum

    proc := NewProcess("add", task)

    wait := proc.Run()

    x <- test.x
    y <- test.y
    got := <-sum
    if got != test.sum {
      t.Errorf("%d + %d != %d", test.x, test.y, test.sum)
    }
    <-wait
  }
}

type adder struct {
  X <-chan int
  Y <-chan int
  Sum chan<- int
}

func (p *adder) Execute() {
  x := <-p.X
  y := <-p.Y
  p.Sum <- x + y
}

// // Test a simple long running Task with one input
// func TestSimpleLongRunningTask(t *testing.T) {
//   tests := []struct {
//     in int
//     expected int
//   }{
//     {12, 24},
//     {5, 10},
//   }
//
//   in := make(chan int)
//   out := make(chan int)
//   proc := &doubler{
//     in,
//     out,
//   }
//
//   wait := Run(proc)
//
//   for _, test := range tests {
//     in <- test.in
//     got := <-out
//
//     if got != test.expected {
//       t.Errorf("%d != %d", got, test.expected)
//     }
//   }
//
//   close(in)
//   <-wait
// }
//
// type doubler struct {
//   In <-chan int
//   Out chan<- int
// }
//
// func (proc *doubler) Execute() {
//   for i := range proc.In {
//     proc.Out <- 2 * i
//   }
// }