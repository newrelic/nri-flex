package aliyun

import (
	"testing"

	"github.com/parnurzeal/gorequest"
	"gotest.tools/assert"
)

func TestSign(t *testing.T) {
	expectedSignedUrl := `http://ecs.ap-southeast-2.aliyuncs.com/?AccessKeyId=testkey1&Action=DescribeRegions&Format=JSON&RegionId=ap-southeast-2&Signature=q7HtOs3jumOb33wnGptb7ihl4ak%3D&SignatureMethod=HMAC-SHA1&SignatureNonce=4751981231928800661&SignatureType=&SignatureVersion=1.0&Timestamp=2021-04-19T05%3A21%3A43Z&Version=2019-01-01`
	reqURL := "http://ecs.ap-southeast-2.aliyuncs.com/?Action=DescribeRegions&Format=JSON&RegionId=ap-southeast-2&SignatureNonce=4751981231928800661&Timestamp=2021-04-19T05%3A21%3A43Z&Version=2019-01-01"

	signer := Signer{
		Key:    "testkey1",
		Secret: "testSecret1",
	}
	request := gorequest.New()
	request = request.Get(reqURL)
	r, _ := request.MakeRequest()
	signedUrl, _ := signer.Sign(r)
	assert.Equal(t, expectedSignedUrl, signedUrl)

}
