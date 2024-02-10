package main

import (
	// "context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/bson"
	_ "github.com/go-sql-driver/mysql"
)

type ApiResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// Declare global variables
var (
	mysqlVarsProvided     bool
	userProvidedMysqlVars *map[string]string
	requiredVars          = []string{"MYSQL_HOST", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_PORT", "MYSQL_DATABASE"}

	mongoVarsProvided     bool
	userProvidedMongoVars *map[string]string
)

// getEnvOrDefault retrieves an environment variable or returns a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Function to check MySQL connection
func checkMySQLConnection() error {
	// Attempt to connect to MySQL database
	db, err := sql.Open("mysql", getDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	// Ping database to check connection status
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}

// Function to get MySQL connection details
func getDSN() string {
	var (
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase string
	)

	if mysqlVarsProvided && userProvidedMysqlVars != nil {
		// If mysqlVarsProvided is true and userProvidedMysqlVars is not nil, use values from userProvidedMysqlVars
		mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase = (*userProvidedMysqlVars)["MYSQL_USER"], (*userProvidedMysqlVars)["MYSQL_PASSWORD"],
			(*userProvidedMysqlVars)["MYSQL_HOST"], (*userProvidedMysqlVars)["MYSQL_PORT"],
			(*userProvidedMysqlVars)["MYSQL_DATABASE"]
	} else {
		// If mysqlVarsProvided is false or userProvidedMysqlVars is nil, use values from environment variables with defaults
		mysqlUser = getEnvOrDefault("MYSQL_USER", "root")
		mysqlPassword = getEnvOrDefault("MYSQL_PASSWORD", "")
		mysqlHost = getEnvOrDefault("MYSQL_HOST", "localhost")
		mysqlPort = getEnvOrDefault("MYSQL_PORT", "3306")
		mysqlDatabase = getEnvOrDefault("MYSQL_DATABASE", "mysql")
	}

	// Construct MySQL data source name (DSN)
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
}

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		// Print request method and URL
		fmt.Printf("Received request: %s %s\n", c.Request.Method, c.Request.URL.Path)

		// Print request headers
		fmt.Println("Headers:")
		for name, values := range c.Request.Header {
			for _, value := range values {
				fmt.Printf("%s: %s\n", name, value)
			}
		}

		// Continue processing the request
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		// Convert headers to map[string]string
		headers := make(map[string]string)
		for name, values := range c.Request.Header {
			headers[name] = values[0]
		}

		// Set Content-Type header
		c.Header("Content-Type", "application/json")

		// Return headers in JSON format
		c.JSON(http.StatusOK, headers)
	})

	// Get headers
	r.POST("/", func(c *gin.Context) {
		// Convert headers to map[string]interface{}
		headers := make(map[string]interface{})
		for name, values := range c.Request.Header {
			headers[name] = values[0]
		}

		// Create response object
		response := ApiResponse{
			Success: true,
			Message: "Headers fetched successfully",
			Data:    headers,
		}

		// Return response in JSON format
		c.JSON(http.StatusOK, response)
	})

	// Get environment variables
	r.POST("/env", func(c *gin.Context) {
		// Retrieve all environment variables
		envMap := make(map[string]interface{})

		// Attempt to retrieve environment variables
		envPairs := os.Environ()
		if envPairs == nil {
			// If unable to retrieve environment variables, create a failed response
			response := ApiResponse{
				Success: false,
				Message: "Failed to fetch environment variables",
				Data:    nil,
			}
			// Return response in JSON format
			c.JSON(http.StatusInternalServerError, response)
			return
		}

		for _, pair := range envPairs {
			keyVal := strings.Split(pair, "=")
			envMap[keyVal[0]] = keyVal[1] // value is converted to interface{}
		}

		// Create response object
		response := ApiResponse{
			Success: true,
			Message: "Environment variables fetched successfully",
			Data:    envMap,
		}

		// Return response in JSON format
		c.JSON(http.StatusOK, response)
	})

	// Get environment variables except os variables
	r.POST("/env-from-dotenv", func(c *gin.Context) {
		// Retrieve environment variables specifically from .env file
		envMap := make(map[string]interface{})

		// Load environment variables from .env file
		envFile, err := godotenv.Read(".env")
		if err != nil {
			// Handle error if .env file cannot be read
			c.JSON(http.StatusInternalServerError, ApiResponse{
				Success: false,
				Message: "Error reading .env file: " + err.Error(),
				Data:    nil,
			})
			return
		}

		for key, value := range envFile {
			envMap[key] = value
		}

		// Create response object
		response := ApiResponse{
			Success: true,
			Message: "Environment variables from .env file fetched successfully",
			Data:    envMap,
		}

		// Return response in JSON format
		c.JSON(http.StatusOK, response)
	})

	// search environment variables with a searchKey
	r.GET("/env/:searchKey", func(c *gin.Context) {
		searchKey := strings.ToLower(c.Param("searchKey"))

		// Check if searchKey has at least two characters
		if len(searchKey) < 2 {
			// Create response object for invalid searchKey
			response := ApiResponse{
				Success: false,
				Message: "At least 2 characters are required to make a search",
				Data:    nil,
			}
			// Return response in JSON format
			c.JSON(http.StatusBadRequest, response)
			return
		}

		envMap := make(map[string]interface{})

		for _, pair := range os.Environ() {
			keyVal := strings.Split(pair, "=")
			key := strings.ToLower(keyVal[0])
			if strings.Contains(key, searchKey) {
				envMap[keyVal[0]] = keyVal[1]
			}
		}

		// Create response object
		response := ApiResponse{
			Success: true,
			Message: "Environment variables with key similar to search key fetched successfully",
			Data:    envMap,
		}

		// Return response in JSON format
		c.JSON(http.StatusOK, response)
	})

	// Set mysql vars & check connecction
	r.POST("/setup-and-check-mysql-connection", func(c *gin.Context) {
		// Parse request body to get provided MYSQLDB variables
		var reqBody map[string]string
		if err := c.ShouldBindJSON(&reqBody); err != nil {
			c.JSON(http.StatusBadRequest, ApiResponse{
				Success: false,
				Message: "Failed to parse request body: " + err.Error(),
			})
			return
		}

		// Validate presence of all required keys
		missingKeys := make([]string, 0)
		for _, key := range requiredVars {
			if _, ok := reqBody[key]; !ok {
				missingKeys = append(missingKeys, key)
			}
		}

		if len(missingKeys) > 0 {
			c.JSON(http.StatusBadRequest, ApiResponse{
				Success: false,
				Message: fmt.Sprintf("Required MYSQLDB variables missing in request: %s", strings.Join(missingKeys, ", ")),
			})
			return
		}

		// Check if all provided variables are empty
		allEmpty := true
		for _, key := range requiredVars {
			if val := reqBody[key]; val != "" {
				allEmpty = false
				break
			}
		}

		// If all provided variables are empty, retrieve credentials from .env file or defaults
		if allEmpty {
			// Retrieve MYSQLDB connection details from .env file or defaults
			reqBody["MYSQL_HOST"] = getEnvOrDefault("MYSQL_HOST", "localhost")
			reqBody["MYSQL_USER"] = getEnvOrDefault("MYSQL_USER", "root")
			reqBody["MYSQL_PASSWORD"] = getEnvOrDefault("MYSQL_PASSWORD", "")
			reqBody["MYSQL_PORT"] = getEnvOrDefault("MYSQL_PORT", "3306")
			reqBody["MYSQL_DATABASE"] = getEnvOrDefault("MYSQL_DATABASE", "myDatabase")

			userProvidedMysqlVars = nil
			mysqlVarsProvided = false
		} else {

			// Validate that non-empty variables are not empty
			emptyKeys := make([]string, 0)
			for _, key := range []string{"MYSQL_HOST", "MYSQL_USER", "MYSQL_PORT", "MYSQL_DATABASE"} {
				if reqBody[key] == "" {
					emptyKeys = append(emptyKeys, key)
				}
			}
			if len(emptyKeys) > 0 {
				c.JSON(http.StatusBadRequest, ApiResponse{
					Success: false,
					Message: fmt.Sprintf("%s cannot be empty", strings.Join(emptyKeys, ", ")),
				})
				return
			}

			// Set global variables to indicate MYSQLDB vars provided by the user
			mysqlVarsProvided = true
			userProvidedMysqlVars = &reqBody
		}

		// Attempt to connect to MySQL database and check the connection status
		if err := checkMySQLConnection(); err != nil {
			c.JSON(http.StatusInternalServerError, ApiResponse{
				Success: false,
				Message: "Failed to connect to MySQL database: " + err.Error(),
			})
			return
		}

		// Connection successful, return success response with information about credentials
		var credentialSource string
		if mysqlVarsProvided {
			credentialSource = "userProvided"
		} else {
			credentialSource = "environment"
		}

		// Include information about the source of credentials in the response message
		responseMessage := fmt.Sprintf("MySQL connection successful. Credentials taken from: %s", credentialSource)

		c.JSON(http.StatusOK, ApiResponse{
			Success: true,
			Message: responseMessage,
		})
	})

	// Check timeout
	r.GET("/timeout/:timeoutValue", func(c *gin.Context) {
		timeoutValueStr := c.Param("timeoutValue")
		timeoutValue, err := strconv.Atoi(timeoutValueStr)
		if err != nil || timeoutValue <= 0 {
			timeoutValue = 30 // Default timeout value of 30 seconds
		}

		// Create a channel to receive notification
		done := make(chan bool)

		// Start a goroutine to simulate processing
		go func() {
			// Simulate processing time
			time.Sleep(time.Duration(timeoutValue) * time.Second)

			// Send notification through the channel
			done <- true
		}()

		// Wait for either timeout or processing completion
		select {
		case <-done:
			// Respond after processing completes
			response := ApiResponse{
				Success: true,
				Message: "Response after timeout",
				Data:    nil,
			}
			c.JSON(http.StatusOK, response)
		case <-time.After(300 * time.Second): // Adjust this timeout as needed
			// Send a response if the request takes too long to process
			response := ApiResponse{
				Success: false,
				Message: "Processing taking longer than expected",
				Data:    nil,
			}
			c.JSON(http.StatusRequestTimeout, response)
		}
	})

	// Setting application port
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}

	// Run the server
	log.Fatal(r.Run(":" + port))
}
