import time
import json
import PyU4V
import sys
from requests.auth import HTTPBasicAuth
import requests


class VMAX():
    def __init__(self, alias, serial, user, ip, port, ssl, secre):
        self.user = user
        self.pw = secre
        self.ip = ip
        self.port = port
        self.ssl = False
        self.arrayid = serial
        self.name = alias
        # Used by simple GET
        self.headers = {'Content-Type': 'application/json'}
        self.URL = 'https://%s:8443/univmax/restapi' % ip

    def connect(self):
        return (PyU4V.U4VConn(
            username=self.user, password=self.pw, server_ip=self.ip,
            port=self.port, verify=self.ssl, array_id=self.arrayid))

    def performance(self):
        conn = self.connect()
        # Current time
        current_time = int(time.time()) * 1000
        # Set start time to last 5 minutes
        start_time = current_time - (60 * 5 * 1000)
        # Get metrics raw
        a = conn.performance.get_array_metrics(start_time, current_time)
        # print(a)
        finalList = []
        # IF perf data available
        if a["perf_data"]:
            # print(a["perf_data"])

            # Select sub dictionary key
            perf = (a["perf_data"])
            if perf:
                del perf[0]["timestamp"]
            # Extract list to dict
            perf_list = perf[0]
            # Add array Id to dict
            perf_list["array"] = self.arrayid
            #perf_list["name"] = self.name
            perf_list["name"] = self.name
            perf_list["collection_error"] = 0

            # Add dict to new list
            finalList = [perf_list]
        else:
            error = {"symmetrixID": self.arrayid, "reporting_level": "array",
                     "error": "Performance data unavailable", "collection_error": 1}
            finalList.append(error)

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)

    def get_directors(self):
        conn = self.connect()
        directors = conn.provisioning.get_director_list()
        finalList = []
        for director in directors:
            a = conn.provisioning.get_director(director)
            a["array"] = self.arrayid
            a["name"] = self.name
            a["collection_error"] = 0

            if a["availability"] == "Online":
                a["status_code"] = 0
            elif a["availability"] == "Offline":
                a["status_code"] = 1
            else:
                a["status_code"] = 99
            finalList.append(a)

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)

    def capacity(self):
        conn = self.connect()
        srps = conn.provisioning.get_srp_list()

        for item in srps:
            srp = conn.provisioning.get_srp(item)
            srp["array"] = self.arrayid
            srp["name"] = self.name
            srp["collection_error"] = 0
            row = {"array": srp["array"], "vp_overall_ratio_to_one": srp["vp_overall_ratio_to_one"],
                   "overall_efficiency_ratio_to_one": srp["overall_efficiency_ratio_to_one"],
                   "vp_saved_percent": srp["vp_saved_percent"], "name": srp["name"],
                   "total_allocated_cap_gb": srp["total_allocated_cap_gb"],
                   "total_subscribed_cap_gb": srp["total_subscribed_cap_gb"],
                   "total_usable_cap_gb": srp["total_usable_cap_gb"],
                   "effective_used_capacity_percent": srp["effective_used_capacity_percent"]}
            finalList = [row]
            b = json.dumps(finalList, indent=4, sort_keys=True)
            print(b)

    def doGet(self, restURInotfull, payload=None):
        restURI = restURInotfull % self.arrayid

        # Form the url by combining the base with the passed in part
        url = "%s/%s" % (self.URL, restURI)

        r = requests.get(url, headers=self.headers, verify=False,
                         auth=HTTPBasicAuth(self.user, self.pw),
                         params=json.dumps(payload))

        # Convert the JSON data into a python dict
        data = r.json()
        return data

    def get_alert(self):
        url = ('90/system/symmetrix/%s/alert')
        response = array_name.doGet(url)
        finalList = []
        count = 0
        for alertID in response:
            for alert in response[alertID]:
                url_id = url + "/" + alert
                # do GET for each alert ID
                response_alerts = self.doGet(url_id)
                # Add array id to each alert
                response_alerts["array"] = self.arrayid
                response_alerts["collection_error"] = 0
                # Filter alert by NEW status (unacknowledged) and above informational
                if response_alerts["state"] == 'NEW':
                    # print(response_alerts)
                    count += 1
                    #print(count)
                    finalList.append(response_alerts)
                    # Read MAX NEW 20 alerts then break the loop (reduce query amount)
                    # Create alert if more than 20 NEW not ack alerts available. Due to performance. Each alert requires individual GET.
                    if count > 20:
                        break

        b = json.dumps(finalList, indent=4, sort_keys=True)
        print(b)

    def get_health(self):
        url = ('90/system/symmetrix/%s/health')
        response = array_name.doGet(url)
        newDict = {}
        # Add json response to list

        for entry in response:
            out = response[entry]

            for x in range(len(out)):
                a = out[x]
                a["array"] = self.arrayid
                a["errors"] = {}
                a["collection_error"] = 0

                b = a["instance_metrics"]
                a.pop("instance_metrics")

                for f in range(len(b)):
                    health_score_instance_metric = b[f]
                    instance_name = health_score_instance_metric["health_score_instance_metric"]

                    for instance_message in instance_name:
                        newDict.update(instance_message)
                        a["errors"] = newDict

        b = json.dumps(response, indent=4, sort_keys=True)
        print(b)


functione = str(sys.argv[1])
alias = str(sys.argv[2])
serial = str(sys.argv[3])
user = str(sys.argv[4])
ip = str(sys.argv[5])
port = str(sys.argv[6])
secre = str(sys.argv[7])
ssl = "false"
#

try:
    array_name = VMAX(alias, serial, user, ip, port, ssl, secre )
    if "performance" in functione:
        array_name.performance()
    elif "get_directors" in functione:
        array_name.get_directors()
    elif "volume_utilization" in functione:
        array_name.volume_utilization()
    elif "capacity" in functione:
        array_name.capacity()
    elif "alert" in functione:
        array_name.get_alert()
    elif "get_health" in functione:
        array_name.get_health()
    else:
        print("No such function / ERROR")
finally:
    pass
