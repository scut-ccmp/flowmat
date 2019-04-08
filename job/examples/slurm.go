package main

import (
	"log"
	"os"

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

	jobMgt := job.JobManager("slurm")
	// submit job
	jobID, err := jobMgt.SubmitJob(conn, pathname)
	if err != nil {
		log.Fatal(err)
	}

	job.Check(jobMgt.CheckDoneFunc, conn, jobID)

	// recive files
	err = job.ReciveFiles(conn.Client, pathname, wd)
	if err != nil {
		log.Fatal(err)
	}
}
