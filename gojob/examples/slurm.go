package main

import (
	"log"
	"os"
	"fmt"
	"text/template"

	"github.com/scut-ccmp/flowmat/gojob"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("$HOME/.config/gojob/")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	user := viper.GetString("server.user")
	pass := viper.GetString("server.password")
	host := viper.GetString("server.host")
	port := viper.GetString("server.port")

	dir := viper.GetString("file.tempDir")
	prefix := viper.GetString("file.dirPrefix")

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

	// prepare _job.sh file
	slurm := gojob.SlurmParameter{
		Name: "flowmat",
		NProc: viper.GetString("job.nproc"),
		NCom: viper.GetString("job.ncom"),
		Partion: viper.GetString("job.partion"),
		Prepend: viper.GetString("job.prepend"),
		ExecCmd: viper.GetString("job.exec"),
	}

	t, err := template.New("slurm job").Parse(gojob.SlurmTmpl)
	if err != nil {
		log.Fatal("Parse tmpl: ", err)
		return
	}

	f, err := os.Create("_job.sh")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = t.Execute(f, slurm)
	if err != nil {
		log.Fatal(err)
		return
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
