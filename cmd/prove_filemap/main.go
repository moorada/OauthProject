package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const PATH = "Person.json"

func main() {

	// create and populate a map from dummy JSON data

	dataStr := `{"domain2":"whoisExample","domain1":"whoisExample2"}`

	personMap := make(map[string]interface{})

	err := json.Unmarshal([]byte(dataStr), &personMap)

	if err != nil {
		panic(err)
	}

	jsonData, err := json.Marshal(personMap)

	if err != nil {
		panic(err)
	}

	// sanity check
	fmt.Println(string(jsonData))

	// write to JSON file

	var jsonFile *os.File

	if !fileExists(PATH) {

		jsonFile, err = os.Create(PATH)
		if err != nil {
			panic(err)
		}

	}
	jsonFile, err = os.Open(PATH)
	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		panic(err)
	}

	fmt.Println("JSON data written to ", jsonFile.Name())

	// Open our jsonFile
	//jsonFile, err := os.Open("./mapToFile/Person.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	jsonFile.Close()

	jsonFile, err = os.Open("./mapToFile/Person.json")

	fmt.Println("Successfully Opened users.json")

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("1", err)
	}

	personMap = make(map[string]interface{})

	fmt.Println("byteValue: ", string(byteValue))

	err = json.Unmarshal(byteValue, &personMap)
	if err != nil {
		fmt.Println("2", err)
	}

	fmt.Println(personMap)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

}

func main2() {
	addCouple(PATH, "dasdsadasdsa", "dasdssada")
}

func addCouple(path string, tld string, whoisAdmin string) {

	tldMap := make(map[string]interface{})

	jsonFile, err := os.Open(PATH)
	if err != nil {
		jsonFile, err = os.Create(PATH)
		if err != nil {
			panic(err)
		}
	} else {
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			panic(err)
		}
		_ = json.Unmarshal(byteValue, &tldMap)

	}

	tldMap[tld] = whoisAdmin
	jsonData, err := json.Marshal(tldMap)

	jsonFile, err = os.Create(PATH)
	if err != nil {
		panic(err)
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		panic(err)
	}

	jsonFile.Close()

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
