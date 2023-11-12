// main.go

/*
	This file contains the backend of HolidAI
	It primarily respons to JSON requests containing trip information
	Uses openAI to respond with a travel itenerary
*/

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sashabaranov/go-openai"
)

// Stores Trip Information
type trip struct {
	City   string
	Days   string
	Budget string
	Age    string
}

var tripObject trip

// Strings used for I/O
var prompt string
var plan string
var plans []string

// Generates prompt based on trip parameters
func makePrompt() string {

	var output string = ""

	// Trip information
	output += "Give me an itinerary for a trip to the city of " + tripObject.City + ". "
	output += "I am going for " + tripObject.Days + " days. "
	output += "My budget is " + tripObject.Budget + " $. "
	output += "I am " + tripObject.Age + ". "

	// Extra information to guide the prompt
	output += "Name every establishment. "
	output += "I only want to do one event every 3 hours, from 9AM to 9PM. "
	output += "I want to see the most iconic things in this city. "
	output += "Say how much every event will cost. "

	return output
}

// Uses openaI to generate a plan based off of the prompt
func makePlan() {
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

	plan = resp.Choices[0].Message.Content
}

// Homepage
func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

// Responds to a trip submission
func submit(response http.ResponseWriter, request *http.Request) {

	// Reads json file in request
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &tripObject)

	// Generates Prompt
	prompt = makePrompt()
	fmt.Println(prompt)

	// Makes Plan
	makePlan()
	fmt.Println(plan)
	plans := strings.Split(plan, "\n")

	// Output plan
	json.NewEncoder(response).Encode(&plans)
}

// Main function
func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/i", submit)

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
