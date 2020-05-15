package api

import (
	// "bufio"
	// "bytes"
	//"errors"
	"fmt"
	"strings"

	//"github.com/astaxie/beego"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	cfg "github.com/pnetwork/sre.ceph.init/internal/config"
	//"io"
	//"os"
)

type CephMgmt struct {
	host       string `ceph host`
	DoName     string `ceph doname`
	bucket_id  string `bucket_id`
	PathStyle  bool   `ceph url style, true means you can use host directly,  false means bucket_id.doname will be used`
	AccessKey  string `aws s3 aceessKey`
	SecretKey  string `aws s3 secretKey`
	block_size int64  `block_size`
}

var Ceph CephMgmt

type MyProvider struct{}

func (m *MyProvider) Retrieve() (credentials.Value, error) {

	return credentials.Value{
		AccessKeyID:     Ceph.AccessKey,
		SecretAccessKey: Ceph.SecretKey,
	}, nil
}
func (m *MyProvider) IsExpired() bool { return false }

func (this *CephMgmt) UpdateCephMgmt(host string, bucket string, accesskey string, secretkey string) error {

	this.host = host
	this.bucket_id = bucket
	this.AccessKey = accesskey
	this.SecretKey = secretkey
	this.PathStyle = true
	this.block_size = 8 * 1024 * 1024

	return nil
}

func CephAPI(config *cfg.Config) {

	Ceph.UpdateCephMgmt(config.BaseConfig.PN_GLOBAL_STORAGE_ENDPOINT, "scripts", config.BaseConfig.PN_GLOBAL_STORAGE_SECRET_ID, config.BaseConfig.PN_GLOBAL_STORAGE_SECRET_KEY)
	cephsession, _ := Ceph.connect()

	cephfile_metrics := GetBucketsFIles(cephsession)
	fmt.Printf("Number of files in rook-ceph : %d \n", cephfile_metrics)

	//svc := s3.New(session.New(session))
	fmt.Println("----------Start to Create Bucket----------")
	CreateBucket(cephsession, config)
	fmt.Println("----------Start to Delete Bucket----------")
	DeleteBucket(cephsession, config)
}

func (this *CephMgmt) connect() (*s3.S3, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			LogLevel: aws.LogLevel(aws.LogDebug),
			Region:   aws.String("US"), //Required 目前尚未分区，填写default即可
			//EndpointResolver: endpoints.ResolverFunc(s3CustResolverFn),
			Endpoint:         &this.host,
			DisableSSL:       aws.Bool(true),
			S3ForcePathStyle: &this.PathStyle,
			Credentials:      credentials.NewCredentials(&MyProvider{}),
			//DisableParamValidation: aws.Bool(false),
		},
	}))
	// Create the S3 service client with the shared session. This will
	// automatically use the S3 custom endpoint configured in the custom
	// endpoint resolver wrapping the default endpoint resolver.
	return s3.New(sess), nil
}

func GetBucketsFIles(s3svc *s3.S3) int {

	//var bucket_name []string
	//i :=0
	total_file := 0
	result, err := s3svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println("Failed to list buckets", err)
		return total_file
	}

	fmt.Println("Buckets:")
	for _, bucket := range result.Buckets {
		fmt.Printf("%s : %s\n", aws.StringValue(bucket.Name), bucket.CreationDate)
		tempbucket := aws.StringValue(bucket.Name)
		total_file = total_file + ListObjects(s3svc, &tempbucket)
	}
	return total_file
}

/*
func testLocationConstraint(s3svc *s3.S3){
    BucketInput := s3.CreateBucketInput{
		Bucket: aws.String("examplebucket"),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String("us-east-1"),
		},
    }


}
*/

func CreateBucket(s3svc *s3.S3, config *cfg.Config) {

	BucketInput := &s3.CreateBucketInput{
		Bucket: aws.String("siang"),
		//GrantRead:    aws.String("GrantRead"),
		//GrantReadACP: aws.String("GrantReadACP"),
		//GrantWrite:   aws.String("GrantWrite"),
		//ACL:          aws.String("public-read | public-read-write"),
		CreateBucketConfiguration: &s3.CreateBucketConfiguration{
			LocationConstraint: aws.String(""),
		},
	}
	/*
		fmt.Printf("BucketInput.Bucket.len : %d\n", len(aws.StringValue(BucketInput.Bucket)))
		err := BucketInput.Validate()
		if err != nil {
			fmt.Println("BucketInput.Validate : ", err)
		}
	*/
	bucketlistArray := strings.Split(config.INIT_BUCKET_LIST, ",")
	if len(bucketlistArray) > 0 {
		for i := 0; i < len(bucketlistArray); i++ {
			if len(bucketlistArray[i]) > 0 {
				BucketInput.SetBucket(bucketlistArray[i])
				s3svc.CreateBucketRequest(BucketInput)
				result, err := s3svc.CreateBucket(BucketInput)
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						switch aerr.Code() {
						case s3.ErrCodeBucketAlreadyExists:
							fmt.Println(s3.ErrCodeBucketAlreadyExists, aerr.Error())
						case s3.ErrCodeBucketAlreadyOwnedByYou:
							fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
						default:
							fmt.Println(aerr)
							//fmt.Println(aerr.Error())
						}
					} else {
						fmt.Println(err)
					}

				}
				fmt.Println(result)
			} else if len(bucketlistArray[i]) == 0 {
				fmt.Println("Bucket name can't be null")
			}
		}
	} else if len(bucketlistArray) == 0 {
		fmt.Println("INIT_BUCKET_LIST is null , this module will not send request to create bucket")
	}
}

func DeleteBucket(s3svc *s3.S3, config *cfg.Config) {

	deletelistArray := strings.Split(config.ROMOVE_BUCKET_LIST, ",")

	if len(deletelistArray) > 0 {
		for i := 0; i < len(deletelistArray); i++ {
			if len(deletelistArray[i]) > 0 {
				bucketName := &deletelistArray[i]
				fmt.Printf("remove bucketName: %s/n", deletelistArray[i])
				_, err := s3svc.DeleteBucket(&s3.DeleteBucketInput{Bucket: bucketName})
				if err != nil {
					if aerr, ok := err.(awserr.Error); ok {
						switch aerr.Code() {
						case s3.ErrCodeBucketAlreadyOwnedByYou:
							fmt.Println(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
						default:
							//fmt.Println("rick")
							fmt.Println(aerr)
							//fmt.Println(aerr.Error())
						}
					}
				}
			} else if len(deletelistArray[i]) == 0 {
				fmt.Println("Bucket name can't be null")
			}

		}
	} else if len(deletelistArray) == 0 {
		fmt.Println("REMOVE_BUCKET_LIST is null , this module will not send request to delete bucket")
	}

}

func ListObjects(s3svc *s3.S3, bucketname *string) int {
	i := 0
	object_number := 0
	err := s3svc.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: bucketname,
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++

		//for _, _ = range p.Contents {
		for _, obj := range p.Contents {
			fmt.Println("Object:", *obj.Key)
			// fmt.Println("Object-SIze ", *obj.Size)
			object_number++
		}
		return
	})
	if err != nil {
		fmt.Println("failed to list objects", err)
		return 0
	}
	fmt.Printf("Object total: %d \n", object_number)
	return object_number
}
