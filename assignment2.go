package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

func main() {

	//task #1
	//first we take the time for comparison and measure the performance
	startTime := time.Now()

	//asking the user for the number
	var number int
	fmt.Println("Enter a positive integer: ")
	fmt.Scanln(&number)

	//create the slice for the factorial
	slice := make([]int, number)
	for i := 0; i < number; i++ {
		slice[i] = i + 1
	}

	//define the number of go routines to calculate the factorial depending on the number of threads
	numThreads := runtime.NumCPU()
	numGoroutines := func() int {
		if 2*number < numThreads {
			return number / 2
		} else {
			return numThreads
		}
	}

	//calculate the size of the divisions of the slice
	size := int(math.Floor(float64(len(slice)) / float64(numGoroutines())))

	//create the channel to collect the factorial results
	chanFactorial := make(chan uint64, numGoroutines())

	//create a wait group that wait for the routines to be done
	var wg sync.WaitGroup

	//Launch the goroutines
	for i := 0; i < numGoroutines(); i++ {
		wg.Add(1)
		startIndex := int(i * size)
		endIndex := (i + 1) * int(size)
		if i == numGoroutines()-1 {
			endIndex = len(slice)
		}
		go factorial(slice[startIndex:endIndex], chanFactorial, &wg)

	}

	//close the result channel
	go func() {
		wg.Wait()
		close(chanFactorial)
	}()

	//collect the result channel to calculate the final result
	factorialResult := uint64(1)
	for result := range chanFactorial {
		factorialResult *= result
	}

	//the end time compared to the starting time to get the duratioin of the code execution
	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("The factorial of %d is %d\n", number, factorialResult)
	fmt.Printf("The execution time: %s\n\n", duration)

	//task #2
	//first we have to create the instances of the rectangle and the circle
	rectangle := Rectangle{Width: 3, Height: 4}
	circle := Circle{Radius: 5}

	fmt.Printf("Area rectangle: %.2f\nPerimeter Rectangle: %2.f\n\n", rectangle.Area(), rectangle.Perimeter())
	fmt.Printf("Area circle: %.2f\nPerimeter circle: %2.f\n\n", circle.Area(), circle.Perimeter())

	//task #3
	//we simulate reading a file
	fileName := "SecretOfLife.go"
	data, err := readFIle(fileName)

	//error handling
	if err != nil {
		switch e := err.(type) {
		case FileNotFoundError:
			fmt.Println("Error caught: ", e)
		default:
			fmt.Println("Unexpected error: ", e)
		}
	}

	//if there's no error
	fmt.Println("File data:", data)
	fmt.Println()

	//Task #4
	//We first read the JSON file containing products informatioin
	inputFileName := "stock.json"
	info, err := os.ReadFile(inputFileName)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	//Unmarshall the JSON data into a Slice of products
	var products []Product
	err = json.Unmarshal(info, &products)
	if err != nil {
		fmt.Println("error unmarshalling JSON data:", err)
		return
	}

	//you can filter the products through different specifications
	filteredProducts := filteredProducts(products, func(p Product) bool {
		return p.Quantity > 15
	})

	//sorting products by name
	sort.Slice(filteredProducts, func(i, j int) bool {
		return filteredProducts[i].Name < filteredProducts[j].Name
	})

	//modify the product's brand
	if len(filteredProducts) > 0 {
		filteredProducts[0].Brand = "Off Brand"
	}

	//write the modified data to a new JSON file
	outputFileName := "modified_products.json"
	outputData, err := json.MarshalIndent(filteredProducts, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON data: ", err)
		return
	}

	err = os.WriteFile(outputFileName, outputData, 0644)
	if err != nil {
		fmt.Println("Error writing to JSON file: ", err)
		return
	}

	//Print a message indicating the operation completion and the output file name
	fmt.Println("Operation completed. Modified data written to ", outputFileName)
	fmt.Println()
}

func factorial(slice []int, chanFactorial chan uint64, wg *sync.WaitGroup) int {
	defer wg.Done()

	fact := uint64(1)

	for _, value := range slice {
		fact *= uint64(value)
	}

	chanFactorial <- fact

	return 0

}

// the interface Shape has the methods for the area and the perimeter
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle implements the Shape interface and contains the width and the height
type Rectangle struct {
	Width  float64
	Height float64
}

// this Area calculates de area of the rectangle
func (rect Rectangle) Area() float64 {
	return rect.Width * rect.Height
}

// this Perimeter calculates de perimeter of a rectangle
func (rect Rectangle) Perimeter() float64 {
	return 2*rect.Width + 2*rect.Height
}

// the same as for rectangle but for the cicrle shape
type Circle struct {
	Radius float64
}

func (cir Circle) Area() float64 {
	return math.Pi * cir.Radius * cir.Radius
}

func (cir Circle) Perimeter() float64 {
	return 2 * math.Pi * cir.Radius
}

// Custom error type
type FileNotFoundError struct {
	FileName string
}

// error returns the error message
func (e FileNotFoundError) Error() string {
	return fmt.Sprintf("The %s file was not found", e.FileName)
}

// this function tries to read an inexistent file to return an error
func readFIle(fileName string) ([]byte, error) {
	return nil, FileNotFoundError{FileName: fileName}
}

// Product tepresents a structure of a product in JSON
type Product struct {
	Name     string `json:"name"`
	Brand    string `json:"brand"`
	Quantity int    `json:"quantity"`
}

// filtering products based on a given condition
func filteredProducts(products []Product, condition func(Product) bool) []Product {
	var result []Product
	for _, product := range products {
		if condition(product) {
			result = append(result, product)
		}
	}
	return result
}
