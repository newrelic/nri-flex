// Aliyun signer

package aliyun

import (
	"crypto/hmac"
	"strconv"

	// #nosec
	"crypto/sha1"
	"encoding/base64"
	"fmt"

	// #nosec
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/newrelic/nri-flex/internal/load"
)

// BasicDateFormat
const (
	BasicDateFormat  = "2006-01-02T15:04:05Z0700"
	SignatureVersion = "1.0"
	APIVersion       = "2019-01-01"
	SignatureMethod  = "HMAC-SHA1"
)

// Signer Key and secret
type Signer struct {
	Key    string
	Secret string
}

// Sign URL
func (s *Signer) Sign(r *http.Request) (string, error) {
	var err error
	// #nosec
	var signedURL string
	reqMethod := r.Method
	paramsToSign := addSignerParams(r, s.Key)
	stringToSign := CanonicalQueryString(paramsToSign)
	stringToSign = reqMethod + "&%2F&" + stringToSign
	signature := ShaHmac1(stringToSign, s.Secret)
	if err != nil {
		return signedURL, err
	}
	paramsToSign.Add("Signature", signature)

	signedURL = r.URL.Scheme + "://" + r.URL.Host + r.URL.Path + "?" + paramsToSign.Encode()
	load.Logrus.Debugf("Aliyun Signer: SignedUrl: %v ", signedURL)
	return signedURL, nil
}

func ShaHmac1(source, secret string) string {
	/*
		https://www.alibabacloud.com/help/doc-detail/25492.htm?spm=a2c63.p38356.b99.700.21246d94fcNh6T

		Note When you calculate the signature, the key value specified by RFC 2104 is your AccessKeySecret with an ampersand (&) which has an ASCII value of 38. For more information, see Create an AccessKey pair.
	*/
	key := []byte(secret + "&")
	hmac := hmac.New(sha1.New, key)
	_, err := hmac.Write([]byte(source))
	if err != nil {
		load.Logrus.Errorf("Aliyun Signer: ShaHmac1 Error: %v", err)
	}
	signedBytes := hmac.Sum(nil)
	signedString := base64.StdEncoding.EncodeToString(signedBytes)

	load.Logrus.Debugf("Aliyun Signer: String to sign: %v", source)
	load.Logrus.Debugf("Aliyun Signer: Signature: %v", signedString)

	return signedString
}

func getNonce() string {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)
	return strconv.FormatInt(r.Int63(), 10)
}

var safeNonce = func(fn func() string) string {
	return fn()
}

func addSignerParams(r *http.Request, key string) url.Values {
	t := time.Now().UTC().Format(BasicDateFormat)
	signatureNonce := safeNonce(getNonce)
	params := r.URL.Query()
	if _, ok := params["Version"]; !ok {
		params.Add("Version", APIVersion)
	}
	if _, ok := params["SignatureVersion"]; !ok {
		params.Add("SignatureVersion", SignatureVersion)
	}
	if _, ok := params["SignatureMethod"]; !ok {
		params.Add("SignatureMethod", SignatureMethod)
	}
	if _, ok := params["AccessKeyId"]; !ok {
		params.Add("AccessKeyId", key)
	}
	if _, ok := params["Timestamp"]; !ok {
		params.Add("Timestamp", t)
	}
	if _, ok := params["SignatureNonce"]; !ok {
		params.Add("SignatureNonce", signatureNonce)
	}
	if _, ok := params["SignatureType"]; !ok {
		params.Add("SignatureType", "")
	}
	return params
}

func CanonicalQueryString(query url.Values) string {
	keys := make([]string, 0)

	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var a []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(query[key])
		for _, v := range query[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			a = append(a, kv)
		}
	}
	queryStr := strings.Join(a, "&")
	return url.QueryEscape(queryStr)
}
