package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/repository/sqlconnect"
	"time"

	"github.com/joho/godotenv"
)

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

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	// secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ResponseTimeMiddleware, rl.Middleware, mw.Cors)
	// secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ResponseTimeMiddleware, mw.Cors)

	server := &http.Server{
		Addr: port,
		// Handler: mux,
		Handler:   mw.Hpp(hppOptions)(rl.Middleware(mw.Compression(mw.ResponseTimeMiddleware(mw.SecurityHeaders(mw.Cors(mux)))))),
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server is running on port: ", port)

	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}

	// err := http.ListenAndServe(port, nil)
}
