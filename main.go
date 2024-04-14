package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mi-gongan/commention_backend/pkg/route"
	"github.com/supabase-community/supabase-go"
)

var (
	SupabaseClient *supabase.Client
	SecretKey      []byte
)

type User struct {
	Email string `json:"email"`
}

type Claims struct {
	User `json:"user"`
	jwt.StandardClaims
}

func initSupbase() {
	SUPABASE_API_URL := os.Getenv("SUPABASE_API_URL")
	SUPABASE_API_KEY := os.Getenv("SUPABASE_API_KEY")

	client, err := supabase.NewClient(SUPABASE_API_URL, SUPABASE_API_KEY, nil)
	if err != nil {
		log.Fatalf("failed to create Supabase client: %v", err)
	}

	SupabaseClient = client
}

func init() {
	godotenv.Load(".env")

	SecretKey = []byte(os.Getenv("JWT_SECRET"))

	initSupbase()
}

func main() {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	route.RegisterAuthRoutes(r, "/auth")

	r.Run(":8080")
}
