//nolint:lll
package fixtures

var UrlTests = []struct {
	Name           string
	Endpoint       string
	Port           string
	Payload        string
	Config         string
	ExpectedStdout string
}{
	{
		Name:     "ECS Task State Change",
		Endpoint: "ecstask",
		Port:     "8000",
		Payload: `
{
  "version": "0",
  "id": "f8a990c2-2f93-4713-8b5d-d5b96f35bfd7",
  "detail-type": "ECS Task State Change",
  "source": "aws.ecs",
  "account": "123456789012",
  "time": "2016-09-15T21:57:35Z",
  "region": "us-east-1",
  "resources": [
    "arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf"
  ],
  "detail": {
    "clusterArn": "arn:aws:ecs:us-east-1:123456789012:cluster/cluster1",
    "containerInstanceArn": "arn:aws:ecs:us-east-1:123456789012:container-instance/04f8c17d-29e0-4711-aa74-852654e477ec",
    "containers": [
      {
        "containerArn": "arn:aws:ecs:us-east-1:123456789012:container/40a3b4bd-79ae-4472-a0be-816e5e0044a0",
        "lastStatus": "PENDING",
        "name": "test",
        "taskArn": "arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf"
      },
      {
        "containerArn": "arn:aws:ecs:us-east-1:123456789012:container/50a3b4bd-eqre-4472-a0be-qeqwe",
        "lastStatus": "PENDING",
        "name": "abc",
        "taskArn": "arn:aws:ecs:us-east-1:123456789012:task/3102878e-adfa-32-b9c1-qwerqe"
      }
    ],
    "createdAt": "2016-09-15T21:30:33.3Z",
    "desiredStatus": "RUNNING",
    "lastStatus": "PENDING",
    "overrides": {
      "containerOverrides": [
        {
          "command": [
            "command1",
            "command2"
          ],
          "environment": [
            {
              "name": "env1",
              "value": "value1"
            },
            {
              "name": "env2",
              "value": "value2"
            }
          ],
          "name": "test"
        }
      ]
    },
    "updatedAt": "2016-09-15T21:30:33.3Z",
    "taskArn": "arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf",
    "taskDefinitionArn": "arn:aws:ecs:us-east-1:123456789012:task-definition/testTD:1",
    "version": 1
  }
}
`,
		Config: `
---
integrations:
  - name: nri-flex
    config:
      name: EcsTaskChangeSample
      custom_attributes: # applies to all apis
        myCustAttr: myCustVal
      apis:
        - event_type: EcsTaskChangeSample
          url: http://127.0.0.1:8000/ecstask
          custom_attributes:
            nestedCustAttr: nestedCustVal # nested custom attributes specific to each api
          add_attribute:
            newAttr: myNewAttr_${detail.desiredStatus} 
`,
		ExpectedStdout: `{"name":"com.newrelic.nri-flex","protocol_version":"3","integration_version":"1.4.4","data":[{"metrics":[{"account":123456789012,"api.StatusCode":200,"containerArn":"arn:aws:ecs:us-east-1:123456789012:container/40a3b4bd-79ae-4472-a0be-816e5e0044a0","detail-type":"ECS Task State Change","detail.clusterArn":"arn:aws:ecs:us-east-1:123456789012:cluster/cluster1","detail.containerInstanceArn":"arn:aws:ecs:us-east-1:123456789012:container-instance/04f8c17d-29e0-4711-aa74-852654e477ec","detail.createdAt":"2016-09-15T21:30:33.3Z","detail.desiredStatus":"RUNNING","detail.lastStatus":"PENDING","detail.taskArn":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf","detail.taskDefinitionArn":"arn:aws:ecs:us-east-1:123456789012:task-definition/testTD:1","detail.updatedAt":"2016-09-15T21:30:33.3Z","detail.version":1,"event_type":"EcsTaskChangeSample","id":"f8a990c2-2f93-4713-8b5d-d5b96f35bfd7","integration_name":"com.newrelic.nri-flex","integration_version":"1.4.4","lastStatus":"PENDING","myCustAttr":"myCustVal","name":"test","nestedCustAttr":"nestedCustVal","newAttr":"myNewAttr_RUNNING","region":"us-east-1","source":"aws.ecs","taskArn":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf","time":"2016-09-15T21:57:35Z","version":0},{"account":123456789012,"api.StatusCode":200,"containerArn":"arn:aws:ecs:us-east-1:123456789012:container/50a3b4bd-eqre-4472-a0be-qeqwe","detail-type":"ECS Task State Change","detail.clusterArn":"arn:aws:ecs:us-east-1:123456789012:cluster/cluster1","detail.containerInstanceArn":"arn:aws:ecs:us-east-1:123456789012:container-instance/04f8c17d-29e0-4711-aa74-852654e477ec","detail.createdAt":"2016-09-15T21:30:33.3Z","detail.desiredStatus":"RUNNING","detail.lastStatus":"PENDING","detail.taskArn":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf","detail.taskDefinitionArn":"arn:aws:ecs:us-east-1:123456789012:task-definition/testTD:1","detail.updatedAt":"2016-09-15T21:30:33.3Z","detail.version":1,"event_type":"EcsTaskChangeSample","id":"f8a990c2-2f93-4713-8b5d-d5b96f35bfd7","integration_name":"com.newrelic.nri-flex","integration_version":"1.4.4","lastStatus":"PENDING","myCustAttr":"myCustVal","name":"abc","nestedCustAttr":"nestedCustVal","newAttr":"myNewAttr_RUNNING","region":"us-east-1","source":"aws.ecs","taskArn":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-adfa-32-b9c1-qwerqe","time":"2016-09-15T21:57:35Z","version":0},{"account":123456789012,"api.StatusCode":200,"commandSamples":"[map[:command1] map[:command2]]","detail-type":"ECS Task State Change","detail.clusterArn":"arn:aws:ecs:us-east-1:123456789012:cluster/cluster1","detail.containerInstanceArn":"arn:aws:ecs:us-east-1:123456789012:container-instance/04f8c17d-29e0-4711-aa74-852654e477ec","detail.createdAt":"2016-09-15T21:30:33.3Z","detail.desiredStatus":"RUNNING","detail.lastStatus":"PENDING","detail.taskArn":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf","detail.taskDefinitionArn":"arn:aws:ecs:us-east-1:123456789012:task-definition/testTD:1","detail.updatedAt":"2016-09-15T21:30:33.3Z","detail.version":1,"environmentSamples":"[map[name:env1 value:value1] map[name:env2 value:value2]]","event_type":"EcsTaskChangeSample","id":"f8a990c2-2f93-4713-8b5d-d5b96f35bfd7","integration_name":"com.newrelic.nri-flex","integration_version":"1.4.4","myCustAttr":"myCustVal","name":"test","nestedCustAttr":"nestedCustVal","newAttr":"myNewAttr_RUNNING","region":"us-east-1","source":"aws.ecs","time":"2016-09-15T21:57:35Z","version":0},{"":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf","account":123456789012,"api.StatusCode":200,"detail-type":"ECS Task State Change","detail.clusterArn":"arn:aws:ecs:us-east-1:123456789012:cluster/cluster1","detail.containerInstanceArn":"arn:aws:ecs:us-east-1:123456789012:container-instance/04f8c17d-29e0-4711-aa74-852654e477ec","detail.createdAt":"2016-09-15T21:30:33.3Z","detail.desiredStatus":"RUNNING","detail.lastStatus":"PENDING","detail.taskArn":"arn:aws:ecs:us-east-1:123456789012:task/3102878e-4af2-4b3c-b9c1-2556b95b2bbf","detail.taskDefinitionArn":"arn:aws:ecs:us-east-1:123456789012:task-definition/testTD:1","detail.updatedAt":"2016-09-15T21:30:33.3Z","detail.version":1,"event_type":"EcsTaskChangeSample","id":"f8a990c2-2f93-4713-8b5d-d5b96f35bfd7","integration_name":"com.newrelic.nri-flex","integration_version":"1.4.4","myCustAttr":"myCustVal","nestedCustAttr":"nestedCustVal","newAttr":"myNewAttr_RUNNING","region":"us-east-1","source":"aws.ecs","time":"2016-09-15T21:57:35Z","version":0},{"event_type":"flexStatusSample","flex.Hostname":"rocky","flex.IntegrationVersion":"1.4.4","flex.counter.ConfigsProcessed":1,"flex.counter.EcsTaskChangeSample":4,"flex.counter.EventCount":4,"flex.counter.EventDropCount":0,"flex.counter.HttpRequests":1,"flex.time.elaspedMs":48,"flex.time.endMs":1654079094326,"flex.time.startMs":1654079094278}],"inventory":{},"events":[]}]}`,
	},
}
