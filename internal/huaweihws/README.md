# huaweihws
--
    import "."


## Usage

```go
const (
	BasicDateFormat     = "20060102T150405Z"
	Algorithm           = "SDK-HMAC-SHA256"
	HeaderXDate         = "X-Sdk-Date"
	HeaderHost          = "host"
	HeaderAuthorization = "Authorization"
	HeaderContentSha256 = "X-Sdk-Content-Sha256"
)
```
BasicDateFormat

#### func  AuthHeaderValue

```go
func AuthHeaderValue(signature, accessKey string, signedHeaders []string) string
```
AuthHeaderValue Get the finalized value for the "Authorization" header. The
signature parameter is the output from SignStringToSign

#### func  CanonicalHeaders

```go
func CanonicalHeaders(r *http.Request, signerHeaders []string) string
```
CanonicalHeaders func

#### func  CanonicalQueryString

```go
func CanonicalQueryString(r *http.Request) string
```
CanonicalQueryString func

#### func  CanonicalRequest

```go
func CanonicalRequest(r *http.Request, signedHeaders []string) (string, error)
```
CanonicalRequest func

#### func  CanonicalURI

```go
func CanonicalURI(r *http.Request) string
```
CanonicalURI returns request uri

#### func  HexEncodeSHA256Hash

```go
func HexEncodeSHA256Hash(body []byte) (string, error)
```
HexEncodeSHA256Hash returns hexcode of sha256

#### func  RequestPayload

```go
func RequestPayload(r *http.Request) ([]byte, error)
```
RequestPayload func

#### func  SignStringToSign

```go
func SignStringToSign(stringToSign string, signingKey []byte) (string, error)
```
SignStringToSign Create the HWS Signature.

#### func  SignedHeaders

```go
func SignedHeaders(r *http.Request) []string
```
SignedHeaders func

#### func  StringToSign

```go
func StringToSign(canonicalRequest string, t time.Time) (string, error)
```
StringToSign Create a "String to Sign".

#### type Signer

```go
type Signer struct {
	Key    string
	Secret string
}
```

Signer Signature HWS meta

#### func (*Signer) Sign

```go
func (s *Signer) Sign(r *http.Request) error
```
Sign SignRequest set Authorization header
