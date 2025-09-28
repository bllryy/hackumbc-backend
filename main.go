package main

import "fmt"
import "google.golang.org/genai"
import "context"
import "github.com/gin-gonic/gin"
import "encoding/json"
import "net/http"
import "strings"
import "github.com/gin-contrib/cors"
import "github.com/joho/godotenv"

type errors struct {
	Line int    `json:"line"`
	Col  int    `json:"column"`
	Desc string `json:"message"`
}

type response struct {
	Errors []errors `json:"errors"`
}

type request struct {
	Data string `json:"file"`
}

func checkFile(c *gin.Context) {
	var r request
	err := c.BindJSON(&r)
	if err != nil {
		return
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  "AIzaSyB7MioCq_eDxEfmUSraYHxRBJBEXaqZ3pc",
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return
	}
	// prompt. add support for files
	parts := []*genai.Part{
		{Text: "Analyze this code for errors. Respond in JSON. If there are no errors return empty JSON. Follow the example: \n {\"errors\": [{\"line\": LINE},{\"column\": COL},{\"message\": MESSAGE},{\"severity\": SEVERITY}]}. Include a reason for the error."},
		{InlineData: &genai.Blob{Data: ([]byte)(r.Data), MIMEType: "text/plain"}},
	}
	// call
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-pro", []*genai.Content{{Parts: parts}}, nil)
	if err != nil {
		return
	}
	// output. replace with json stuff for frontend
	if result.Candidates != nil {
		for _, v := range result.Candidates {
			for _, k := range v.Content.Parts {
				// remove formatting
				f := strings.Replace(k.Text, "```json", "", 1)
				s := strings.Replace(f, "```", "", 1)

				var data response
				err := json.Unmarshal([]byte(s), &data)
				if err != nil {
					return
				}

				c.PureJSON(http.StatusOK, data)
				return
			}
		}
	}

	return
}

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/api/analyze", checkFile)

	router.Run("0.0.0.0:8888")
}
