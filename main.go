package main

//Importing all tha file from the package main
//The Downloaded gorilla/mux package also lies in the main package

import (
	"encoding/json" //for working with json encoding and decoding
	"fmt"
	"log" //for console log
	"math/rand"
	"net/http" //For Working with Routes and server
	"strconv"

	"github.com/gorilla/mux" //the Mux for Routing and and managing the server
)

//Declaring a Movie Structure, with Json as Data structure for the api to send data
type Movie struct {
	ID       string    `json:"id"`       //movie ID, random genrates
	Isbn     string    `json:"isbn"`     //This used to test api on PostMan
	Title    string    `json:"title"`    //Title of movie
	Director *Director `json:"Director"` //Meta data og the movie
}

//Meta data of Movies
//Declaring a Director Structure, with Json and and Acossiate a pointer to movies Structure
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

//As we are not using any Database, so we are using a Slice to store movies data
var movies []Movie //A movie Slice

//function to Get Movies for User.
//Function take a two arguments, Respone and Request.
//Respone and Request exits in http libary.
func getmovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //setting up the incoming Json File formate using the Key in Json
	json.NewEncoder(w).Encode(movies)                  //changing the Slice to Json Format

}

//function for Deleting the The movies from the Slice
func deleteMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //setting up the incoming Json File formate using the Key in Json
	params := mux.Vars(r)                              //the Id Required For the Deletion of the movies Lies the request, getting it using the Mux.Vars

	for index, item := range movies { //finding  the movie using the For loop
		//The Index return the position and item return the Movie

		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //Deleting using the append Function,refer below for Explanation
			break                                                //Stopping the loop after last movie Shift
		}
	}
	json.NewEncoder(w).Encode(movies) //Returing the all the movies after deletion
}

/*Deleting using append Function :=
movies = append(movies[:index],movies[index+1:]...)
Here index is postion from which the movie have to remove
The pointer goes to the till the index position,
then it sees the next postion
It overwriters the movies till the 'index+1' to index position.
it does not delete the movies at index postion, it just overwrite the moves upto that index position.*/

func getmovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //setting up the incoming Json File formate using the Key in Json
	params := mux.Vars(r)                              //the Id Required For the Deletion of the movies Lies the request, getting it using the Mux.Vars

	for _, item := range movies { //here we do not require the index,so inplace use a blank
		if item.ID == params["id"] { //checking the movie id and requested id
			json.NewEncoder(w).Encode(item) //returning the movie data in json format
			return
		}
	}
}

func createmovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //setting up the incoming Json File formate using the Key in Json

	var movie Movie                               //making a temp variable of type movie Struture,for movies data that user entered
	_ = json.NewDecoder(r.Body).Decode(&movie)    //Decoding the entered movie Data, storing it in temp variable
	movie.ID = strconv.Itoa(rand.Intn(100000000)) //making a Random movie ID by rand.Intn and converting to string
	movies = append(movies, movie)                //adding the movie data entered to main Slice
	json.NewEncoder(w).Encode(movie)              //Returing the all the movies after Adding new movie data
}

func updatemovies(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	//if movie id exits,PUT data is empty it will be deleted
	//else if movie exits the older movie data delete and new movie with same id is appened
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...) //Deleting using the append Function

			var movie Movie //making a temp variable of type movie Struture,for movies data that user entered
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]        //new movie data with same id
			movies = append(movies, movie) //adding the movie data entered to main Slice
			json.NewEncoder(w).Encode(movie)
			return

		}
	}

}

func main() {

	//creating the Router for Server
	r := mux.NewRouter()

	//Creating the some  Movies data for the testing
	movies = append(movies, Movie{ID: "1", Isbn: "4512", Title: "Superman:Word of Man", Director: &Director{Firstname: "harry", Lastname: "got"}})
	movies = append(movies, Movie{ID: "2", Isbn: "412454", Title: "Thor:Love and Thunder", Director: &Director{Firstname: "Akash", Lastname: "Mercy"}})

	//Declaring the  Handle function for Routes
	//classifiying the route by Methods

	r.HandleFunc("/movies", getmovies).Methods("GET") //handler to get all the movies from the Slice

	r.HandleFunc("/movies/{id}", getmovie).Methods("GET") //handler to get the movies by id

	r.HandleFunc("/movies", createmovies).Methods("POST") //handler for adding the New movie and Data

	r.HandleFunc("/movies/{id}", updatemovies).Methods("PUT") //handler for updating the movie data using movie id

	r.HandleFunc("/Movies/{id}", deleteMovies).Methods("DELETE") //handler for deleting the movie bu movie id

	fmt.Printf("Server start at port 8000 \n") //console log for verification of server
	log.Fatal(http.ListenAndServe(":8000", r)) //logging out the server on port 8000.

}
