package main

import (
	"attdapp/Auth"
	"attdapp/Models"
	"attdapp/Store"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var faceAPIKey string
var faceAPISecret string

func main() {
	mongoStore := &Store.MongoStore{}
	mongoStore.OpenConnectionWithMongoDB()
	faceAPIKey = "Oy-1LgVH6ZkET9FytDhQAmVW2dzUMf0E"
	faceAPISecret = "jrNQCmO8v1ULJNiaGTL8QCC6Xys3gmBc"

	router := gin.Default()

	// Configure CORS to allow requests from the frontend
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/register", Auth.TokenAuthMiddleware(), Auth.IsAdmin(), func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
			return
		}

		// Extracting form values
		name := c.Request.FormValue("name")
		handle := c.Request.FormValue("handle")
		pass := c.Request.FormValue("pass")

		fmt.Println("Received values - Name:", name, "Handle:", handle, "Password:", pass)

		// Reading the image file
		file, header, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read image file"})
			return
		}
		defer file.Close()

		imageData, err := ioutil.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image data"})
			return
		}

		newUser := Models.User{
			Name:     name,
			Handle:   handle,
			Password: pass,
			Role:     "student",
			Image: Models.Image{
				Data: imageData,
			},
		}

		// Storing user data and image in the database
		err = mongoStore.StoreUserData(newUser, imageData, header.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
	})

	router.POST("/login", func(c *gin.Context) {
		var loginCredentials Models.User
		err := c.ShouldBindJSON(&loginCredentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		isAuthenticated, role := mongoStore.UserLogin(loginCredentials.Handle, loginCredentials.Password)
		if isAuthenticated {
			fmt.Println("user is authenticated")
			signedToken, signedRefreshToken, err := Auth.GenerateAllTokens(loginCredentials.Handle, role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"token": signedToken, "refreshToken": signedRefreshToken, "role": role})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		}
	})

	router.POST("/mark-attendance", markAttendance(mongoStore))

	port := ":8080"
	fmt.Printf("Server running on port %s\n", port)
	if err := router.Run(port); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}

func markAttendance(mongoStore *Store.MongoStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req Models.RecognitionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		base64Data := req.Image
		if base64Data == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Empty image data"})
			return
		}

		// Log the received base64 image data
		fmt.Println("Received Base64 Image Data:", base64Data[:30], "...")

		imageData, err := base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image data"})
			return
		}
		fmt.Println(imageData)

		users := mongoStore.GetAllUsers()
		if users == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "No users found"})
			return
		}

		for _, user := range users {
			storedImageData := user.Image.Data
			storedImageBase64 := base64.StdEncoding.EncodeToString(storedImageData)

			if storedImageBase64 == "" {
				fmt.Println("Stored image data is empty, skipping user.")
				continue
			}

			payload := strings.NewReader(fmt.Sprintf(
				"api_key=%s&api_secret=%s&image_base64_1=%s&image_base64_2=%s",
				faceAPIKey, faceAPISecret, base64Data, storedImageBase64))

			req, err := http.NewRequest("POST", "https://api-us.faceplusplus.com/facepp/v3/compare", payload)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
				return
			}
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

			var res *http.Response
			var body []byte
			for retries := 0; retries < 3; retries++ {
				res, err = http.DefaultClient.Do(req)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform request"})
					return
				}
				defer res.Body.Close()

				body, err = ioutil.ReadAll(res.Body)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
					return
				}

				fmt.Println("Raw response body:", string(body))

				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err != nil {
					fmt.Println("Failed to parse response:", err)
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
					return
				}

				if result["error_message"] == "CONCURRENCY_LIMIT_EXCEEDED" {
					time.Sleep(time.Duration(retries+1) * time.Second)
					continue
				}

				fmt.Println("Face++ API Response:", result)

				if confidence, ok := result["confidence"].(float64); ok {
					fmt.Printf("Comparing face. Confidence: %f\n", confidence)
					if confidence > 80.0 {
						c.JSON(http.StatusOK, gin.H{"message": "Attendance marked as present", "user": user})
						return
					}
				} else {
					fmt.Println("Confidence score not found in API response.")
				}
				break
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Attendance marked as absent"})
	}
}
