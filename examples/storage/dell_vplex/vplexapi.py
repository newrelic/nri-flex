import sys
import json
from requests.auth import HTTPBasicAuth
import requests
import urllib3


class VpLex:
    def __init__(self, ip, port, user, secret, arrayid):
        self.user = user
        self.pw = pw
        self.ip = ip
        self.port = port
        self.arrayid = arrayid
        self.headers = {'Content-Type': 'application/json;format=1;prettyprint=1'}
        self.URL = 'https://%s:%s' % (ip, port)

    def doGet(self, restURI, payload=None):
        urllib3.disable_warnings()
        # Form the url by combining the base with the passed in part
        url = "%s/%s" % (self.URL, restURI)

        try:
            r = requests.get(url, headers=self.headers, verify=False,
                             auth=HTTPBasicAuth(self.user, self.pw),
                             )
            if r.status_code == 200:
                return r.json()
            elif r.status_code == 401:
                print(r.status_code, ": Authentication has failed. The client did not provide the correct username and password")
                return None
            elif r.status_code == 404:
                print(r.status_code, ": Nonexistent URI. Context or command not found")
                return None
            elif r.status_code == 403:
                print(r.status_code , ": Forbidden for PUT request to change a read-only attribute")
                return None
        except requests.RequestException as ex:
            return ex

    def process(self, response, keyid, grep, splitid):
        list = []

        try:
            for item in response["response"]["context"]:
                row = {}
                all_row = {}

                for item_values in item["attributes"]:
                    key = item_values["name"]
                    value = item_values["value"]
                    rowtemp = {key: value}
                    row.update(rowtemp)

                all_row.update(row)
                all_row["array"] = self.arrayid
                all_row["key"] = keyid

                # IF key "operational-status" exist, insert additional key
                if "operational-status" in all_row:
                    if all_row["operational-status"] == "error":
                        all_row["operational-status-id"] = 1
                    elif all_row["operational-status"] == "failure":
                        all_row["operational-status-id"] = 1
                    elif all_row["operational-status"] == "ok":
                        all_row["operational-status-id"] = 0
                    elif all_row["operational-status"] == "online":
                        all_row["operational-status-id"] = 0
                    # If operational status not available (NULL), but health is OK
                    elif all_row["operational-status"] is None:
                        if all_row["health-state"] == "ok":
                            all_row["operational-status-id"] = 0
                    else:
                        all_row["operational-status-id"] = 99

                # If director ID is required and "yes" was passed to process function, grep splitid array element from
                # "parent" key. Ex: "parent": "/engines/engine-1-1/directors/director-1-1-B/hardware/sfps",
                if grep == "yes":
                    directorID = item["parent"].split("/")[splitid]
                    all_row["directorid"] = directorID



                list.append(all_row)

            final_list = [list]
            return final_list
        except TypeError as e:
            print(e, "<- No data to process")

    def psu(self):
        url = ('vplex/engines/*/power-supplies/*')
        key = "psu"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 2)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def ups(self):
        url = ('vplex/clusters/cluster-*/uninterruptible-power-supplies/ups-*-*')
        key = "ups"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 2)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def devices(self):
        url = ('vplex/clusters/cluster-*/devices/*')
        key = "devices"

        response = self.doGet(url)
        final_list = self.process(response, key, "no", 0)

        # Remove Bytes "B" in capacity value
        for entry in final_list[0]:
            entry["capacity"] = entry["capacity"].split("B")[0]

        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def engines(self):
        url = ('vplex/engines/*')
        key = "engines"

        response = self.doGet(url)
        final_list = self.process(response, key, "no", 0)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def alerts(self):
        url = ('vplex/alerts/triggered-alerts')
        key = "alerts"

        response = self.doGet(url)
        final_list = self.process(response, key, "no", 0)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def battery(self):
        url = ('vplex/engines/*/battery-backup-units/*')
        key = "battery"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 2)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def engine_fan(self):
        url = ('vplex/engines/*/fans/*')
        key = "engine_fan"

        response = self.doGet(url)
        final_list = self.process(response, key, "no", 0)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def director_fan(self):
        url = ('vplex/engines/*/directors/*/hardware/fan-modules/*')
        key = "director_fan"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 4)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def ports(self):
        url = ('vplex/engines/*/directors/*/hardware/ports/*')
        key = "directors_ports"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 4)

        # For "Ports" function remove list of port protocol
        for entry in final_list[0]:
            entry["protocols"] = entry["protocols"][0]

        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def iomodule(self):
        url = ('vplex/engines/*/directors/*/hardware/io-modules/*')
        key = "directory_iomodule"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 4)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def sfps(self):
        url = ('vplex/engines/*/directors/*/hardware/sfps/*')
        key = "director_sfps"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 4)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def dimms(self):
        url = ('vplex/engines/*/directors/*/hardware/dimms/*')
        key = "director_dimms"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 4)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)

    def mgmt_modules(self):
        url = ('vplex/engines/*/mgmt-modules/*')
        key = "mgmt_modules"

        response = self.doGet(url)
        final_list = self.process(response, key, "yes", 2)
        b = json.dumps(final_list, indent=4, sort_keys=True)
        print(b)


def main():
    functione = str(sys.argv[1])
    ip = str(sys.argv[2])
    port = str(sys.argv[3])
    user = str(sys.argv[4])
    arrayid = str(sys.argv[5])
    secret = str(sys.argv[6])

    try:
        array_name = VpLex(ip, port, user, secret, arrayid)
        if "psu" in functione:
            array_name.psu()
        if "ups" in functione:
            array_name.ups()
        if "devices" in functione:
            array_name.devices()
        if "engines" in functione:
            array_name.engines()
        if "alerts" in functione:
            array_name.alerts()
        if "battery" in functione:
            array_name.battery()
        # VERY SLOW
        if "director_fan" in functione:
            array_name.director_fan()
        if "engine_fan" in functione:
            array_name.engine_fan()
        if "ports" in functione:
            array_name.ports()
        if "iomodule" in functione:
            array_name.iomodule()
        if "sfps" in functione:
            array_name.sfps()
        if "dimms" in functione:
            array_name.dimms()
        if "mgmt_modules" in functione:
            array_name.mgmt_modules()
        #else:
        #    print("No such function in script / ERROR")
    finally:
        pass


if __name__ == '__main__':
     main()

