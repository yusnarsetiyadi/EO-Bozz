package thirdparty

import (
	cfg "capstone-alta1/config"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

// CREATE RANDOM STRING

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func autoGenerate(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return autoGenerate(length, charset)
}

func Upload(c echo.Context, imageField string, folderPath string) (string, error) {
	var errStr string

	fmt.Println("\n\n image upload ", imageField, " ---  ", folderPath)
	// upload foto
	file, _ := c.FormFile(imageField)
	if file != nil {
		file, fileheader, errFileHeader := c.Request().FormFile(imageField)
		if errFileHeader != nil {
			errStr = "Upload file failed. Field < " + imageField + " >, failed get File Header. Sent default image url. Error detail : " + errFileHeader.Error()
			fmt.Println(errStr)
			return "", errors.New("Upload file failed. Please try again later.")
		}

		randomStr := String(20)

		godotenv.Load(".env")

		s3Config := &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY_IAM"), os.Getenv("SECRET_KEY_IAM"), ""),
		}
		s3Session := session.New(s3Config)

		uploader := s3manager.NewUploader(s3Session)

		input := &s3manager.UploadInput{
			Bucket:      aws.String(os.Getenv("AWS_BUCKET_NAME")),                             // bucket's name
			Key:         aws.String(folderPath + "/" + randomStr + "-" + fileheader.Filename), // files destination location
			Body:        file,                                                                 // content of the file
			ContentType: aws.String("image/jpg"),                                              // content type
		}
		res, errUploadWithContext := uploader.UploadWithContext(context.Background(), input)
		if errUploadWithContext != nil {
			errStr = "Upload file failed. Field < " + imageField + " >, failed upload to server. Sent default image url. Error detail : " + errUploadWithContext.Error()
			fmt.Println(errStr)
			return "", errors.New("Upload file failed. Please try again later.")
		}
		// RETURN URL LOCATION IN AWS
		return res.Location, nil

	} else {
		fmt.Println("\n\n Upload file failed. Field < " + imageField + " > not found. Sent defafult image url.")
		return cfg.DEFAULT_IMAGE_URL, nil
	}
}

func UploadForUpdate(c echo.Context, imageField string, folderPath string) (string, error) {
	var errStr string
	// upload foto
	file, fileheader, errFileHeader := c.Request().FormFile(imageField)
	if errFileHeader != nil {
		errStr = "Upload file failed. Field < " + imageField + " >, failed get File Header. Sent default image url. Error detail : " + errFileHeader.Error()
		fmt.Println(errStr)
		return "", errors.New("Upload file failed. Please try again later.")
	}

	randomStr := String(20)

	godotenv.Load(".env")

	s3Config := &aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("ACCESS_KEY_IAM"), os.Getenv("SECRET_KEY_IAM"), ""),
	}
	s3Session := session.New(s3Config)

	uploader := s3manager.NewUploader(s3Session)

	input := &s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("AWS_BUCKET_NAME")),                             // bucket's name
		Key:         aws.String(folderPath + "/" + randomStr + "-" + fileheader.Filename), // files destination location
		Body:        file,                                                                 // content of the file
		ContentType: aws.String("image/jpg"),                                              // content type
	}
	res, errUploadWithContext := uploader.UploadWithContext(context.Background(), input)
	if errUploadWithContext != nil {
		errStr = "Upload file failed. Field < " + imageField + " >, failed upload to server. Sent default image url. Error detail : " + errUploadWithContext.Error()
		fmt.Println(errStr)
		return "", errors.New("Upload file failed. Please try again later.")
	}
	// RETURN URL LOCATION IN AWS
	return res.Location, nil
}
