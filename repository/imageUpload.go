package imageUpload

import (
	"context"
	"errors"
	"fmt"
	"os"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type SpacesConfig struct {
    Key      string
    Secret   string
    Endpoint string
    Region   string
    Bucket   string
    CDNUrl   string
}

type SpacesUploader struct {
    client *s3.Client
    config SpacesConfig
}

func NewSpacesUploader() (*manager.Uploader, error) {

  cfg, err := config.LoadDefaultConfig(context.TODO())

  fmt.Println("config in new uploader maker function is: ", cfg)
  if err != nil {
    return nil, errors.New("cant make config using stuff")
  }

  //
  // Go TODO
  //
  client := s3.NewFromConfig(cfg)
  fmt.Println("client: ", client)

  uploader := manager.NewUploader(client)
  return uploader, nil

}


func UploadImagesHandler(c *gin.Context) {
    // Get uploaded files

    err := godotenv.Load()
    if err != nil {

        c.JSON(500, gin.H{"error":"Error loading .env file"})
        return
    }

    spacesKey := os.Getenv("SPACES_KEY")
    // spacesSecret := os.Getenv("SPACES_SECRET")
    spacesBucket := os.Getenv("SPACES_BUCKET")
    // spacesEndpoint := os.Getenv("SPACES_ENDPOINT")
    // spacesBucketEndpoint := os.Getenv("SPACES_BUCKET_ENDPOINT")
    // spacesRegion := os.Getenv("SPACES_REGION")

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "Failed to parse multipart form"})
        return
    }
    fmt.Println("form is: ", form);

    files := form.File["images"]
    if len(files) == 0 {
        c.JSON(400, gin.H{"error": "No images provided"})
        return
    }
    fmt.Println("files is: ", files);

    // Initialize uploader (you might want to do this once at startup)
    uploader, err := NewSpacesUploader()

    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to initialize uploader"})
        return
    }


    var results []string

    for _, file := range files {

      fileO, err := file.Open()
      if err != nil {
          c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
          return
      }

      fmt.Println("file opened: ")
      fmt.Println(fileO)
      fmt.Println("spaces key is: ")
      fmt.Println(spacesKey)

      result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(spacesBucket),
        Key:    aws.String(spacesKey),
        Body:   fileO,
      })
      fmt.Println("result is: ", result)

      if err != nil {
          c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to upload images: %v", err)})
          return
      }

      results = append(results, result.Location)

    }

    if len(results) == len(files) {
      c.JSON(200, gin.H{"urls": results})
      return
    }
}
