package giida

import (
  "fmt"
	"time"
)

// Tasker is a unit run a Execute
type Tasker interface {
  Execute()
}

type Processer interface{
  Run() Wait
}

type Process struct {
  Name string
  task  Tasker
	Ports map[string]chan int
}

func NewProcess(name string, task Tasker) *Process {
  p := &Process{
    Name: name,
    task: task,
  }
  return p
}

// func portToTask(p *Process) err {
// 	// Set value to task's fields
// 	t = p.task
// 	val := reflect.ValueOf(p.task).Elem()
// 	for i := 0; i < val.NumField(); i++ {
// 		field := val.Field(i)
// 		fieldType := field.Type()
// 		fieldName := filed.Name
// 		p.Ports[fieldName] =
// 	}
// }

// Done notifies that the process is finished
type Done struct{}

// Wait is a channel signalling of a completion
type Wait chan struct{}

func (p *Process) Run() Wait {
	t := p.task
	wait := make(Wait)
	go func() {
		// fmt.Printf("%s | Running %s\n", timeStamp(), p.Name)
		t.Execute()

		wait <- Done{}
		// fmt.Printf("%s | %s Finished\n", timeStamp(), p.Name)
	}()
	return wait
}

func timeStamp() string {
	t := time.Now()
	return fmt.Sprintf(t.Format("2006/01-02/15:04:05"))
}

// type InputGuard struct {
//   ports map[string]bool
//   complete int
// }
//
// //
// func NewInputGuard(ports ...string) *InputGuard {
//   portMap := make(map[string]bool, len(ports))
//   for _, p := range ports {
//     portMap[p] = false
//   }
//   return &InputGuard{portMap, 0}
// }
//
// // Complete is called when a port is closed and returns true when all the ports have been closed
// func (g *InputGuard) Complete(port string) bool {
//   if !g.ports[port] {
//     g.ports[port] = true
//     g.complete++
//   }
//   return g.complete >= len(g.ports)
// }