# Flex-Siebel-Integration
Siebel Extension with New Relic Flex

## Technical Details:
This does NOT support Siebel prior to IP17 due the SMC (Siebel Management Console) not being included. 
The below was tested on **Siebel 19.6** and has not been tested on newer versions.

### Example integration

> Note that this example uses features that are [experimental](https://github.com/newrelic/nri-flex/tree/master/docs/experimental) (not officially supported).

This integration extracts metrics from your existing Siebel Installation into New Relic Platform. Metrics are extrected as custom events - as demonstrated in the screenshot below.

<img width="1851" alt="siebel_extension" src="https://user-images.githubusercontent.com/52437445/198934301-714b8e7d-40af-4215-9bef-c43a614be6d6.png">

### Sample of the Payload you receive from each of the endpoints are as below: 

## Tasks
This URL will retrieve the tasks, across all servers across the enterprise.
This is the same as doing ```list tasks``` in Siebel server manager. 

```https://tstgw1.gateway.siebel.net/siebel/v1.0/cloudgateway/enterprises/ENTERPRISENAME/tasks/```

### Sample Payload Event Type: Tasks

```json
[
  {
    "results": [
      {
        "events": [
          {
            "CC_ALIAS": "SCBroker",
            "CC_INCARN_NO": 0,
            "CC_NAME": "Siebel Connection Broker",
            "CC_RUNMODE": "Background",
            "CT_NAME": "SCBroker",
            "EN_NAME": "TSTENT",
            "SV_NAME": "TSTSERVER01",
            "TK_DISP_RUNSTATE": "Running",
            "TK_END_TIME": "",
            "TK_LABEL": "",
            "TK_PARENT_TASKNUM": "",
            "TK_PID": 20488,
            "TK_RUNSTATE": "Running",
            "TK_START_TIME": "2022-10-28 14:51:25",
            "TK_STATUS": "",
            "TK_TASKID": 2097155,
            "TK_TASKTYPE": "Normal",
            "TK_TID": 22316,
            "agentName": "Infrastructure",
            "agentVersion": "1.28.4",
            "api.StatusCode": 200,
            "baseUrl": "https://tstgw1.gateway.usertest.net/siebel/v1.0/cloudgateway/enterprises/TSTENT/",
            "coreCount": "2",
            "criticalViolationCount": 0,
            "entityGuid": "MTYxOTg4MHxJTkZSQXxOQXw2NjA1MjY1MTQxNDEwNTI3NTQy",
            "entityId": "6605265141410527542",
            "entityKey": "AUTSTSBLTSTAPP1.usertest.net",
            "entityName": "AUTSTSBLTSTAPP1",
            "event_type": "Siebel_Server_Tasks",
            "fullHostname": "AUTSTSBLTSTAPP1.usertest.net",
            "hostStatus": "running",
            "hostname": "AUTSTSBLTSTAPP1",
            "instanceType": "VMware, Inc. VMware7,1",
            "integrationName": "com.newrelic.nri-flex",
            "integrationVersion": "1.5.1",
            "nr.entityType": "HOST",
            "nr.ingestTimeMs": 1667178571000,
            "nr.invalidAttributeCount": 2,
            "operatingSystem": "windows",
            "processorCount": "2",
            "systemMemoryBytes": "25768775680",
            "timestamp": 1667178619000,
            "warningViolationCount": 0,
            "windowsFamily": "Server",
            "windowsPlatform": "Microsoft Windows Server 2016 Datacenter",
            "windowsVersion": "10.0.14393 Build 14393"
          }
        ]
      }
    ],
    "metadata": {
      "contents": [
        {
          "function": "events",
          "limit": 1,
          "order": {
            "column": "timestamp",
            "descending": true
          }
        }
      ],
      "eventTypes": [
        "Siebel_Server_Tasks"
      ],
      "eventType": "Siebel_Server_Tasks",
      "openEnded": true,
      "messages": [],
      "beginTimeMillis": 1667178618520,
      "endTimeMillis": 1667178678520,
      "beginTime": "2022-10-31T01:10:18Z",
      "endTime": "2022-10-31T01:11:18Z",
      "guid": "1dc634f6-c936-1e3f-5bfc-33d51ed17f87",
      "routerGuid": "1dc634f6-c936-1e3f-5bfc-33d51ed17f87",
      "rawSince": "1 MINUTES AGO",
      "rawUntil": "NOW",
      "rawCompareWith": "",
      "accounts": [
        1234567
      ]
    },
    "performanceStats": {
      "inspectedCount": 1202,
      "responseTime": 29
    }
  }
]
```

## Sessions
This URL will retrieve the sessions, across all servers across the enterprise.
This is the same as doing ```list sessions``` in Siebel server manager. 

```https://tstgw1.gateway.siebel.net/siebel/v1.0/cloudgateway/enterprises/ENTERPRISENAME/sesions/```

### Sample Payload Event Type: Session

```json
[
  {
    "results": [
      {
        "events": [
          {
            "CC_ALIAS": "SRProc",
            "CG_ALIAS": "SystemAux",
            "DB_SESSION_ID": "",
            "OM_APPLET": "",
            "OM_BUSCOMP": "",
            "OM_BUSSVC": "",
            "OM_LOGIN": "Forwarding Task",
            "OM_VIEW": "",
            "SV_NAME": "TSTSERVER01",
            "TK_DISP_RUNSTATE": "Running",
            "TK_HUNG_STATE": "",
            "TK_IDLE_STATE": "FALSE",
            "TK_PID": 13648,
            "TK_PING_TIME": "",
            "TK_TASKID": 3145733,
            "agentName": "Infrastructure",
            "agentVersion": "1.28.4",
            "api.StatusCode": 200,
            "baseUrl": "https://tstgw1.gateway.usertest.net/siebel/v1.0/cloudgateway/enterprises/TSTENT/",
            "coreCount": "2",
            "criticalViolationCount": 0,
            "entityGuid": "MTYxOTg4MHxJTkZSQXxOQXw2NjA1MjY1MTQxNDEwNTI3NTQy",
            "entityId": "6605265141410527542",
            "entityKey": "AUTSTSBLTSTAPP1.usertest.net",
            "entityName": "AUTSTSBLTSTAPP1",
            "event_type": "Siebel_Server_Sessions",
            "fullHostname": "AUTSTSBLTSTAPP1.usertest.net",
            "hostStatus": "running",
            "hostname": "AUTSTSBLTSTAPP1",
            "instanceType": "VMware, Inc. VMware7,1",
            "integrationName": "com.newrelic.nri-flex",
            "integrationVersion": "1.5.1",
            "nr.entityType": "HOST",
            "nr.ingestTimeMs": 1667191082000,
            "nr.invalidAttributeCount": 2,
            "operatingSystem": "windows",
            "processorCount": "2",
            "systemMemoryBytes": "25768775680",
            "timestamp": 1667191129000,
            "warningViolationCount": 0,
            "windowsFamily": "Server",
            "windowsPlatform": "Microsoft Windows Server 2016 Datacenter",
            "windowsVersion": "10.0.14393 Build 14393"
          }
        ]
      }
    ],
    "metadata": {
      "contents": [
        {
          "function": "events",
          "limit": 1,
          "order": {
            "column": "timestamp",
            "descending": true
          }
        }
      ],
      "eventTypes": [
        "Siebel_Server_Sessions"
      ],
      "eventType": "Siebel_Server_Sessions",
      "openEnded": true,
      "messages": [],
      "beginTimeMillis": 1667187553175,
      "endTimeMillis": 1667191153175,
      "beginTime": "2022-10-31T03:39:13Z",
      "endTime": "2022-10-31T04:39:13Z",
      "guid": "73832c0d-3049-d284-7cf0-6cb416814f7d",
      "routerGuid": "73832c0d-3049-d284-7cf0-6cb416814f7d",
      "rawSince": "60 MINUTES AGO",
      "rawUntil": "NOW",
      "rawCompareWith": "",
      "accounts": [
        1234567
      ]
    },
    "performanceStats": {
      "inspectedCount": 11234,
      "responseTime": 33
    }
  }
]
```

## Server Status
This URL will retrieve the sever status, across all servers across the enterprise.
This is the same as doing ```server status``` in Siebel server manager. 

```https://tstgw1.gateway.siebel.net/siebel/v1.0/cloudgateway/enterprises/ENTERPRISENAME/servers/```

### Sample Payload Event Type: Server Status

```json
[
  {
    "results": [
      {
        "events": [
          {
            "END_TIME": "",
            "ENTSRVR_NAME": "TSTENT",
            "HOST_NAME": "AUTSTSBLTSTAPP3",
            "INSTALL_DIR": "d:\\Siebel\\ent\\siebsrvr",
            "ROW_ID": 7,
            "SBLMGR_PID": 10704,
            "SBLSRVR_GROUP_NAME": "",
            "SBLSRVR_NAME": "TSTSERVER03",
            "SBLSRVR_STATE": "Running",
            "SBLSRVR_STATUS": "19.6.0.0 [23073] LANG_INDEPENDENT",
            "START_TIME": "2022-10-28 22:56:14",
            "SV_DISP_STATE": "Running",
            "agentName": "Infrastructure",
            "agentVersion": "1.28.4",
            "api.StatusCode": 200,
            "baseUrl": "https://tstgw1.gateway.usertest.net/siebel/v1.0/cloudgateway/enterprises/TSTENT/",
            "coreCount": "2",
            "criticalViolationCount": 0,
            "entityGuid": "MTYxOTg4MHxJTkZSQXxOQXw2NjA1MjY1MTQxNDEwNTI3NTQy",
            "entityId": "6605265141410527542",
            "entityKey": "AUTSTSBLTSTAPP1.usertest.net",
            "entityName": "AUTSTSBLTSTAPP1",
            "event_type": "Siebel_Server_Status",
            "fullHostname": "AUTSTSBLTSTAPP1.usertest.net",
            "hostStatus": "running",
            "hostname": "AUTSTSBLTSTAPP1",
            "instanceType": "VMware, Inc. VMware7,1",
            "integrationName": "com.newrelic.nri-flex",
            "integrationVersion": "1.5.1",
            "nr.entityType": "HOST",
            "nr.ingestTimeMs": 1667191261000,
            "nr.invalidAttributeCount": 2,
            "operatingSystem": "windows",
            "processorCount": "2",
            "systemMemoryBytes": "25768775680",
            "timestamp": 1667191309000,
            "warningViolationCount": 0,
            "windowsFamily": "Server",
            "windowsPlatform": "Microsoft Windows Server 2016 Datacenter",
            "windowsVersion": "10.0.14393 Build 14393"
          }
        ]
      }
    ],
    "metadata": {
      "contents": [
        {
          "function": "events",
          "limit": 1,
          "order": {
            "column": "timestamp",
            "descending": true
          }
        }
      ],
      "eventTypes": [
        "Siebel_Server_Status"
      ],
      "eventType": "Siebel_Server_Status",
      "openEnded": true,
      "messages": [],
      "beginTimeMillis": 1667187794037,
      "endTimeMillis": 1667191394037,
      "beginTime": "2022-10-31T03:43:14Z",
      "endTime": "2022-10-31T04:43:14Z",
      "guid": "263bcd7f-a990-8da7-3bb8-68d8bc25c9ca",
      "routerGuid": "263bcd7f-a990-8da7-3bb8-68d8bc25c9ca",
      "rawSince": "60 MINUTES AGO",
      "rawUntil": "NOW",
      "rawCompareWith": "",
      "accounts": [
        1234567
      ]
    },
    "performanceStats": {
      "inspectedCount": 444,
      "responseTime": 35
    }
  }
]
```

## Components
This URL will retrieve the component status, across all servers across the enterprise.
This is the same as doing ```list components``` in Siebel server manager. 

```https://tstgw1.gateway.siebel.net/siebel/v1.0/cloudgateway/enterprises/ENTERPRISENAME/components/```

### Sample Payload Event Type: Components

```json
[
  {
    "results": [
      {
        "events": [
          {
            "CC_ALIAS": "XMLPReportServer",
            "CC_DESC_TEXT": "",
            "CC_INCARN_NO": "",
            "CC_NAME": "XMLP Report Server",
            "CC_RUNMODE": "Batch",
            "CG_ALIAS": "XMLPReport",
            "CP_ACTV_MTS_PROCS": 1,
            "CP_DISP_RUN_STATE": "Online",
            "CP_END_TIME": "",
            "CP_MAX_MTS_PROCS": 1,
            "CP_MAX_TASKS": 20,
            "CP_NUM_RUN_TASKS": 0,
            "CP_RUN_STATE": "Online",
            "CP_STARTMODE": "Auto",
            "CP_START_TIME": "2022-10-28 14:51:31",
            "CP_STATUS": "Enabled",
            "CT_ALIAS": "BusSvcMgr",
            "ROW_ID": "1:24",
            "SV_NAME": "TSTSERVER01",
            "agentName": "Infrastructure",
            "agentVersion": "1.28.4",
            "api.StatusCode": 200,
            "baseUrl": "https://tstgw1.gateway.usertest.net/siebel/v1.0/cloudgateway/enterprises/TSTENT/",
            "coreCount": "2",
            "criticalViolationCount": 0,
            "entityGuid": "MTYxOTg4MHxJTkZSQXxOQXw2NjA1MjY1MTQxNDEwNTI3NTQy",
            "entityId": "6605265141410527542",
            "entityKey": "AUTSTSBLTSTAPP1.usertest.net",
            "entityName": "AUTSTSBLTSTAPP1",
            "event_type": "Siebel_Server_Components",
            "fullHostname": "AUTSTSBLTSTAPP1.usertest.net",
            "hostStatus": "running",
            "hostname": "AUTSTSBLTSTAPP1",
            "instanceType": "VMware, Inc. VMware7,1",
            "integrationName": "com.newrelic.nri-flex",
            "integrationVersion": "1.5.1",
            "nr.entityType": "HOST",
            "nr.ingestTimeMs": 1667191441000,
            "nr.invalidAttributeCount": 2,
            "operatingSystem": "windows",
            "processorCount": "2",
            "systemMemoryBytes": "25768775680",
            "timestamp": 1667191489000,
            "warningViolationCount": 0,
            "windowsFamily": "Server",
            "windowsPlatform": "Microsoft Windows Server 2016 Datacenter",
            "windowsVersion": "10.0.14393 Build 14393"
          }
        ]
      }
    ],
    "metadata": {
      "contents": [
        {
          "function": "events",
          "limit": 1,
          "order": {
            "column": "timestamp",
            "descending": true
          }
        }
      ],
      "eventTypes": [
        "Siebel_Server_Components"
      ],
      "eventType": "Siebel_Server_Components",
      "openEnded": true,
      "messages": [],
      "beginTimeMillis": 1667187940648,
      "endTimeMillis": 1667191540648,
      "beginTime": "2022-10-31T03:45:40Z",
      "endTime": "2022-10-31T04:45:40Z",
      "guid": "6b4f1dc1-b6b9-3a4c-1f34-725b62f75cba",
      "routerGuid": "6b4f1dc1-b6b9-3a4c-1f34-725b62f75cba",
      "rawSince": "60 MINUTES AGO",
      "rawUntil": "NOW",
      "rawCompareWith": "",
      "accounts": [
        1234567
      ]
    },
    "performanceStats": {
      "inspectedCount": 12351,
      "responseTime": 43
    }
  }
]
```

## Server
This URL will retrieve the Server Statistics, for a specific Server in your Enterprise.
This is the same as doing ```<SERVERNAME> statistics``` in Siebel server manager. 

```https://tstgw1.gateway.siebel.net/siebel/v1.0/cloudgateway/enterprises/ENTERPRISENAME/SERVERNAME/statistics/```

### Sample Payload Event Type: Server Statistics
```json
[
  {
    "results": [
      {
        "events": [
          {
            "CURR_VAL": 0,
            "SBLSRVR_NAME": "TSTSERVER01",
            "SD_DATATYPE": "Integer",
            "SD_DESC": "Number of inprocess transactions from PIM",
            "SD_SUBSYSTEM": "PIMSIEng",
            "SD_VISIBILITY": "Basic",
            "STAT_ALIAS": "NumPIMTransInProcess",
            "STAT_NAME": "Transactions from PIM in Progress",
            "agentName": "Infrastructure",
            "agentVersion": "1.28.4",
            "api.StatusCode": 200,
            "baseUrl": "https://tstgw1.gateway.usertest.net/siebel/v1.0/cloudgateway/enterprises/TSTENT/",
            "coreCount": "2",
            "criticalViolationCount": 0,
            "entityGuid": "MTYxOTg4MHxJTkZSQXxOQXw2NjA1MjY1MTQxNDEwNTI3NTQy",
            "entityId": "6605265141410527542",
            "entityKey": "AUTSTSBLTSTAPP1.usertest.net",
            "entityName": "AUTSTSBLTSTAPP1",
            "event_type": "Siebel_Server_Statistics_TSTSERVER01",
            "fullHostname": "AUTSTSBLTSTAPP1.usertest.net",
            "hostStatus": "running",
            "hostname": "AUTSTSBLTSTAPP1",
            "instanceType": "VMware, Inc. VMware7,1",
            "integrationName": "com.newrelic.nri-flex",
            "integrationVersion": "1.5.1",
            "nr.entityType": "HOST",
            "nr.ingestTimeMs": 1667191531000,
            "nr.invalidAttributeCount": 2,
            "operatingSystem": "windows",
            "processorCount": "2",
            "systemMemoryBytes": "25768775680",
            "timestamp": 1667191579000,
            "warningViolationCount": 0,
            "windowsFamily": "Server",
            "windowsPlatform": "Microsoft Windows Server 2016 Datacenter",
            "windowsVersion": "10.0.14393 Build 14393"
          }
        ]
      }
    ],
    "metadata": {
      "contents": [
        {
          "function": "events",
          "limit": 1,
          "order": {
            "column": "timestamp",
            "descending": true
          }
        }
      ],
      "eventTypes": [
        "Siebel_Server_Statistics_TSTSERVER01"
      ],
      "eventType": "Siebel_Server_Statistics_TSTSERVER01",
      "openEnded": true,
      "messages": [],
      "beginTimeMillis": 1667188043357,
      "endTimeMillis": 1667191643357,
      "beginTime": "2022-10-31T03:47:23Z",
      "endTime": "2022-10-31T04:47:23Z",
      "guid": "2a86b022-547d-4b2b-2f9c-4793af8ec5cd",
      "routerGuid": "2a86b022-547d-4b2b-2f9c-4793af8ec5cd",
      "rawSince": "60 MINUTES AGO",
      "rawUntil": "NOW",
      "rawCompareWith": "",
      "accounts": [
        1234567
      ]
    },
    "performanceStats": {
      "inspectedCount": 18392,
      "responseTime": 39
    }
  }
]
```

## Dashboard Included
We have included a .json for the Dashboard we have created based on the metrics received from this integration. You can fine tune it based on your requirements. 

<img width="1689" alt="siebel_dashboard" src="https://user-images.githubusercontent.com/52437445/199653062-20d79b80-5b28-4be5-940a-075fc739890a.png">

