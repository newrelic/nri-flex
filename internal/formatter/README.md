# formatter
--
    import "github.com/newrelic/nri-flex/internal/formatter"


## Usage

#### func  KvFinder

```go
func KvFinder(mode string, k1 string, k2 string) bool
```
KvFinder tests with multiple modes, whether k1 satisfies k2

#### func  PercToDecimal

```go
func PercToDecimal(v *interface{})
```
PercToDecimal convert percentage to decimal

#### func  RegMatch

```go
func RegMatch(text string, regexmatch string) []string
```
RegMatch Perform regex matching

#### func  RegSplit

```go
func RegSplit(text string, delimiter string) []string
```
RegSplit Split by Regex

#### func  SnakeCaseToCamelCase

```go
func SnakeCaseToCamelCase(key *string)
```
SnakeCaseToCamelCase converts snake_case to camelCase

#### func  SplitKey

```go
func SplitKey(key, splitChar string) (string, string, bool)
```
SplitKey simple key value pair splitter

#### func  ValueParse

```go
func ValueParse(v interface{}, regex string) string
```
ValueParse Plucks first found value out with regex, if nothing found send back
the value
