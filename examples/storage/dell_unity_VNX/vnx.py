import json
from requests.auth import HTTPBasicAuth
import requests
import urllib3
import sys

# This script polls Dell Unity's REST API to collect performance stats
# Arguments specified in `vnx_arrays.json`

class VNX:
    def __init__(self, user, pw, ip, port, ssl, clusterCI):
        self.user = user
        self.pw = pw
        self.ip = ip
        self.port = port
        self.ssl = False
        self.clusterCI = clusterCI
        self.headers = {'Content-Type': 'application/json', 'X-EMC-REST-CLIENT': 'true'}
        self.URL = 'https://%s:%s' % (ip, port)

    def doGet(self, restURI, payload=None):
        urllib3.disable_warnings()
        # Form the url by combining the base with the passed in part
        url = "%s/%s" % (self.URL, restURI)
        # Do the GET. verify=False lets us ignore the SSL Cert error

        r = requests.get(url, headers=self.headers, verify=False,
                         auth=HTTPBasicAuth(self.user, self.pw),
                         )
        if r.status_code != 200:
            print(r.text)
            raise Exception
        else:
            response = r.text
            response = json.loads(response)
        return response

    def doPost(self, restURI, *argv):
        urllib3.disable_warnings()

        if argv:
            payload = argv
        else:
            payload = None

        # Form the url by combining the base with the passed in part
        url = "%s/%s" % (self.URL, restURI)

        r = requests.post(url, headers=self.headers, verify=False,
                          auth=HTTPBasicAuth(self.user, self.pw), json=payload,
                          )
        if r.status_code != 200:
            response = r.text
            response = json.loads(response)
            print(response)
            raise Exception
        else:
            response = r.text
            response = json.loads(response)
            #print(response)

    def alerts(self):
        url = ('api/types/alert/instances?fields=id,timestamp,severity,component,messageId,message,descriptionId,description,resolutionId,resolution')
        key = "alert"
        response = self.doGet(url)
        finalList = []
        for alert in response["entries"]:
            alert["content"]["array"] = self.clusterCI
            alert["content"]["key"] = key
            alert["content"]["time_created"] = alert["content"]["timestamp"]
            del alert["content"]["timestamp"]
            if "ALRT_SW_UPGRADE_AVAILABLE_WARNING" not in alert["content"]["descriptionId"]:
                finalList.append(alert["content"])

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)
        self.logout()

    def drives(self):
        url = ('api/types/disk/instances?fields=health,id,name,slotNumber')
        key = "drives"
        response = self.doGet(url)
        finalList = []
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key
            # item["content"]["operationalStatus"] = item["content"]["operationalStatus"][0]
            item["content"]["health_descr"] = item["content"]["health"]["descriptions"][0]
            item["content"]["health_descr_status"] = item["content"]["health"]["descriptionIds"][0]
            item["content"]["health_ID"] = item["content"]["health"]["value"]

            del item["content"]["health"]

            finalList.append(item["content"])

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)
        self.logout()

    def capacity(self):
        # Cluster capacity = by pool ( customer name )
        # Pool capacity = by Tier (constant names)->( extreme performance , capacity, performance )
        url = ('api/types/pool/instances?fields=name,health,sizeFree,sizeTotal,sizeUsed,tiers,snapSizeUsed,raidType')
        key = "cluster_capacity"
        response = self.doGet(url)
        finalList = []
        for pool in response["entries"]:
            pool["content"]["health"]["descriptionIds"] = pool["content"]["health"]["descriptionIds"][0]
            pool["content"]["health"]["descriptions"] = pool["content"]["health"]["descriptions"][0]
            pool["content"]["key"] = key
            pool["content"]["array"] = self.clusterCI
            poolTier = pool["content"]["name"]

            for i in pool["content"]["tiers"]:
                # Each Pool has 3 Tiers. Create name as identifier of Pool Tier
                poolTierName = poolTier + "_" + i["name"]

                row = {"array": self.clusterCI, "poolName": poolTierName, "poolSizeFree": i["sizeFree"],
                       "poolSizeTotal": i["sizeTotal"],
                       "poolSizeUsed": i["sizeUsed"], "key": "pool_capacity"}
                finalList.append(row)

            del pool["content"]["tiers"]
            finalList.append(pool["content"])

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)
        self.logout()

    def diskGroup(self):
        url = ('api/types/diskGroup/instances?fields=advertisedSize,diskSize,diskTechnology,emcPartNumber,hotSparePolicyStatus,id,minHotSpareCandidates,name,rpm,speed,tierType,totalDisks,unconfiguredDisks')
        key = "health"
        response = self.doGet(url)
        finalList = []
        for alert in response["entries"]:
            alert["content"]["array"] = self.clusterCI
            alert["content"]["key"] = key
            finalList.append(alert["content"])

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)
        self.logout()

    def hardware(self):

        allhardware = []

        ################### battery ###################
        batteryList = []
        url = ('api/types/battery/instances?fields=health,needsReplacement,parent,slotNumber,name')
        key = "battery"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key,
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   "slotNumber": item["content"]["slotNumber"]
                   }

            batteryList.append(row)

        # b = json.dumps(batteryList, indent=4, sort_keys=True)
        # print(b)

        ################### dpe ###################
        dpeList = []
        url = ('api/types/dpe/instances?fields=id,health,needsReplacement,parent,slotNumber,name,avgTemperature')
        key = "dpe"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "avgTemperature": item["content"]["avgTemperature"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   "slotNumber": item["content"]["slotNumber"]
                   }

            dpeList.append(row)

        # b = json.dumps(dpeList, indent=4, sort_keys=True)
        # print(b)

        ################### fan ###################
        fanList = []
        url = ('api/types/fan/instances?fields=health,needsReplacement,parent,slotNumber,id')
        key = "fan"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "id": item["content"]["id"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   "slotNumber": item["content"]["slotNumber"]
                   }
            #
            fanList.append(row)

        # b = json.dumps(fanList, indent=4, sort_keys=True)
        # print(b)

        ################### memoryModule ###################
        memoryModuleList = []
        url = ('api/types/memoryModule/instances?fields=health,needsReplacement,parent,slotNumber,name,id')
        key = "memoryModule"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "id": item["content"]["id"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   "slotNumber": item["content"]["slotNumber"]
                   }

            memoryModuleList.append(row)

        # b = json.dumps(memoryModuleList, indent=4, sort_keys=True)
        # print(b)

        ################### PSU ###################
        psuList = []
        url = ('api/types/powerSupply/instances?fields=health,needsReplacement,parent,slotNumber,name,id')
        key = "psu"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "id": item["content"]["id"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   "slotNumber": item["content"]["slotNumber"]
                   }

            psuList.append(row)

        # b = json.dumps(psuList, indent=4, sort_keys=True)
        # print(b)

        ################### sasPort ###################
        sasPortList = []
        url = ('api/types/sasPort/instances?fields=health,needsReplacement,parent,name,id,port')
        key = "sasPort"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "id": item["content"]["id"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   "port": item["content"]["port"]
                   }

            if "NOT_IN_USE" not in row["descriptionIds"]:
                sasPortList.append(row)

        # b = json.dumps(sasPortList, indent=4, sort_keys=True)
        # print(b)


        ################### ssd ###################
        ssdList = []
        url = ('api/types/ssd/instances?fields=health,needsReplacement,parent,slotNumber,name,id')
        key = "ssd"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "id": item["content"]["id"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   }

            ssdList.append(row)

        # b = json.dumps(ssdList, indent=4, sort_keys=True)
        # print(b)

        ################### SP (StorageProcessor ###################
        spList = []
        url = ('api/types/storageProcessor/instances?fields=health,needsReplacement,parent,slotNumber,name,id')
        key = "storageProcessor"
        response = self.doGet(url)
        for item in response["entries"]:
            item["content"]["array"] = self.clusterCI
            item["content"]["key"] = key

            row = {"array": self.clusterCI, "key": key, "id": item["content"]["id"],
                   "descriptionIds": item["content"]["health"]["descriptionIds"][0],
                   "descriptions": item["content"]["health"]["descriptions"][0],
                   "healthID": item["content"]["health"]["value"],
                   "needsReplacement": item["content"]["needsReplacement"],
                   "parentId": item["content"]["parent"]["id"],
                   "parentResource": item["content"]["parent"]["resource"],
                   }

            spList.append(row)

        # b = json.dumps(spList, indent=4, sort_keys=True)
        # print(b)

