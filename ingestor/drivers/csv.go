package drivers

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

type driver interface {
	Parse(file *os.File, ch chan interface{}) error
}

type csvdriver struct{}

var driversmap = map[string]driver{
	"CSV": csvdriver{},
}

func Driver(name string) driver {
	return driversmap[name]
}

func (d csvdriver) Parse(file *os.File, ch chan interface{}) error {
	reader := csv.NewReader(bufio.NewReader(file))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
			return error
		}

		fmt.Println(line)
	}
	return nil
}
