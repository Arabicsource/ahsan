package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nwaples/rardecode"
)

// extract a rar file and save content into a bok directory
func extract(f string) (err error) {

	err = os.Setenv("MDB_JET3_CHARSET", "cp1256")
	if err != nil {
		return err
	}

	if _, err := os.Stat("bok"); os.IsNotExist(err) {
		if err := os.MkdirAll("bok", 0755); err != nil {
			log.Println(err)
			return err
		}
	} else if err != nil {
		return err
	}

	//	cmd := exec.Command("unrar", "e", "../downloads/"+f)
	//	cmd.Dir = "bok"
	//	cmd.Stdout = os.Stdout
	//	cmd.Stderr = os.Stderr
	//	if err := cmd.Run(); err != nil {
	//		log.Println(err)
	//		return err
	//	}

	// id of the book
	id := strings.Split(f, ".")[0]

	nfn := id + ".bok"

	rarfile, err := os.Open(filepath.Join("downloads", f))
	if err != nil {
		return err
	}

	rdr, err := rardecode.NewReader(rarfile, "")
	if err != nil {
		return err
	}

	_, err = rdr.Next()
	if err != nil {
		return err
	}

	newbok, err := os.Create(filepath.Join("bok", nfn))
	defer newbok.Close()
	if err != nil {
		return err
	}

	fmt.Printf("%v being extracted ....", nfn)
	_, err = io.Copy(newbok, rdr)
	if err != nil {
		newbok.Close()
		fmt.Printf("error!")
		rm := exec.Command("rm", nfn)
		rm.Stderr = os.Stderr
		if err = rm.Run(); err != nil {
			fmt.Printf(" - Could not remove file [" + nfn + "]\n")
		}

		return err
	}
	fmt.Printf("done!\n")

	return nil
}

func main() {

	files, err := ioutil.ReadDir("downloads")
	if err != nil {
		log.Println(err)
	}

	for _, file := range files {
		if err = extract(file.Name()); err != nil {
			log.Println(err)
			continue
		}
	}
}
