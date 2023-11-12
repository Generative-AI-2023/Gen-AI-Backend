package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

func (this trip) plan() string {

	var plan string = this.prompt()

	// GPT request

	return plan
}

func submission(response http.ResponseWriter, request *http.Request) {
	var newTrip trip
	reqBody, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(reqBody, &newTrip)
	json.NewEncoder(response).Encode(newTrip)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)

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
