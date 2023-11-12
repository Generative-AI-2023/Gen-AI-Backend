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

var holiday trip
var prompt string
var plan string
var plans []string

// Generates prompt based on trip parameters
func makePrompt() string {

	var output string = ""

	output += "Give me an itinerary for a trip to the city of " + holiday.City + ". "
	output += "I am going for " + holiday.Days + " days. "
	output += "My budget is " + holiday.Budget + " $. "
	output += "I am " + holiday.Age + ". "

	/*for i := 0; i < len(holiday.traveller.wants); i++ {

		if holiday.traveller.wants[i] == "true" {
			output += "I want " + holiday.traveller.styles[i] + ". "
		} else {
			output += "I do not want " + holiday.traveller.styles[i] + ". "
		}
	}*/

	output += "Name every establishment. "
	output += "I only want to do one event every 3 hours, from 9AM to 9PM. "
	output += "I want to see the most iconic things in this city. "
	output += "Say how much every event will cost. "

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

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Homepage!")
	fmt.Println("Endpoint Hit: Homepage")
}

func submit(response http.ResponseWriter, request *http.Request) {

	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &holiday)

	// Generates Prompt
	prompt = makePrompt()
	fmt.Println(prompt)

	// Makes Plan
	plan = makePlan()
	plans := strings.Split(plan, "\n")

	// Output plan
	json.NewEncoder(response).Encode(&plans)
}

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
