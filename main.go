package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	// "github.com/gorilla/mux"
)

type Movie struct {
	ID string `json:"id"`
	Isbn string `json:"isbn"`
	Title string `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName string `json:"lastname"`
}

var movies []Movie

func getMovies(res http.ResponseWriter, req *http.Request){
	if req.URL.Path != "/movies"{
		http.Error(res, "route not found", http.StatusNotFound)
		return
	}
	if req.Method != http.MethodGet {
		http.Error(res, "method not supported", http.StatusNotFound)
		return
	}
	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(movies)
	if err != nil {
		http.Error(res, "Error formatting JSON", http.StatusInternalServerError)
		return
	}
}

func getMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	movieId := chi.URLParam(req, "id")
	if movieId == ""{
		http.Error(res, "please send the movie ID", http.StatusNotFound)
		return
	}
	for _, value := range movies {
		if value.ID == movieId {
			res.WriteHeader(http.StatusOK)
			response := map[string]Movie{"movie": value}
			json.NewEncoder(res).Encode(response)
			return
		}
	}
	res.WriteHeader(http.StatusOK)
	response := map[string]string{"error": "Movie not found"}
	json.NewEncoder(res).Encode(response)
}

func createMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	var movie Movie
	err := json.NewDecoder(req.Body).Decode(&movie)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(map[string]string{"error": "invalid request payload"})
		return 
	}
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)
	res.WriteHeader(http.StatusOK)
	// response := map[string]string{"message": "Movie created successfully"}
	// json.NewEncoder(res).Encode(response)
	// movieRes := map[string]Movie{"movie": movie}
	// json.NewEncoder(res).Encode(movieRes)
	response := map[string]interface{}{
		"message": "Movie created successfully",
		"movie": movie,
	}
	json.NewEncoder(res).Encode(response)
}

func updateMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	movieId := chi.URLParam(req, "id")
	var movie Movie
	err := json.NewDecoder(req.Body).Decode(&movie)
	if err != nil {
		http.Error(res, "error decoding movie data", http.StatusInternalServerError)
		return
	}
	movie.ID = strconv.Itoa(rand.Intn(10000000))
	for index, value := range movies{
		if value.ID == movieId {
			movies = append(movies[:index], movies[index+1:]...)
			movies = append(movies, movie)
			res.WriteHeader(http.StatusOK)
			response := map[string]interface{}{
				"message": "Movie updated successfully",
				"movie": movie,
			}
			json.NewEncoder(res).Encode(response)
			return
		}
	}
	// if movie id not found
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(map[string]string{"message": "Movie not found"})
}

func deleteMovie(res http.ResponseWriter, req *http.Request){
	res.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(req)
	// idToDelete := req.PathValue("id")
	idToDelete := chi.URLParam(req, "id")
	for index, item := range movies {
		if item.ID == idToDelete {
			movies = append(movies[:index], movies[index+1:]...)
			res.WriteHeader(http.StatusOK)
			response := map[string]string{"message": "Movie deleted successfully"}
			json.NewEncoder(res).Encode(response) // this will print a pretty looking json
			return
		}
	}
	
	// if loop finishes and we havent returned, means the movie didnot exists
	res.WriteHeader(http.StatusNotFound)
	json.NewEncoder(res).Encode(map[string]string{"error": "Movie not found"})
}

func main(){
	// r := mux.NewRouter()
	// r := http.NewServeMux()
	r := chi.NewRouter()
	r.Use(middleware.Logger) // Logs every request to the terminal
	r.Use(middleware.Recoverer) // prevent the server from crashing if there is any panic

	movies = append(movies, Movie{ID: "1", Isbn: "2345234", Title: "Movie One", Director: &Director{FirstName: "John", LastName: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "4545674", Title: "Movie Two", Director: &Director{FirstName: "Steve", LastName: "Smith"}})
	movies = append(movies, Movie{ID: "3", Isbn: "9876543", Title: "The Matrix", Director: &Director{FirstName: "Lana", LastName: "Wachowski"}})
	movies = append(movies, Movie{ID: "4", Isbn: "1122334", Title: "Interstellar", Director: &Director{FirstName: "Christopher", LastName: "Nolan"}})
	movies = append(movies, Movie{ID: "5", Isbn: "5566778", Title: "Pulp Fiction", Director: &Director{FirstName: "Quentin", LastName: "Tarantino"}})
	movies = append(movies, Movie{ID: "6", Isbn: "9988776", Title: "Jurassic Park", Director: &Director{FirstName: "Steven", LastName: "Spielberg"}})

/* using Gorilla Mux
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
*/
/* using inbuilt http.NewServe()
	r.HandleFunc("GET /movies", getMovies)
	r.HandleFunc("GET /movies/{id}", getMovie)
	r.HandleFunc("POST /movies", createMovie)
	r.HandleFunc("PUT /movies", updateMovie)
	r.HandleFunc("DELETE /movies", deleteMovie)
*/
	// go-chi/chi which is popular in Go community
	r.Get("/movies", getMovies)
	r.Get("/movies/{id}", getMovie)
	r.Post("/movies", createMovie)
	r.Put("/movies/{id}", updateMovie)
	r.Delete("/movies/{id}", deleteMovie)

	fmt.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}