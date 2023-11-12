package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sashabaranov/go-openai"
)

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

// Stores Person Information
type person struct {
	name  string
	age   string
	style string
}

// Stores Trip Information
type trip struct {
	city      string
	days      string
	budget    string
	traveller person
}

var holiday trip
var prompt string
var plan string
var plans [1]string

// Generates prompt based on trip parameters
func makePrompt() string {

	var output string = ""

	output += "Give me an itinerary for a trip to the city of " + holiday.city + ". "
	output += "I am going for " + holiday.days + " days. "
	output += "My budget is " + holiday.budget + " $. "
	output += "I am " + holiday.traveller.age + ". "
	output += "I am " + holiday.traveller.style + ". "
	output += "Name every establishment. "
	output += "I only want to do one event every 3 hours, from 9AM to 9PM. "
	output += "I want to see the most iconic things in this city. "
	output += "Say how much every event will cost. "
	output += "Say the word banana after"

	return output
}

func makePlan() string {
	client := openai.NewClient(os.Getenv("API_KEY"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		panic(err)
	}

	return resp.Choices[0].Message.Content
}

func submit(response http.ResponseWriter, request *http.Request) {

	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &holiday)

	prompt = makePrompt()
	fmt.Println(prompt)

	plan = makePlan()
	fmt.Println(plan)

	plans[0] = plan
	json.NewEncoder(response).Encode(&plans)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/i", submit)

	bob := person{"Bob", "73", "Adventurous"}
	holiday = trip{"London, Enlgand", "2", "35", bob}
	prompt = makePrompt()
	prompt = "Give me a word"
	fmt.Println(prompt)

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	port := os.Getenv("PORT")

	if port == "" {
		port = "10000"
	}

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
