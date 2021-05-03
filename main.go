package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type structStreet struct {
	lon float64
	lat float64
	//number   int
	// street   string
	// unit     string
	// city     string
	// district string
	// region   string
	postcode string
}

// ExampleStats demonstrates how to read a full file and gather some stats.
// This is similar to `osmconvert --out-statistics`
func main() {
	mapDistrict := make(map[string]map[string]map[string]structStreet)

	csvfile, err := os.Open("./data/sweden.csv")
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	r := csv.NewReader(csvfile)
	r.Comma = ','
	r.Comment = '#'

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		lon, err := strconv.ParseFloat(record[0], 32)
		if err != nil {
			continue
		}
		lat, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			continue
		}
		if len(strings.TrimSpace(record[3])) == 0 ||
			len(strings.TrimSpace(record[5])) == 0 ||
			len(strings.TrimSpace(record[8])) == 0 {
			continue
		}

		if street, found := mapDistrict[record[5]]; found {
			if streetNo, found := street[record[3]]; found {
				if _, found := streetNo[record[2]]; found {
					continue
				} else {
					streetNo[record[2]] = structStreet{lon, lat, record[8]}
					street[record[3]] = streetNo
					mapDistrict[record[5]] = street
				}

			} else {
				mapStreetNo := make(map[string]structStreet)
				mapStreetNo[record[2]] = structStreet{lon, lat, record[7]}
				street[record[3]] = mapStreetNo
				mapDistrict[record[5]] = street
			}

		} else {
			mapStreetNo := make(map[string]structStreet)
			mapStreetNo[record[2]] = structStreet{lon, lat, record[8]}
			street := make(map[string]map[string]structStreet)
			street[record[3]] = mapStreetNo
			mapDistrict[record[5]] = street
		}
	}

	for city, streets := range mapDistrict {
		for street, numbers := range streets {
			for no, _ := range numbers {
				fmt.Printf("%s %s %s\n", city, street, no)
			}
		}

	}

	fmt.Println("Klar")
}
