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

	// submit job
	jobID, err := job.SubmitJob(conn, pathname)
	if err != nil {
		log.Fatal(err)
	}

	Check(func() (done bool) {
		state, err := job.FindJobState(conn, jobID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(state)
		if state == "NOJOBFOUND" {
			return true
		}
		return false
	})

	// recive files
	err = job.ReciveFiles(conn.Client, pathname, wd)
	if err != nil {
		log.Fatal(err)
	}
}


func Check(proc func() bool) {
	timeout := time.After(10000 * time.Second)
	tick := time.Tick(5 * time.Second)

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
