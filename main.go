package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ExampleStats demonstrates how to read a full file and gather some stats.
// This is similar to `osmconvert --out-statistics`
func main() {
	mapDistrict := make(map[string]map[string]map[string]StructStreetDetail)

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

		street := &record[3]
		city := &record[5]
		postalcode := &record[8]
		streetNo := &record[2]

		if len(strings.TrimSpace(*city)) == 0 ||
			len(strings.TrimSpace(*street)) == 0 ||
			len(strings.TrimSpace(*postalcode)) == 0 {
			continue
		}

		if mapStreet, found := mapDistrict[*city]; found {
			if mapStreetNo, found := mapStreet[*street]; found {
				if _, found := mapStreetNo[*streetNo]; found {
					continue
				} else {
					mapStreetNo[*streetNo] = StructStreetDetail{lon, lat, *postalcode}
					mapStreet[*street] = mapStreetNo
					mapDistrict[*city] = mapStreet
				}

			} else {
				mapStreetNo := make(map[string]StructStreetDetail)
				mapStreetNo[*streetNo] = StructStreetDetail{lon, lat, *postalcode}
				mapStreet[*street] = mapStreetNo
				mapDistrict[*city] = mapStreet
			}

		} else {
			mapStreetNo := make(map[string]StructStreetDetail)
			mapStreetNo[*streetNo] = StructStreetDetail{lon, lat, *postalcode}
			mapStreet := make(map[string]map[string]StructStreetDetail)
			mapStreet[*street] = mapStreetNo
			mapDistrict[*city] = mapStreet
		}
	}

	arr := makeStruct(&mapDistrict)

	fmt.Println(len(arr))
	fmt.Println("Klar")
}

func makeStruct(mapDistrict *map[string]map[string]map[string]StructStreetDetail) []StructCity {
	col := connectMongo("sweden")
	var arrCity = []StructCity{}
	for city, streets := range *mapDistrict {
		var arrStreet = []StructStreet{}
		for street, numbers := range streets {
			var arrStreetNo = []StructStreetNo{}
			for no, detail := range numbers {
				arrStreetNo = append(arrStreetNo, StructStreetNo{no, detail})
			}
			arrStreet = append(arrStreet, StructStreet{street, arrStreetNo})

		}

		arrCity = append(arrCity, StructCity{primitive.NewObjectID(), city, arrStreet})
	}

	citys := make([]interface{}, len(arrCity))
	for i, s := range arrCity {
		citys[i] = s
	}

	_, err := col.InsertMany(context.Background(), citys)

	if err != nil {
		log.Fatal(err)
	}

	return arrCity
}

func connectMongo(country string) *mongo.Collection {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)

	quickstartDatabase := client.Database("street-golang")
	return quickstartDatabase.Collection(country)
}
