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

func dump(file os.FileInfo) (SQLFile string, err error) {
	err = os.Setenv("MDB_JET3_CHARSET", "cp1256")
	if err != nil {
		return "", err
	}

	cmd := exec.Command("mdb-tables", "bok/"+file.Name())

	output, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	r := bufio.NewReader(output)
	output.Close()

	line, _, err := r.ReadLine()
	if err != nil {
		return "", err
	}

	l := string(line)
	tables := strings.Fields(l)

	var f *os.File
	defer f.Close()

	fn := strings.Split(file.Name(), ".")[0]

	err = os.MkdirAll("sql", 0755)
	if err != nil {
		return "", err
	}

	f, err = os.Create("sql/" + fn + ".sql")
	if err != nil {
		return "", err
	}

	f.Write([]byte("BEGIN;"))

	for _, table := range tables {

		export := exec.Command("mdb-export", "-I", "mysql", "bok/"+file.Name(), table)

		//export.Stdout = os.Stdout
		//export.Stderr = os.Stderr

		out, err := export.StdoutPipe()
		if err != nil {
			return "", err
		}

		if err = export.Start(); err != nil {
			return "", err
		}

		rd := bufio.NewReader(out)

		if err = export.Wait(); err != nil {
			return "", err
		}

		_, err = rd.WriteTo(f)
		if err != nil {
			return "", err
		}

		out.Close()

	}

	f.Write([]byte("COMMIT;"))
	return f.Name(), nil
}

func main() {

	files, err := ioutil.ReadDir("bok")
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		SQLFile, err := dump(file)
		if err != nil {
			log.Fatal(err)

		}
		fmt.Printf("Completed SQL file: %v", SQLFile)

	}

}
