package cmd

import (
	"log"
	"os"

	"github.com/scut-ccmp/flowmat/gojob"
)

func RunVasp(host, port, user, pass, dir, prefix string) {
	conn, err := gojob.NewConnect(user, pass, host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	pathname, err := gojob.TempDir(conn.Client, dir, prefix)
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// send files
	// bugs need execute mod
	err = gojob.SendFiles(conn.Client, wd, pathname)
	if err != nil {
		log.Fatal(err)
	}

	jobMgt := gojob.JobManager("slurm")
	// submit job
	jobID, err := jobMgt.SubmitJob(conn, pathname)
	if err != nil {
		log.Fatal(err)
	}

	gojob.Check(jobMgt.CheckDoneFunc, conn, jobID)

	// recive files
	err = gojob.ReciveFiles(conn.Client, pathname, wd)
	if err != nil {
		log.Fatal(err)
	}
}
