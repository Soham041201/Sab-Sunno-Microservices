package gemini

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetupGeminiServiceToGenerateResponsee() (*GeminiClient, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		log.Fatal("GOOGLE_API_KEY environment variable not set")
	}
	client, err := NewGeminiClient(apiKey)
	
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendSetup()

	if err != nil {
		log.Fatal("setup:", err)
	}

	return client, nil

}
