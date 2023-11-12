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
	plan      string
}

// Generates prompt based on trip parameters
func (this trip) prompt() string {

	var prompt string = ""

	prompt += "Give me an itinerary for " + this.city + ". "
	prompt += "I am going for " + this.days + " days. "
	prompt += "My budget is " + this.budget + " $. "
	prompt += "I am " + this.traveller.age + ". "
	prompt += "I am " + this.traveller.style + ". "
	prompt += "Name every establishment. "
	prompt += "I only want to do one event every 3 hours, from 9AM to 9PM. "
	prompt += "I want to see the most iconic things in this city. "
	prompt += "// Say how much every event will cost."

	return prompt
}

func (this trip) makePlan() {

	var GPToutput string = ""

	// GPT request

	this.plan = GPToutput
}

func submission(response http.ResponseWriter, request *http.Request) {
	var newTrip trip
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &newTrip)
	json.NewEncoder(response).Encode(newTrip.plan)
}

func newItinerary(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Endpoint Hit: Itinerary")
	prompt := "Give me one word"
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
	fmt.Println(resp.Choices[0].Message.Content)
	fmt.Fprintf(response, resp.Choices[0].Message.Content)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/i", newItinerary)

	c := cors.New(cors.Options{
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	port := os.Getenv("PORT")

	if port == "" {
		port = "10000"
	}

	log.Fatal(http.ListenAndServe(":"+port, handler))

	/*
		bob := person{"Bob", "73", "Adventurous"}
		holiday := trip{"Halifax", "2", "35", bob}
		fmt.Println(holiday.prompt())
	*/

}