##########  Merge all hardware metrics to single list ##################
        allhardware.append(batteryList)
        allhardware.append(dpeList)
        allhardware.append(fanList)
        allhardware.append(memoryModuleList)
        allhardware.append(psuList)
        allhardware.append(sasPortList)
        allhardware.append(ssdList)
        allhardware.append(spList)

        b = json.dumps(allhardware, indent=4, sort_keys=True)
        print(b)

        self.logout()

    def lun(self):
        url = ('api/types/lun/instances?fields=health,name,id,sizeAllocated,sizeTotal,wwn,pool')
        key = "lun"
        response = self.doGet(url)
        finalList = []
        # print("LUN \t Allocated \t Total \t Free \t WWN \t Pool ")
        for lun in response["entries"]:
            lun["content"]["array"] = self.clusterCI
            lun["content"]["key"] = key

            row = {"array": self.clusterCI,
                   "descriptionIds": lun["content"]["health"]["descriptionIds"][0],
                   "descriptions": lun["content"]["health"]["descriptions"][0],
                   "healthID": lun["content"]["health"]["value"],
                   "id": lun["content"]["id"],
                   "key": lun["content"]["key"],
                   "name": lun["content"]["name"],
                   "sizeAllocated": lun["content"]["sizeAllocated"],
                   "sizeTotal": lun["content"]["sizeTotal"],
                   "sizeFree": int(lun["content"]["sizeTotal"]-lun["content"]["sizeAllocated"]),
                   "wwn": lun["content"]["wwn"].replace(":", ""),
                   "pool": lun["content"]["pool"],
                   #"sizeUsed": lun["content"]["sizeUsed"]
                   }
            finalList.append(row)

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)

        self.logout()

    def logout(self):
        url = ('api/types/loginSessionInfo/action/logout')
        self.doPost(url)


clusterCI = str(sys.argv[1])
user = str(sys.argv[2])
ip = str(sys.argv[3])
port = str(sys.argv[4])
ssl = str(sys.argv[5])
functione = str(sys.argv[6])
secre = str(sys.argv[7])


try:

    array_name = VNX(user, secre, ip, port, ssl, clusterCI)
    if "alerts" in functione:
        array_name.alerts()
    if "capacity" in functione:
        array_name.capacity()
    if "diskGroup" in functione:
        array_name.diskGroup()
    if "drives" in functione:
        array_name.drives()
    if "hardware" in functione:
        array_name.hardware()
    if "lun" in functione:
        array_name.lun()

    # else:
    #     print("No such function / ERROR")
finally:
    pass
