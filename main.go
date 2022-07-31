package main

// матеріали які я використовував
// https://habr.com/ru/post/502176/
// https://tutorialedge.net/golang/parsing-json-with-golang/
// https://golangify.com/parsing-string-date

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float64   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

type StringTrain struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float64
	ArrivalTime        string
	DepartureTime      string
}

const (
	price         string = "price"
	arrivalTime   string = "arrival-time"
	departureTime string = "departure-time"
)

func main() {

	f, _ := parsingJsonFile()

	for _, v := range f {
		fmt.Println(v)
	}

	//	... запит даних від користувача
	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	//	... обробка помилки
	//	... друк result
}

// переоприділяю метод UnmarshalJSON щоб зчитувалися time
func (t *StringTrain) UnmarshalJSON(data []byte) error {

	var stringTrain StringTrain
	err := json.Unmarshal(data, &stringTrain)
	if err != nil {
		fmt.Println(err)
	}
	parsingArrivalTime, err := time.Parse("00:50:00", stringTrain.ArrivalTime)
	if err != nil {
		fmt.Println(err)
	}
	parsingDepartureTime, err := time.Parse("00:50:00", stringTrain.DepartureTime)
	if err != nil {
		fmt.Println(err)
	}

	t.ArrivalTime = parsingArrivalTime.Format("00:50:00")
	t.DepartureTime = parsingDepartureTime.Format("00:50:00")
	t.TrainID = stringTrain.TrainID
	t.DepartureStationID = stringTrain.DepartureStationID
	t.ArrivalStationID = stringTrain.ArrivalStationID
	t.Price = stringTrain.Price

	return nil
}

func parsingJsonFile() ([]Train, error) {

	// відкриваєм наш jsonFile
	jsonFile, err := os.Open("data.json")
	// якщо виникає помилка при відкриванні, обробляєм помилку
	if err != nil {
		fmt.Println(err)
	}

	// перетворюю jsonFile в масив байтів щоб передати його нашому методу json.Unmarshal().
	byteFile, _ := ioutil.ReadAll(jsonFile)
	// створив файл типу слайс train
	var train []Train
	err = json.Unmarshal(byteFile, &train) // передаю спарсені дані у файл
	if err != nil {
		fmt.Println(err)
	}

	// закриваєм файл після того як функція відпрацювала
	defer jsonFile.Close()

	return train, nil
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	// ... код
	return nil, nil // маєте повернути правильні значення
}
