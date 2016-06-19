package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func dump(file os.FileInfo) {
	err := os.Setenv("MDB_JET3_CHARSET", "cp1256")
	if err != nil {
		log.Println(err)
	}

	cmd := exec.Command("mdb-tables", "bok/"+file.Name())

	output, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
	}

	r := bufio.NewReader(output)

	line, _, err := r.ReadLine()
	if err != nil {
		log.Println(err)
	}

	l := string(line)
	tables := strings.Fields(l)
	for _, table := range tables {
		if strings.HasPrefix(table, "b") || strings.HasPrefix(table, "t") {

			cmd = exec.Command("mdb-export", "-I", "mysql", "bok/"+file.Name(), table)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				log.Println(err)
			}

		}
	}

}

func main() {

	fmt.Println("initialising the process")

	files, err := ioutil.ReadDir("bok")
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		dump(file)
	}

}
