package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/repository/sqlconnect"
	"time"

	"github.com/joho/godotenv"
)

type user struct {
	Name string `json:"name"`
	Age string `json:"age"`
	City string `json:"city"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	_, err = sqlconnect.ConnectDb()
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	} 
	port := os.Getenv("API_PORT")

	cert := "cert.pem"
	key := "key.pem"

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the API"))
		fmt.Println("Hello from the API")
	})

	mux.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Hello GET method on teachers route"))
			fmt.Println("Hello GET method on teachers route")
		case http.MethodPost:
			// Parse form data (necessary for x-www-form-urlencoded)
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}
			fmt.Println("Form data: ", r.Form)

			// Prepare response data
			response := make(map[string]interface{})
			for key, value := range r.Form {
				response[key] =value[0]
			}
			fmt.Println("Processed response map: ", response)

			// Raw body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()
			fmt.Println("Raw body: ", body)
			fmt.Println("Raw body: ", string(body))

			var userInstance user
			err = json.Unmarshal(body, &userInstance)
			if err != nil {
				return
			}
			fmt.Println("Unmarshaled JSON into an instance user struct", userInstance)
			fmt.Println("Received username as:", userInstance.Name)

			response1 := make(map[string]interface{})
			for key, value := range r.Form {
				response[key] =value[0]
			}

			err = json.Unmarshal(body, &response1)
			if err != nil {
				return
			}
			fmt.Println("Unmarshaled JSON into a map", response1)


			w.Write([]byte("Hello POST method on teachers route"))
			fmt.Println("Hello POST method on teachers route")
		case http.MethodPut:
			w.Write([]byte("Hello PUT method on teachers route"))
			fmt.Println("Hello PUT method on teachers route")
		case http.MethodPatch:
			w.Write([]byte("Hello PATCH method on teachers route"))
			fmt.Println("Hello PATCH method on teachers route")
		case http.MethodDelete:
			w.Write([]byte("Hello DELETE method on teachers route"))
			fmt.Println("Hello DELETE method on teachers route")
		default:
			w.Write([]byte("Hello from the teachers API"))
			fmt.Println("Hello from the teachers API")
		}
	})

	mux.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the students API"))
		fmt.Println("Hello from the students API")
	})

	mux.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the execs API"))
		fmt.Println("Hello from the execs API")
	})

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery: true,
		CheckBody:  true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist: []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	// secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)
	// secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ResponseTimeMiddleware, mw.Cors)

	server := &http.Server{
		Addr: port,
		// Handler: mux,
		Handler: mw.Hpp(hppOptions)(rl.Middleware(mw.Compression(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux)))))),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}

	// err := http.ListenAndServe(port, nil)
}