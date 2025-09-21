package imageUpload

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
  "coldwheels/lib"
  "strings"
  "net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type S3Config struct {
    Key      string
    Secret   string
    Endpoint string
    Region   string
    Bucket   string
    CDNUrl   string
}

type S3Uploader struct {
    client *s3.Client
    config S3Config
}

type Presigner struct {
	PresignClient *s3.PresignClient
}


func PresignUrls(ctx *gin.Context, cars []lib.Car) ([]lib.Car, error) {
  // takes a list of cars, signs each url
  for i := range cars {
    if len(cars[i].ImageURLs) > 0 {
      // sign all the urls in this car
      for j := range cars[i].ImageURLs {

        // get object key from object url
        u, err := url.Parse(cars[i].ImageURLs[j])
        if err != nil {
          return nil, err
        }

        objectKey := strings.TrimPrefix(u.Path, "/")

	// DO THIS extra trimming ONLY IF youre using some s3 endpoint like https://domain/somethingextra
	pathParts := strings.Split(objectKey, "/")
	remainingParts := pathParts[1:]
	objectKey = strings.Join(remainingParts, "/")

        fmt.Println("object key from url : ", objectKey)

        client, err := NewS3Client()
        if err != nil {
          return nil, err
        }

        presignClient := s3.NewPresignClient(client)
        presigner := Presigner{PresignClient: presignClient}
        s3Bucket := os.Getenv("S3_BUCKET")
        signedRequest, err := presigner.GetObject(ctx, s3Bucket, objectKey)
        if err != nil {
          return nil, err
        }


        // in place modification
        cars[i].ImageURLs[j] = signedRequest.URL

      }
    }
  }
  return cars, nil
}



func NewS3Client() (*s3.Client, error) {

  cfg, err := config.LoadDefaultConfig(context.TODO())

  // fmt.Println("config in new client maker function is: ", cfg)
  if err != nil {
    fmt.Println(err)
    return nil, errors.New("cant make config using stuff")
  }
  // fmt.Println("cfg loading no error: ", cfg)

  s3Region := os.Getenv("S3_REGION")
  s3Endpoint := os.Getenv("S3_ENDPOINT")
  client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(s3Endpoint)
    o.Region = *aws.String(s3Region)
  })
  // fmt.Println("client: ", client)

  return client, nil

}

func NewS3Uploader() (*manager.Uploader, error) {

  cfg, err := config.LoadDefaultConfig(context.TODO())

  // fmt.Println("config in new uploader maker function is: ", cfg)
  if err != nil {
    fmt.Println(err)
    return nil, errors.New("cant make config using stuff")
  }
  // fmt.Println("cfg loading no error: ", cfg)

  s3Region := os.Getenv("S3_REGION")
  s3Endpoint := os.Getenv("S3_ENDPOINT")
  client := s3.NewFromConfig(cfg, func(o *s3.Options) {
    o.BaseEndpoint = aws.String(s3Endpoint)
    o.Region = *aws.String(s3Region)
  })
  // fmt.Println("client: ", client)

  uploader := manager.NewUploader(client)

  return uploader, nil

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
        // c.JSON(500, gin.H{"error":"Error loading .env file"})
  	fmt.Println("Error loading .env file while in upload images handler")
    }

    s3Bucket := os.Getenv("S3_BUCKET")

    form, err := c.MultipartForm()
    if err != nil {
        c.JSON(400, gin.H{"error": "Failed to parse multipart form"})
        return
    }
    // fmt.Println("form is: ", form);

    files := form.File["images"]
    if len(files) == 0 {
        c.JSON(400, gin.H{"error": "No images provided"})
        return
    }
    
    // Validate image URLs
    if len(files) > lib.ImagesPerCarLimit {
      c.JSON(http.StatusBadRequest, gin.H{"error": "Maximum images limit crossed"})
      return
    }

    // fmt.Println("files is: ", files);

    // Initialize uploader (you might want to do this once at startup)
    uploader, err := NewS3Uploader()
    // presignClient := s3.NewPresignClient(client)
    // presigner := Presigner{PresignClient: presignClient}

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

      fmt.Println("file opened")
      // fmt.Println(fileO)


      uid := c.MustGet("UserID").(string)
      result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(s3Bucket),
        Key:    aws.String(uid + "/" + file.Filename),
        Body:   fileO,
      })
      fmt.Println("result is: ", result)

      // signedRequest, err := presigner.GetObject(c, s3Bucket, *result.Key)

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
