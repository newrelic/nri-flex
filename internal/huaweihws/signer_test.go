package huaweihws

import (
	"testing"

	"github.com/parnurzeal/gorequest"
	"gotest.tools/assert"
)

func TestSign(t *testing.T) {
	expectedSignature := `SDK-HMAC-SHA256 Access=keyabc111, SignedHeaders=content-type;header1;x-sdk-date, Signature=c5f664f53d5f02d79428e5d7d188ecb9f4018ce0fbdf0ac46673362d952cce81`

	reqURL := "dummyURL.com"
	signer := Signer{
		Key:    "keyabc111",
		Secret: "secretxyz999",
	}
	request := gorequest.New()
	request = request.Get(reqURL)
	request = request.Set("header1", "headerValue1")
	request = request.Set("X-Sdk-Date", "20201117T025305Z")
	r, _ := request.MakeRequest()
	_ = signer.Sign(r)

	assert.Equal(t, expectedSignature, r.Header.Get(HeaderAuthorization))

}
