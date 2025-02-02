package imageUpload

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	// "time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
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

type Presigner struct {
	PresignClient *s3.PresignClient
}

func NewSpacesUploader() (*manager.Uploader, *s3.Client, error) {

  cfg, err := config.LoadDefaultConfig(context.TODO())

  fmt.Println("config in new uploader maker function is: ", cfg)
  if err != nil {
    fmt.Println(err)
    return nil, nil, errors.New("cant make config using stuff")
  }
  fmt.Println("cfg loading no error: ", cfg)

  //
  // Go TODO
  //

  spacesRegion := os.Getenv("SPACES_REGION")
  spacesEndpoint := os.Getenv("SPACES_ENDPOINT")
  client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(spacesEndpoint)
    o.Region = *aws.String(spacesRegion)
  })
  fmt.Println("client: ", client)

  uploader := manager.NewUploader(client)

  return uploader, client, nil

}


func (presigner Presigner) GetObject(
	ctx context.Context, bucketName string, objectKey string) (*v4.PresignedHTTPRequest, error) {
	request, err := presigner.PresignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
    // opts.Expires = time.Duration(lifetimeSecs * int64(time.Second))
	})
	if err != nil {
		fmt.Println("Couldn't get a presigned request to get %v:%v. Here's why: %v\n",
			bucketName, objectKey, err)
	}
	return request, err
}


func UploadImagesHandler(c *gin.Context) {
    // Get uploaded files

    err := godotenv.Load()
    if err != nil {

        c.JSON(500, gin.H{"error":"Error loading .env file"})
        return
    }

    // spacesKey := os.Getenv("SPACES_KEY")
    // spacesSecret := os.Getenv("SPACES_SECRET")
    spacesBucket := os.Getenv("SPACES_BUCKET")
    // spacesBucketEndpoint := os.Getenv("SPACES_BUCKET_ENDPOINT")

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
    uploader, client, err := NewSpacesUploader()
    presignClient := s3.NewPresignClient(client)
    presigner := Presigner{PresignClient: presignClient}

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
      defer fileO.Close()

      fmt.Println("file opened: ")
      fmt.Println(fileO)
      // fmt.Println("spaces key is: ")
      // fmt.Println(spacesKey)


      // TODO: the Key is being used as a file name to save instead of actually using as credential

      uid := c.MustGet("UserID").(string)
      result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(spacesBucket),
        Key:    aws.String("coldwheels/" + uid + "/" + file.Filename),
        Body:   fileO,
      })
      fmt.Println("result is: ", result)

      signedRequest, err := presigner.GetObject(c, spacesBucket, *result.Key)

      if err != nil {
          c.JSON(500, gin.H{"error": fmt.Sprintf("Failed to upload images: %v", err)})
          return
      }

      // results = append(results, result.Location)
      results = append(results, signedRequest.URL)

    }

    if len(results) == len(files) {
      c.JSON(200, gin.H{"urls": results})
      return
    }
}
