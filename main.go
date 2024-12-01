package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Movie struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Director          *Director `json:"director"`
	ProductionCompany string    `json:"production_company"`
	Language          string    `json:"language"`
	Rating            float64   `json:"rating"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	if movie.Title == "" || movie.Director == nil || movie.ProductionCompany == "" || movie.Language == "" || movie.Rating == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields"})
		return
	}

	movie.ID = uuid.New().String()

	movies = append(movies, movie)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)

	log.Printf("Created a new movie: %s", movie.Title)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			var updatedMovie Movie
			if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
				return
			}

			if updatedMovie.Title == "" || updatedMovie.Director == nil || updatedMovie.ProductionCompany == "" || updatedMovie.Language == "" || updatedMovie.Rating == 0 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "Missing required fields"})
				return
			}

			updatedMovie.ID = params["id"]
			movies[index] = updatedMovie

			json.NewEncoder(w).Encode(updatedMovie)
			log.Printf("Updated movie: %s", updatedMovie.Title)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "Movie not found"})
}

func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{ID: "1", Title: "Movie One", Director: &Director{Firstname: "Chris", Lastname: "Smith"}, ProductionCompany: "ABC Studios", Language: "English", Rating: 7.5})
	movies = append(movies, Movie{ID: "2", Title: "Movie Two", Director: &Director{Firstname: "Dan", Lastname: "Smith"}, ProductionCompany: "XYZ Films", Language: "Spanish", Rating: 6.2})
	movies = append(movies, Movie{ID: "3", Title: "Movie Three", Director: &Director{Firstname: "Steve", Lastname: "Smith"}, ProductionCompany: "Global Entertainment", Language: "French", Rating: 8.0})

	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")

	fmt.Printf("Starting server at port 8090\n")
	log.Fatal(http.ListenAndServe(":8090", r))
}
