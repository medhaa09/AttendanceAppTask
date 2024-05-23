package Auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type SignedDetails struct {
	Handle string
	Role   string
	jwt.StandardClaims
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("attdapp").Collection("data")
	return collection
}

func DBinstance() *mongo.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MongoDb := os.Getenv("MONGODB_URL")
	client, err := mongo.NewClient(options.Client().ApplyURI(MongoDb))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")

	return client
}

var Client *mongo.Client = DBinstance()
var userCollection *mongo.Collection = OpenCollection(Client, "user")
var SECRET_KEY string = os.Getenv("SECRET_KEY")

// generating both the detailed token and refresh token
func GenerateAllTokens(Handle string, Role string) (signedToken string, signedRefreshToken string, err error) {
	fmt.Println("generate all token function: ", Handle, Role)
	claims := &SignedDetails{
		Handle: Handle,
		Role:   Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}
	return token, refreshToken, err
}

func ValidateToken(signedToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil {
		return nil, err
	}
	return token, nil
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		fmt.Println("inTokenAuthMiddleware function ", authHeader)
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		} else {
			temp := strings.Split(authHeader, "Bearer ")
			if len(temp) < 2 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
				return
			}
			tokenString := temp[1]
			token, err := ValidateToken(tokenString)
			if err != nil || !token.Valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			} else {

				fmt.Println("token validated")
			}
			///extracting claims from token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				return
			}
			//extracting role from claims
			fmt.Println("claims  ", claims)
			role, ok := claims["Role"].(string)
			fmt.Println("role :=", role)
			if !ok {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found in token"})
				return
			}

			c.Set("role", role)
		}
		c.Next()
	}
}

// Middleware function to check if user is allowed to access register route
func IsAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		fmt.Println("role: ", role)
		if !exists || role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
			fmt.Println("error is in isadmin function")
			c.Abort()
			return
		} else {
			fmt.Println("admin access granted")
		}

		c.Next()
	}
}
