# processor
--
    import "github.com/newrelic/nri-flex/internal/processor"


## Usage

#### func  AutoSetMetricAPI

```go
func AutoSetMetricAPI(currentSample *map[string]interface{}, api *load.API)
```
AutoSetMetricAPI automatically set metrics for use with the metric api

#### func  AutoSetMetricInfra

```go
func AutoSetMetricInfra(k string, v interface{}, metricSet *metric.Set, metrics map[string]string, autoSet bool, mode string)
```
AutoSetMetricInfra parse to number

#### func  AutoSetStandard

```go
func AutoSetStandard(currentSample *map[string]interface{}, api *load.API, workingEntity *integration.Entity, eventType string, config *load.Config)
```
AutoSetStandard x

#### func  CreateMetricSets

```go
func CreateMetricSets(samples []interface{}, config *load.Config, i int)
```
CreateMetricSets creates metric sets

#### func  FinalMerge

```go
func FinalMerge(data map[string]interface{}) []interface{}
```
FinalMerge Perform final data merging Separates detected samples and already
flattened attributes

#### func  FindStartKey

```go
func FindStartKey(mainDataset *map[string]interface{}, startKeys []string, inheritAttributes bool)
```
FindStartKey start at a different section of a payload

#### func  FlattenData

```go
func FlattenData(unknown interface{}, data map[string]interface{}, key string, sampleKeys map[string]string, api *load.API) map[string]interface{}
```
FlattenData flatten an interface

#### func  ProcessSamplesToMerge

```go
func ProcessSamplesToMerge(samplesToMerge *map[string][]interface{}, yml *load.Config)
```
ProcessSamplesToMerge used to merge multiple samples together

#### func  RunDataHandler

```go
func RunDataHandler(dataSets []interface{}, samplesToMerge *map[string][]interface{}, i int, cfg *load.Config)
```
RunDataHandler handles the data received for processing

#### func  RunEventFilter

```go
func RunEventFilter(filters []load.Filter, createEvent *bool, k string, v interface{})
```
RunEventFilter filters events generated

#### func  RunKeepKeys

```go
func RunKeepKeys(keepKeys []string, key *string, currentSample *map[string]interface{}, k *string)
```
RunKeepKeys Removes all other keys/attributes and keep only those defined in
keep_keys

#### func  RunKeyConversion

```go
func RunKeyConversion(key *string, api load.API, v interface{}, SkipProcessing *[]string)
```
RunKeyConversion handles to lower and snake to camel case for keys

#### func  RunKeyFilter

```go
func RunKeyFilter(filters []load.Filter, currentSample *map[string]interface{}, k string)
```
RunKeyFilter filters keys generated

#### func  RunKeyRemover

```go
func RunKeyRemover(removeKeys []string, key *string, progress *bool, currentSample *map[string]interface{})
```
RunKeyRemover Remove unwanted keys with regex

#### func  RunKeyRenamer

```go
func RunKeyRenamer(renameKeys map[string]string, key *string, originalKey string)
```
RunKeyRenamer find keys with regex, and replace the value

#### func  RunLazyFlatten

```go
func RunLazyFlatten(ds *map[string]interface{}, cfg *load.Config, api int)
```
RunLazyFlatten lazy flattens the payload

#### func  RunMathCalculations

```go
func RunMathCalculations(math *map[string]string, currentSample *map[string]interface{})
```
RunMathCalculations performs math calculations

#### func  RunPluckNumbers

```go
func RunPluckNumbers(v *interface{}, api load.API, key *string)
```
RunPluckNumbers pluck numbers out automatically with ValueParser eg.
"sample_start_time = 1552864614.137869 (Sun, 17 Mar 2019 23:16:54 GMT)" returns
1552864614.137869

#### func  RunSampleFilter

```go
func RunSampleFilter(sampleFilters []map[string]string, createSample *bool, key string, v interface{})
```
RunSampleFilter Filters samples generated

#### func  RunSampleRenamer

```go
func RunSampleRenamer(renameSamples map[string]string, currentSample *map[string]interface{}, key string, eventType *string)
```
RunSampleRenamer using regex if sample has a key that matches, make that a
different sample (event_type)

#### func  RunSubParse

```go
func RunSubParse(subParse []load.Parse, currentSample *map[string]interface{}, key string, v interface{})
```
RunSubParse splits nested values out from one line eg.
db0:keys=1,expires=0,avg_ttl=0

#### func  RunValConversion

```go
func RunValConversion(v *interface{}, api load.API, key *string)
```
RunValConversion performs percentage to decimal & nano second to millisecond

#### func  RunValueParser

```go
func RunValueParser(v *interface{}, api load.API, key *string)
```
RunValueParser use regex to find a key, and pluck out its value by regex

#### func  RunValueTransformer

```go
func RunValueTransformer(v *interface{}, api load.API, key *string)
```
RunValueTransformer use regex to find a key, and then transform the value eg.
key: world key: hello-${value} == key: hello-world

#### func  SetEventType

```go
func SetEventType(currentSample *map[string]interface{}, eventType *string, apiEventType string, apiMerge string, apiName string)
```
SetEventType sets the metricSet's eventType

#### func  StoreLookups

```go
func StoreLookups(storeLookups map[string]string, key *string, lookupStore *map[string][]string, v *interface{})
```
StoreLookups if key is found (using regex), store the values in the lookupStore
as the defined lookupStoreKey for later use

#### func  StripKeys

```go
func StripKeys(ds *map[string]interface{}, stripKeys []string)
```
StripKeys strip defined keys out

#### func  VariableLookups

```go
func VariableLookups(variableLookups map[string]string, key *string, variableStore *map[string]string, v *interface{})
```
VariableLookups if key is found (using regex), store the value in the
variableStore, as the defined by the variableStoreKey for later use
