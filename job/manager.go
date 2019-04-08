package job

type Manager interface{
  FindJobState(c *Conn, id int) (state string, err error)
  SubmitJob(c *Conn, path string) (id int, err error)
}

func JobManager(name string) Manager {
  table := map[string]Manager {
    "slurm": NewSlurmMgt(),
  }
  return table[name]
}
