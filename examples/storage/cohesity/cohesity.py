import json
import sys
import requests
import urllib3
import datetime
import os

# This script utilizes the Cohesity REST API to fetch performance stats: https://developer.cohesity.com/apidocs-641.html#/rest/getting-started
# Arguments specified in `cohesity.json`

os.environ['no_proxy'] = '*'


class cohesity():
    arrayid: object

    def __init__(self, ip, arrayid, cluster_alias, apikey):
        self.ip = ip
        self.arrayid = arrayid
        self.alias = cluster_alias
        # Used by simple GET
        self.headers = {'apikey': apikey}
        self.URL = 'https://%s' % ip
        self.token = 1
        urllib3.disable_warnings()

    def doGet(self, restURI):
        # Form the url by combining the base with the passed in part
        url = "%s/%s" % (self.URL, restURI)

        # print("URL in doGet", url)
        r = requests.get(url, headers=self.headers, verify=False, timeout=15)

        # Convert the JSON data into a python dict
        response = json.loads(r.text)
        return response

    def capacity(self):
        restURI = 'irisservices/api/v1/public/stats/storage'
        response = self.doGet(restURI)
        response["key"] = "capacity"
        response["array"] = self.arrayid
        response["used %"] = round(response["localUsageBytes"]/response["totalCapacityBytes"]*100, 2)
        inList = [response]
        return inList

    def protectionSummary(self):
        restURI = 'irisservices/api/v1/public/stats/protectionSummary'
        response = self.doGet(restURI)
        response["array"] = self.arrayid
        inList = []

        summary_block = {"key": "protectionSummary", "totalnumObjectsProtected": response["numObjectsProtected"],
                         "totalnumObjectsUnprotected": response["numObjectsUnprotected"],
                         "totalprotectedSizeBytes": response["protectedSizeBytes"],
                         "totalunprotectedSizeBytes": response["unprotectedSizeBytes"]}
        inList.append(summary_block)

        for block in response["statsByEnv"]:
            row = {"array": self.arrayid, "key": "protectionSummaryEnv", "env": block["environment"],
                   "numObjectsProtected": block["numObjectsProtected"],
                   "numObjectsUnprotected": block["numObjectsUnprotected"],
                   "protectedSizeBytes": block["protectedSizeBytes"],
                   "unprotectedSizeBytes": block["unprotectedSizeBytes"]
                   }
            inList.append(row)

        return inList

    def get_protection_jobs(self):
        restURI = 'irisservices/api/v1/public/protectionJobs'
        filerkeys = ["environment=kVMware",
                     "isActive=True"]
        if filerkeys:
            restURI = restURI + "?"
            for filteringkey in filerkeys:
                restURI = restURI + "&"
                restURI = restURI + filteringkey
        # print(restURI)
        inList = []
        response = self.doGet(restURI)
        for block in response:
            block["key"] = "protectionJobs"
            block["array"] = self.arrayid
            # print(block)
            try:
                filter = {  # "description": block["description"],
                    "environment": block["environment"],
                    "id": block["id"],
                    "isPaused": block["isPaused"],
                    "key": block["key"],
                    "name": block["name"]}
                inList.append(filter)
            except Exception as e:
                # print("Exception in get_protection_jobs:", e)
                pass

        return inList

    def get_alerts(self):
        restURI = 'irisservices/api/v1/public/alerts'
        response = self.doGet(restURI)
        for block in response:
            block["key"] = "alert"
            block["array"] = self.arrayid
            propertyListdic = {}
            propertyListAll = block["propertyList"]
            for property in propertyListAll:
                key = property["key"]
                value = property["value"]
                propertyListdic[key] = value

            block["properties"] = propertyListdic
            block["latest_alert_timestamp"] = int(block["latestTimestampUsecs"]/1000)

            if "Open" in block["alertState"]:
                block["stateID"] = 1
            elif "Resolved" in block["alertState"]:
                block["stateID"] = 0
            else:
                block["stateID"] = 99

            ### Removing unnessesary keys
            to_delete = ["cluster_id_str", "run_id", "run_url"]
            for key_to_delete in to_delete:
                if key_to_delete in block["properties"]:
                    del block["properties"][key_to_delete]

            del block["propertyList"]

        # inList = [response]
        return response

    def get_cluster_performance(self):
        restURI = 'irisservices/api/v1/public/cluster?fetchStats=true'
        response = self.doGet(restURI)
        res = response["stats"]
        res["key"] = "stats"
        res["array"] = self.arrayid

        filter = {"array": self.arrayid, "key": "stats",
                  "numBytesRead": res["usagePerfStats"]["numBytesRead"],
                  "numBytesWritten": res["usagePerfStats"]["numBytesWritten"],
                  "readIos": res["usagePerfStats"]["readIos"],
                  "readLatencyMsecs": res["usagePerfStats"]["readLatencyMsecs"],
                  "writeIos": res["usagePerfStats"]["writeIos"],
                  "writeLatencyMsecs": res["usagePerfStats"]["writeLatencyMsecs"],
                  }

        if "usagePerfStats" in res:
            if "dataInBytes" in res["usagePerfStats"]:
                filter.update({"dataInBytes": res["usagePerfStats"]["dataInBytes"]})
        if "usagePerfStats" in res:
            if "dataInBytesAfterReduction" in res["usagePerfStats"]:
                filter.update({"dataInBytesAfterReduction": res["usagePerfStats"]["dataInBytesAfterReduction"]})
        if "dataReductionRatio" in res:
            filter.update({"dataReductionRatio": res["dataReductionRatio"]})

        inList = [filter]
        return inList

    def get_services(self):
        restURI = 'irisservices/api/v1/public/clusters/services/states'
        response = self.doGet(restURI)

        for block in response:
            block["key"] = "services"
            block["array"] = self.arrayid

            if "Running" in block["state"]:
                block["stateID"] = 0
            elif "Restarting" in block["state"]:
                block["stateID"] = 1
            elif "Stopped" in block["state"]:
                block["stateID"] = 2
            else:
                block["stateID"] = 99

        inList = [response]
        return inList

    def get_license(self):
        restURI = 'irisservices/api/v1/public/licenseUsage'
        response = self.doGet(restURI)["licensedUsage"]
        for block in response:
            block["key"] = "license"
            block["array"] = self.arrayid
            block["expiryDate"] = str(datetime.datetime.fromtimestamp(block["expiryTime"]))
        inList = [response]
        return inList

    def get_scheduled(self):
        restURI = 'irisservices/api/v1/public/stats/protectionSummary'
        response = self.doGet(restURI)
        inList = [response]
        return inList

    def get_failed_objects_report(self):
        restURI = 'irisservices/api/v1/public/reports/protectionSourcesJobsSummary'
        try:
            response = self.doGet(restURI)
            for block in response["protectionSourcesJobsSummary"]:
                block["key"] = "failed_objects_report"
                block["array"] = self.arrayid

                if "Success" in block["lastRunStatus"]:
                    block["lastRunStatusID"] = 1
                elif "Warning" in block["lastRunStatus"]:
                    block["lastRunStatusID"] = 2
                elif "Error" in block["lastRunStatus"]:
                    block["lastRunStatusID"] = 3
                else:
                    block["lastRunStatusID"] = 0

                keys_to_remove = ["lastRunType", "numDataReadBytes", "numLogicalBytesProtected"]
                for key_to_remove in keys_to_remove:
                    if key_to_remove in block:
                        del block[key_to_remove]

                # Converting Microsec to milisec
                keys = ["firstSuccessfulRunTimeUsecs", "lastRunEndTimeUsecs", "lastRunStartTimeUsecs",
                        "lastSuccessfulRunTimeUsecs"]
                for key in keys:
                    if key in block:
                        block[key] = int(block[key] / 1000)

                # For Virtual machines only
                tags = ""

                if "vmWareProtectionSource" in block["protectionSource"]:

                    if "tagAttributes" in block["protectionSource"]["vmWareProtectionSource"]:
                        tagscount = len(block["protectionSource"]["vmWareProtectionSource"]["tagAttributes"])
                        i = 1
                        for tag in block["protectionSource"]["vmWareProtectionSource"]["tagAttributes"]:
                            if i == 1:
                                tags = str(tag["name"])
                                i = i + 1
                            elif i != tagscount:
                                tags = tags + "," + str(tag["name"])
                                i = i + 1
                            elif i == tagscount:
                                tags = tags + "," + str(tag["name"])


                    if "Connected" in block["protectionSource"]["vmWareProtectionSource"]["connectionState"]:
                        block["protectionSource"]["vmWareProtectionSource"]["connectionStateID"] = 1
                    else:
                        block["protectionSource"]["vmWareProtectionSource"]["connectionStateID"] = 2

                    keys_to_remove = ["virtualDisks", "id", "toolsRunningStatus", "version", "tagAttributes"]
                    for key_to_remove in keys_to_remove:
                        if key_to_remove in block["protectionSource"]["vmWareProtectionSource"]:
                            del block["protectionSource"]["vmWareProtectionSource"][key_to_remove]

                # Create new tags key
                block["protectionSource"]["tags"] = tags

                # For physical machines only:
                if "physicalProtectionSource" in block["protectionSource"]:
                    block["protectionSource"]["tags"] = "client_based_backup"
                    if "agents" in block["protectionSource"]["physicalProtectionSource"]:
                        if "Healthy" in block["protectionSource"]["physicalProtectionSource"]["agents"][0]["status"]:
                            block["protectionSource"]["physicalProtectionSource"]["healthID"] = 1
                        else:
                            block["protectionSource"]["physicalProtectionSource"]["healthID"] = 2
                    else:
                        # "Agent" key is not found
                        block["protectionSource"]["physicalProtectionSource"]["healthID"] = 3

                    # Delete unnecessary keys
                    keys_to_remove = ["networkingInfo", "volumes", "agents", "memorySizeBytes",
                                      "numProcessors", "type", "vsswriters"]
                    for key_to_remove in keys_to_remove:
                        if key_to_remove in block["protectionSource"]["physicalProtectionSource"]:
                            del block["protectionSource"]["physicalProtectionSource"][key_to_remove]

            # inList = [response]
            return response
        except KeyError as e:
            response = {"array": self.arrayid, "key": "failed_objects_report",
                        "error": "Keyerror", "errorID": 1}
            return response

    def get_nodes_disks(self):
        restURI = 'irisservices/api/v1/public/nodes?showSystemDisks=true'
        response = self.doGet(restURI)
        inList = []
        for block in response:
            # For drive
            for item in block["systemDisks"]:
                # If offline = 2, online 1
                if item["offline"]:
                    item["offlineID"] = 2
                elif not item["offline"]:
                    item["offlineID"] = 1
                else:
                    item["offlineID"] = 0

            filter = {"key": "node_disk", "array": self.arrayid,
                      "chassisID": block["chassisInfo"]["chassisId"],
                      "node_hostname": block["hostName"],
                      "node_serial": block["cohesityNodeSerial"],
                      "systemDisks": block["systemDisks"]}

            inList.append(filter)
        finalList = [inList]
        return finalList

    def get_cluster_nodes(self):
        restURI = 'irisservices/api/v1/public/cluster/status'
        response = self.doGet(restURI)
        inList = []
        for block in response["nodeStatuses"]:
            block["clusterCI"] = self.alias
            block["key"] = "nodes"

            if "services" in block:
                del block["services"]

            if block["inCluster"]:
                block["online"] = 1
            else:
                block["online"] = 2

        finalList = [response["nodeStatuses"]]
        return finalList

    def get_excluded_VMs(self):
        restURI = 'irisservices/api/v1/public/protectionJobs?isDeleted=false'

        response = self.doGet(restURI)
        alllist = []
        for block in response:
            if "excludeSourceIds" in block:
                for id in block["excludeSourceIds"]:
                    alllist.append({"id": id, "protectionGroupName":  block["name"]})

        return alllist

    def list_protection_source(self):
        rest = "irisservices/api/v1/public/protectionSources/objects"
        response = self.doGet(rest)

        for item in response:
            if "vmWareProtectionSource" in item:
                del item["vmWareProtectionSource"]
            if "parentId" in item:
                del item["parentId"]

        excludedVMlist = []
        elist = self.get_excluded_VMs()

        for item in response:
            for item2 in elist:

                if item["id"] == item2["id"]:
                    # print("match", item)
                    item["key"] = "excludedVM"
                    item["array"] = self.arrayid
                    item["protectionGroupName"] = item2["protectionGroupName"]
                    excludedVMlist.append(item)

        return excludedVMlist

    def main(self):
        finalList = []
        finalList.append(self.capacity())
        finalList.append(self.protectionSummary())
        finalList.append(self.get_protection_jobs())
        finalList.append(self.get_alerts())
        finalList.append(self.get_cluster_performance())
        finalList.append(self.get_services())
        finalList.append(self.get_license())
        finalList.append(self.get_failed_objects_report())
        finalList.append(self.get_nodes_disks())
        finalList.append(self.get_cluster_nodes())
        finalList.append(self.list_protection_source())
        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)


fqdn = str(sys.argv[1])
cluster_name = str(sys.argv[2])
cluster_alias = str(sys.argv[3])
func = str(sys.argv[4])
secret = str(sys.argv[5])


try:
    array_name = cohesity(fqdn, cluster_name, cluster_alias, secret)
    if "main" in func:
        array_name.main()
    if "protectionSummary" in func:
        print(json.dumps(array_name.protectionSummary(), indent=4, sort_keys=True))
    if "capacity" in func:
        print(json.dumps(array_name.capacity(), indent=4, sort_keys=True))
    if "get_scheduled" in func:
        print(json.dumps(array_name.get_scheduled(), indent=4, sort_keys=True))

except NameError as r:
    print("NameError error", r)
