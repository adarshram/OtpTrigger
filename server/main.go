package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"url-trigger/database"
)

func Logger(s interface{}) {
	log.Printf("%v\n", s)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World!")
	})
	http.HandleFunc("/otps", func(w http.ResponseWriter, r *http.Request) {
		authenticated, err := authenticateRequest(w, r, Logger)
		if err != nil {
			Logger(err)
			return
		}
		if authenticated {
			otpDisplayHandler(w, r, Logger)
		}
	})
	http.HandleFunc("/notifyOtp", func(w http.ResponseWriter, r *http.Request) {
		otpHandler(w, r, Logger)
	})
	log.Println("Starting server1 on :8066")
	if err := http.ListenAndServe(":8066", nil); err != nil {
		log.Fatal(err)
	}

}

func otpHandler(w http.ResponseWriter, r *http.Request, logFn func(interface{})) {
	enableCors(&w)
	ua := r.Header.Get("Authorization")
	logFn(ua)

	if r.Body == nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("No request body")
		fmt.Fprintln(w, "No request body")
		return
	}

	var t struct {
		Message string `json:"message"`
	}
	db := database.NewDataBase()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading request body")
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("No request body")
		fmt.Fprintln(w, "No request body")
		return
	}
	err = json.Unmarshal(body, &t)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("Error while unmarshalling request body")
		fmt.Fprintln(w, "No request body")
		return

	}

	message := t.Message
	if message != "" {
		err = db.InsertOtpToTable("otps", message)
		if err != nil {
			log.Println("Error while inserting data to table")
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintln(w, "Not Successful")
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, fmt.Sprintf("Successfully inserted data: %s", message))
	}

}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization,Origin")
}

func otpDisplayHandler(w http.ResponseWriter, r *http.Request, logFn func(interface{})) {
	db := database.NewDataBase()
	otps, err := db.RetrieveLatest("otps")
	if err != nil && err.Error() != "no data found" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		log.Println("Error while retrieving data from table")
		fmt.Fprintf(w, "error while retrieving data from table: %v", err)
		return
	}

	if len(otps) == 0 {
		log.Println("No data found")
	}
	outputData := struct {
		Otps    []string `json:"otps"`
		Success bool     `json:"success"`
	}{
		Otps:    otps,
		Success: true,
	}

	jsonData, err := json.Marshal(outputData)
	if err != nil {
		panic(err)
	}
	jsonString := string(jsonData)
	fmt.Fprintln(w, jsonString)
}

func authenticateRequest(w http.ResponseWriter, r *http.Request, logFn func(interface{})) (bool, error) {
	enableCors(&w)
	method := r.Method
	if method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return false, nil
	}

	db := database.NewDataBase()
	authHeader := r.Header.Get("Authorization")

	_, _, err := db.AuthenticateBearer(authHeader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return false, err
	}
	return true, nil
}
