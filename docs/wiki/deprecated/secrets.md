### Secret Management

> ⚠️ **Notice** ⚠️: this document contains a deprecated functionality that is still
> supported by New Relic for backwards compatibility. However, we encourage you to
> use the improved, fully-supported [secrets management functionality for On-Host Integrations](https://docs.newrelic.com/docs/integrations/host-integrations/installation/secrets-management). 

Flex Secret Management supports:
* Hashicorp Vault
* AWS KMS

#### AWS KMS
```yaml
### kms secret example
---
name: awskmsExample
secrets:
  myawskey: # specify a namespace for your key
    kind: aws-kms # set kind of secret
    region: ap-southeast-2 # region required
    ### can specify your own credential and/or config file
    # credential_file: "./myAWSCredentialFile"
    # config_file: "./myAWSConfigFile"
    ### need to set one of three options, http, file, or data to decrypt
    http: 
      url: https://some-special-s3-bucket.s3-ap-southeast-2.amazonaws.com/some.file ### download a file containing a secret
      ### other optional configuration for http
      # headers:
      #   myHeader: ABC
      # tls_config:
      #   insecure_skip_verify: true
      #   ca: "./pathToCa"
    # file: ./some.file ### path to file
    # data: AQICAHgdUAUlK7RGdwKFdBCTfRCsNk3oNtTv7FWsrP5VgoCFUgHPmcKgBpjzi7sdnDvV+RipAAAAZzBlBgkqhkiG9w0BBwagWDBWAgEAMFEGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMDUJMrVgA36DJmdflAgEQgCQIlM5F0zpUbH0MgWUskXjb4GJmAOvLgJUkMf5SJDmE2sceuCs=
    ### the exact returned value will be accessible via ${secret.<namespace>:value} eg. ${secret.myawskey:value} 
    ### if you want to unpack the decrypted contents further, set a "type"
    ### if the type is set to json, the json attributes will become available to select 
    ### eg. {"abc":"def"} returned the "def" value can be accessed by setting ${secret.myawskey:abc} and "def" will be substituted in
    ### likewise if you have a value that is split by equals that looks like "abc=def" or even "abc=def,hello=hi"
    ### ${secret.myawskey:abc} will be substituted to "def" and ${secret.myawskey:hello} will be substituted "hi"
    type: equal 
custom_attributes: # applies to all apis
  myCustAttr: myCustVal
apis: 
  - event_type: awskmsExample
    url: https://jsonplaceholder.typicode.com/todos/1
    custom_attributes:
      nestedCustAttr: nestedCustVal # nested custom attributes specific to each api
      secretFull: ${secret.myawskey:value} # 
      secretOther: ${secret.myawskey:hello} #
```

#### Vault
```yaml
### Hashicorp Vault Secret Example
---
name: vaultExample
secrets:
  myVaultKey: # specify a namespace for your key
    kind: vault # set kind of secret
    http: 
      url: http://localhost:1234/v1/newengine/data/secondsecret ### v2 engine format
      headers:
        X-Vault-Token: myroot
      # V2 Engine example payload returned
      #  ${secret.myVaultKey:abc} would return "hello"
      #  ${secret.myVaultKey:def} would return "bye"
      # {
      #   "request_id": "96bd7c0a-3a93-cd41-363f-e43710e05ac8",
      #   "lease_id": "",
      #   "renewable": false,
      #   "lease_duration": 0,
      #   "data": {
      #     "data": {
      #       "abc": "hello",
      #       "def": "bye"
      #     },
      #     "metadata": {
      #       "created_time": "2019-08-08T14:04:12.755825Z",
      #       "deletion_time": "",
      #       "destroyed": false,
      #       "version": 2
      #     }
      #   },
      #   "wrap_info": null,
      #   "warnings": null,
      #   "auth": null
      # }
      ### V1 Engine
      # url: http://localhost:1234/v1/kv/somepath ### v1 engine format
      # headers:
      #   X-Vault-Token: myroot
      # V1 Engine example payload returned
      #  ${secret.myVaultKey:key1} would return "hello"
      #  ${secret.myVaultKey:key2} would return "bye"
      # {
      #   "request_id": "b80581c9-69ed-2f48-13a2-c1c7a064f26f",
      #   "lease_id": "",
      #   "renewable": false,
      #   "lease_duration": 2764800,
      #   "data": {
      #     "key1": "hello",
      #     "key2": "bye"
      #   },
      #   "wrap_info": null,
      #   "warnings": null,
      #   "auth": null
      # }
apis: 
  - event_type: vaultExample
    url: https://jsonplaceholder.typicode.com/todos/1
    custom_attributes:
      nestedCustAttr: nestedCustVal # nested custom attributes specific to each api
      secretOne: ${secret.myVaultKey:abc} 
      secretTwo: ${secret.myVaultKey:def}
```