package main

import (
	"log"
	"os"
	"time"
	"fmt"

	"github.com/scut-ccmp/flowmat/job"
)


func main() {

	host := "202.38.220.15"
	port := "22"
	user := "unkcpz"
	pass := "sunshine"

	// get host public key
	// hostKey := getHostKey(host)

	conn, err := job.NewConnect(user, pass, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	dir := "/home/unkcpz/giida"
	prefix := "tmp"
	pathname, err := job.TempDir(conn.Client, dir, prefix)
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// send files
	// bugs need execute mod
	err = job.SendFiles(conn.Client, wd, pathname)
	if err != nil {
		log.Fatal(err)
	}

	// execute slurm job
	sess := conn.Session
	if err != nil {
		log.Fatal(err)
	}
	defer sess.Close()

	// sess.Stdout = os.Stdout
	// sess.Stderr = os.Stderr

	// Start remote shell
	// cmd := "cd " + pathname + "; module load vasp/5.4.4-impi-mkl; mpirun -n 4 vasp_std;"
	cmd := "cd " + pathname + ";sbatch job.sh"
 	out, err := sess.Output(cmd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(out))
	// c := 0
	// Run(func() (done bool) {
	// 	time.Sleep(2 * time.Second)
	// 	fmt.Println("doing")
	// 	c++
	// 	if c > 2 {
	// 		return true
	// 	}
	// 	return false
	// })

	// recive files
	err = job.ReciveFiles(conn.Client, pathname, wd)
	if err != nil {
		log.Fatal(err)
	}
}


func Run(proc func() bool) {
	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <- timeout:
			fmt.Println("timeout!")
			return
		case <- tick:
			done := proc()
			if done {
				fmt.Println("DONE")
				return
			}
		}
	}
}
