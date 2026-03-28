package cos

import (
	"context"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	"time"
	"os"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/astaxie/beego"
)

type CosConfig struct {
	SecretId  string
	SecretKey string
	CosUrl    string
}

var (
	client = &cos.Client{}
)

func InitCos(cfg *CosConfig) error {
	u, _ := url.Parse(cfg.CosUrl)
	b := &cos.BaseURL{BucketURL: u}
	client = cos.NewClient(b, &http.Client{
		Timeout: 60 * time.Second,
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretId,
			SecretKey: cfg.SecretKey,
		},
	})
	return nil
}

//name 为存储桶中的路径,例：/img/123.jpg，就是img文件夹下的123.jpg。访问路径就是CosUrl+/img/123.jpg
//filePath 源文件路径 服务器保存的路径
func Upload(name, filePath string) bool {
	isTue := true
	if client == nil {
		isTue = false
		logs.Info("Cos client is nil")
	}
	_, err := client.Object.PutFromFile(context.Background(), name, filePath, nil)
	if err != nil {
		isTue = false
		logs.Error(fmt.Sprintf("Cos upload err:%s", err.Error()))
	}
	return isTue
}


//亚马逊上传存储桶
func UploadAwsFile(fileUrl,fileName string) string {
	var awsConfig *aws.Config
	accessKey := beego.AppConfig.String("cos::accessKey")
	accessSecret := beego.AppConfig.String("cos::accessSecret")
	myBucket := beego.AppConfig.String("cos::myBucket")
	region := beego.AppConfig.String("cos::region")
	if region == "" {
		region = "us-east-1"
	}
	if accessKey == "" || accessSecret == "" {
		//load default credentials
		awsConfig = &aws.Config{
			Region: aws.String(region),

		}
	} else {
		awsConfig = &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewStaticCredentials(accessKey, accessSecret, ""),
		}
	}

	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession(awsConfig))

	// Create an uploader with the session and default options
	//uploader := s3manager.NewUploader(sess)

	// Create an uploader with the session and custom options
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		u.PartSize = 5 * 1024 * 1024 // The minimum/default allowed part size is 5MB
		u.Concurrency = 2            // default is 5
	})

	//open the file
	f, err := os.Open(fileUrl)
	if err != nil {
		fmt.Printf("failed to open file %q, %v", fileUrl, err)
		return ""
	}
	//defer f.Close()

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(fileName),
		Body:   f,
	})

	//in case it fails to upload
	if err != nil {
		fmt.Printf("failed to upload file, %v", err)
		return ""
	}
	//fmt.Printf("file uploaded to, %s\n", result.Location)
	return result.Location
}

