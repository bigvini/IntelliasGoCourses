package main

// матеріали які я використовував
// https://habr.com/ru/post/502176/
// https://tutorialedge.net/golang/parsing-json-with-golang/
// https://golangify.com/parsing-string-date
// https://www.geeksforgeeks.org/time-time-date-function-in-golang-with-examples/
// тут зразок викоритстанн date https://goplay.tools/snippet/ptIFIzj0tDR

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
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

var (
	departureStation string
	arrivalStation   string
	criteria         string
	result           []Train
)

var (
	UnsupportedCriteria      = errors.New("unsupported criteria")
	EmptyStation             = errors.New("empty station")
	EmptyDepartureStation    = errors.New("empty departure station")
	EmptyArrivalStation      = errors.New("empty arrival station")
	BadStationInput          = errors.New("bad station input")
	BadDepartureStationInput = errors.New("bad departure station input")
	BadArrivalStationInput   = errors.New("bad arrival station input")
)

func main() {

	fmt.Println("Enter departure station ID")
	departureStation = userInput()
	fmt.Println("Enter arrival station ID")
	arrivalStation = userInput()
	fmt.Println("Enter criteria")
	criteria = userInput()

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		err = fmt.Errorf("entered incorrect parameters: %v", err)
		fmt.Println(err)
		return
	}

	for _, v := range result {
		// fmt.Printf("TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, "+
		// 	"Price: %v, ArrivalTime: time.Date(%v) , DepartureTime: %v \n",
		// 	v.TrainID, v.DepartureStationID, v.ArrivalStationID, v.Price,
		// 	v.ArrivalTime.Format("15:04:05"), v.DepartureTime.Format("15:04:05"))

		fmt.Printf("TrainID: %v, DepartureStationID: %v, ArrivalStationID: %v, Price: %v, ArrivalTime: %v,"+
			" DepartureTime: %v \n", v.TrainID, v.DepartureStationID, v.ArrivalStationID, v.Price,
			v.ArrivalTime.Format("15:04:05"), v.DepartureTime.Format("15:04:05"),
		)
	}
	//	... запит даних від користувача
	//result, err := FindTrains(departureStation, arrivalStation, criteria))
	//	... обробка помилки
	//	... друк result
}

func checkErrorStation(s string) (int, error) {

	if s == "" {
		return 0, EmptyStation
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0, BadStationInput
	}
	if result < 1 {
		return 0, BadStationInput
	}
	return result, nil
}

func userInput() (userInput string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput = scanner.Text()
	return userInput
}

// переоприділяю метод UnmarshalJSON щоб зчитувалися time
func (t *Train) UnmarshalJSON(data []byte) error {

	var stringTrain StringTrain
	err := json.Unmarshal(data, &stringTrain)
	if err != nil {
		fmt.Println(err)
	}

	// парсінг строки з конвертацією в time
	parsingArrivalTime, err := time.Parse("15:04:05", stringTrain.ArrivalTime)
	if err != nil {
		fmt.Println(err)
	}
	parsingDepartureTime, err := time.Parse("15:04:05", stringTrain.DepartureTime)
	if err != nil {
		fmt.Println(err)
	}
	// присвоюю значення у форматі time

	t.ArrivalTime = time.Date(0, time.January, 1, parsingArrivalTime.Hour(), parsingArrivalTime.Minute(), parsingArrivalTime.Second(), 0, time.UTC)
	t.DepartureTime = time.Date(0, time.January, 1, parsingDepartureTime.Hour(), parsingDepartureTime.Minute(), parsingDepartureTime.Second(), 0, time.UTC)

	// присвоєю значення зміним структурі без змін
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

func priceSort(f []Train) []Train {
	// сортую по ціні від найменшоі
	sort.SliceStable(f, func(i, j int) bool {
		return f[i].Price < f[j].Price
	})
	return f[:3]
}

func sortDepartureTime(f []Train) []Train {

	// var sortArrival []Train

	sort.SliceStable(f, func(i, j int) bool {
		return f[i].DepartureTime.Before(f[j].DepartureTime)
	})

	return f[:3]
}

func sortArrivalTime(f []Train) []Train {

	// var sortArrival []Train

	sort.SliceStable(f, func(i, j int) bool {
		return f[i].ArrivalTime.Before(f[j].ArrivalTime)
	})

	return f[:3]
}

func findStation(departureStation, arrivalStation int) []Train {
	f, _ := parsingJsonFile()
	var sliceTrain []Train = nil

	for _, v := range f {
		if v.ArrivalStationID == arrivalStation && v.DepartureStationID == departureStation {
			sliceTrain = append(sliceTrain, v)
		}

	}

	return sliceTrain

}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	departure, err := checkErrorStation(departureStation)
	if err != nil {
		if errors.Is(err, BadStationInput) {
			return nil, BadDepartureStationInput
		}
		return nil, EmptyDepartureStation
	}

	arrival, err := checkErrorStation(arrivalStation)
	if err != nil {
		if errors.Is(err, BadStationInput) {
			return nil, BadArrivalStationInput
		}
		return nil, EmptyArrivalStation
	}

	var sortTrain []Train

	switch criteria {
	case price:
		return priceSort(findStation(departure, arrival)), nil
	case arrivalTime:
		return sortArrivalTime(findStation(departure, arrival)), nil
	case departureTime:
		return sortDepartureTime(findStation(departure, arrival)), nil
	default:
		return sortTrain, UnsupportedCriteria
	}

	// маєте повернути правильні значення
}
