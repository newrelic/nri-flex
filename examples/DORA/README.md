# DORA Metrics Integrations 

## Flex Integrations
- GitLab DORA API
    - Group: query dora metrics in all repos in a single group  
    - Repo: query dora metrics in a single repo

## Requirements
[DORA Entity Definition](https://github.com/newrelic/entity-definitions/blob/main/definitions/ext-dora/definition.yml)
- Entity Creation
    - eventType: DoraMetricsSample
    - doraName: user defined doraName
        - doraName is used as an encapsulation, send metrics from mulitple sources to create a complete view on DORA metrics
    - must contain 1 of the 4 DORA Metrics to create an entity (recommended to send an event with all 4)
        - leadTimeForChanges
        - deploymentFrequency
        - timeToRestoreService
        - changeFailureRate
    - add a source attribute to allow filtering of different doraSources under the same entity
    - Optional: 
        - Send supporting DORA metrics from different sources.
            - i.e. ci/cd tools, custom scripts, different repositories

### **IMPORTANT** 
**doraName** needs to be a unique identifier for your entity creation

### Example Payloads
``` json
[
	{
		"eventType": "DoraMetricsSample",
		"doraName": "Opencart",
		"projectId": 123456,
		"repoUrl": "https://gitlab.com/your/repo",
		"repoName": "opencart-tf",
		"deploymentFrequency": 15,
		"leadTimeForChanges": 12,
		"timeToRestoreService": 0.4,
		"changeFailureRate": 12.6,
		"team": "dev-rel-apj",
		"pipeline": "production-opencart",
		"org": "dev-rel",
		"source": "GitLab"
	},
	{
		"eventType": "DoraMetricsSample",
		"doraName": "Opencart",
		"repoName": "opencart-tf-infra",
		"deploymentFrequency": 15,
		"team": "dev-rel-sre",
		"pipeline": "production-opencart-infra",
		"org": "dev-rel",
		"source": "GitLab"
	},
		{
		"eventType": "DoraMetricsSample",
		"doraName": "Opencart",
		"deploymentFrequency": 1,
		"team": "dev-rel-sre",
		"pipeline": "production-opencart-ansible",
		"org": "dev-rel",
		"source": "ansible"
	}
]
```

