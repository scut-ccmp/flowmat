package flowmat

import (
  "testing"
)

type PlusOneWF struct {
  In int
  Out int
}

func (t *PlusOneWF) Execute() {
  t.Out = t.In + 1
}

func TestSimpleWorkflow(t *testing.T) {
  tests := []struct {
    in int
    out int
  }{
    {0, 2},
    {-2, 0},
    {199, 201},
  }

  for _, test := range tests {
    p1 := NewProcess("p1", new(PlusOneWF))
    p2 := NewProcess("p2", new(PlusOneWF))

    wf := NewWorkflow("test_wf")
    wf.Add(p1)
    wf.Add(p2)
    wf.Connect("p1", "Out", "p2", "In")

    wf.ExposeIn("wfIn", "p1", "In")
    wf.ExposeOut("wfOut", "p2", "Out")

    wf.In("wfIn", test.in)
    wf.Run()

    got := wf.Out("wfOut")
    if got.(int) != test.out {
      t.Errorf("%d + 2 = %d", test.in, got)
    }
  }
}
