package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"slices"
	"sort"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	Init()

	//Create("Test task")
	//result = ReadLast()
	//result = ReadAll()
	//fmt.Printf("First available index is: %d", FirstAvailableIndex(ReadAll()))

	//fmt.Println(text)
	for {
		exec.Command("cls")
		MainPrint()
		input, _ := reader.ReadString('\n')

		switch input = strings.TrimSpace(input); input {
		case "1":
			fmt.Println("What would you like to name your task?")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			Create(input)
			_, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
		case "2":
			fmt.Println("Please write the ID of the task you would like to be updated: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			id, err := strconv.Atoi(input)
			if err != nil {
				log.Println("The value you entered is not valid, please try again...")
			} else {
				Update(id)
			}
			_, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
		case "3":
			fmt.Println("Please write the ID of the task you would like to be deleted: ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			id, err := strconv.Atoi(input)
			if err != nil {
				log.Println("The value you entered is not valid, please try again...")
			} else {
				Delete(id)
			}
			_, err = reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
		case "0":
			os.Exit(0)
		default:
			fmt.Println("Unknown option please try again. Press enter to continue...")
			_, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
		}
	}
}

func Init() {
	file, err := os.OpenFile("list.csv", os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		log.Printf("Failed to create file: %s", err)
	}
	defer file.Close()
}

func Create(name string) {
	file, err := os.OpenFile("list.csv", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	id := FirstAvailableIndex(ReadAll())

	record := []string{strconv.Itoa(id), name, "False"}

	if err := w.Write(record); err != nil {
		log.Fatalf("error writing record to file: %s", err)
	}

	defer w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Created a new task ID-%d called '%s'. Press enter to continue...\n", id, name)
}

func ReadLast() []string {
	file, err := os.OpenFile("list.csv", os.O_APPEND|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	var result []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		result = record
	}

	fmt.Println(result)

	return result
}

func ReadAll() [][]string {
	file, err := os.OpenFile("list.csv", os.O_APPEND|os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	var result [][]string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, record)
	}

	return result
}

func Update(id int) {
	var data [][]string = ReadAll()

	file, err := os.OpenFile("list.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	reader := bufio.NewReader(os.Stdin)
	w := csv.NewWriter(file)

	missing := true

	for _, row := range data {
		curr, _ := strconv.Atoi(row[0])
		if id == curr {
			fmt.Printf("Editing task id %d...\n", id)
			fmt.Println("Please enter a new name (keep blank if you want to keep the same name):")

			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input != "" {
				row[1] = input
				fmt.Println("Name successfully changed...")
			} else {
				fmt.Println("Name unchanged...")
			}

			fmt.Println("Do you want to change the task complete status? ('Y' to change, 'N' or blank to leave current status):")

			input, _ = reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input == "Y" || input == "y" {
				if row[2] == "True" {
					row[2] = "False"
				} else {
					row[2] = "True"
				}
				fmt.Println("Status changed successfully...")
			} else {
				fmt.Println("Status unchanged...")
			}

			w.WriteAll(data)
			missing = false
			fmt.Printf("Task with id %d has been sucessfully updated. Press enter to continue...\n", id)
			break
		}
	}
	if missing {
		fmt.Println("ID you inputed was not found...")
	}
}

func Delete(id int) {
	var data [][]string = ReadAll()

	file, err := os.OpenFile("list.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	missing := true

	for i, row := range data {
		curr, _ := strconv.Atoi(row[0])
		if id == curr {
			data = slices.Delete(data, i, i+1)
			w.WriteAll(data)
			missing = false
			fmt.Printf("Task with id %d has been sucessfully deleted. Press enter to continue...\n", id)
			break
		}
	}
	if missing {
		fmt.Println("ID you inputed was not found...")
	}
}

// Function returns the first index number that is not already in use
func FirstAvailableIndex(data [][]string) int {
	result := 0
	var indexes []int

	for _, task := range data {
		index, err := strconv.Atoi(task[0])
		if err != nil {
			log.Fatalf("Unable to parse index number: %s", err)
		}
		indexes = append(indexes, index)
	}
	sort.Slice(indexes, func(i, j int) bool {
		return indexes[i] < indexes[j]
	})

	for _, index := range indexes {
		result++
		if result != index {
			return result
		}
	}

	result++

	return result
}

func RowPrint(row []string) {
	par1, par2, par3 := row[0], row[1], row[2]

	fmt.Printf("%3s|%100s|%9s\n", par1, par2, par3)
}

func MainPrint() {
	var data [][]string
	header := []string{"ID", "Description", "Complete"}

	fmt.Println("Your current TODO list:")
	RowPrint(header)
	data = ReadAll()
	for _, task := range data {
		RowPrint(task)
	}
	fmt.Println("What would you like to do next? (1 - Create, 2 - Update, 3 - Delete, 0 - Exit)")
}
