package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

// checks if file is found and return it's format or error
func fileValidation(fileName string) (bool, string) {

	// get file metaData
	fileInfo, err := os.Stat(fileName)

	// split file-name and it's format
	fileNameParts := strings.Split(fileName, ".")

	// check errors
	if os.IsNotExist(err) {
		fmt.Println("file not found")
		return false, ""
	} else if err != nil {
		fmt.Println("error:", err)
		return false, ""
	} else if fileInfo.Size() == 0 {
		fmt.Println("file is empty")
		return false, ""
	} else if len(fileNameParts) == 1 {
		fmt.Println("file has no format")
		return false, ""
	}

	// possible file types
	textFileExtensions := map[string]bool{
		"json": true,
		// new formats can be added here
	}

	// check file's format
	if textFileExtensions[fileNameParts[len(fileNameParts)-1]] {
		return true, fileNameParts[len(fileNameParts)-1]
	}

	fmt.Println("unsupported file format. Please use any of next formats: TXT, CSV or JSON")
	return false, ""

}

// checks if file contatins valid JSON data and can be used for dictionary
func jsonValidation(inputString string) (inputMap map[string]string) {

	//v bool, byteValue []byte, err error

	byteValue, err := os.ReadFile(inputString)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully read json-file")
	}

	err = json.Unmarshal([]byte(byteValue), &inputMap)
	if err != nil || len(inputMap) == 0 {
		fmt.Println(err)
		fmt.Println("Data from JSON wasn't parsed")
	}

	return inputMap

}

// provides dictation from map for user and exits the program when finish
func dictation(m map[string]string) {

	for {

		scanner := bufio.NewScanner(os.Stdin)
		v := randomWord(m)

	jumpTo:

		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()

		// Check if the user wants to finish dictation
		if input == "exit" {
			break
		} else if input == v {
			fmt.Println("Correct! Try next word")
			continue
		} else {
			fmt.Println("Wrong, try again")
			goto jumpTo
		}

	}

}

// provides random word from map for dictation
func randomWord(m map[string]string) string {
	//k := rand.Intn(len(m))

	for i, v := range m {
		s := fmt.Sprintf("Переведи %s:", i)
		fmt.Println(s)
		return v
	}

	return ""

}

// Define your PostgreSQL connection parameters
const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func main() {

	// Connect to the PostgreSQL database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

jumpTo:

	fmt.Println("Enter file to parse:")

	// Scan for the next token (which is a line)
	scanner.Scan()

	inputString := scanner.Text()

	if scanner.Text() == "exit" {
		os.Exit(0)
	}

	fValid, _ := fileValidation(inputString)

	if fValid == false {
		goto jumpTo
	}

	m := jsonValidation(inputString)

	if len(m) == 0 {
		goto jumpTo
	}

	// Insert data into the PostgreSQL table
	for key, value := range m {
		_, err := db.Exec("INSERT INTO dictionary (key, value) VALUES ($1, $2)", key, value)
		if err != nil {
			fmt.Println("Error inserting data into the database:", err)
			return
		}
	}

	fmt.Println("You've uploaded next words from dictionary:")
	for key, value := range m {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println("Now the dictation starts.")

	dictation(m)

	fmt.Println("It was a pleasure to work with you, see ya!")

}
