import purestorage #purestorage/rest-client/
import urllib3
import json
import sys
import base64
import os

os.environ['no_proxy'] = '*'

# This script polls Pure's REST API via py client: https://pure-storage-python-rest-client.readthedocs.io/en/stable/api.html
# Arguments specified in `pure_hosts.json`


class PureAPI:
    def __init__(self, IP, token, host, clusterCI):
        self.IP = IP
        self.token = token
        self.host = host
        self.clusterCI = clusterCI
        urllib3.disable_warnings()
        self.array = purestorage.FlashArray(self.IP,
                                            api_token=self.token)

    def add_keys(self, dictionary):
        for elements in dictionary:
            elements['array'] = self.host
            elements['clusterCI'] = self.clusterCI
        return dictionary

    def capacity(self):
        array_info = self.array.get(space="true")

        # Add "array" name key
        array_info = self.add_keys(array_info)

        # Make json pretty
        print(json.dumps(array_info, indent=4, sort_keys=True))

    def drives(self):
        # Read the list of drives
        drives = self.array.list_drives()

        # Convert list of arrays to json format
        jsonDrives = json.dumps(drives)
        parsed = json.loads(jsonDrives)

        # Add "array" name key
        parsed = self.add_keys(parsed)

        # Insert statusID
        for elements in parsed:
            if "healthy" in elements['status']:
                elements['statusID'] = 1
            elif "unused" in elements['status']:
                elements['statusID'] = 0
            elif "evacuat" in elements['status']:
                elements['statusID'] = 0
            else:
                elements['statusID'] = 2

        # Make json pretty
        a = json.dumps(parsed, indent=4, sort_keys=True)
        print(a)

    def int_stats(self):
        # list_hardware = self.array.get_network_interface(action="monitor", interface="CT1.FC8")
        list_interfaces = self.array.list_network_interfaces(action="monitor")
        jsonResponse = json.dumps(list_interfaces)
        parsed = json.loads(jsonResponse)

        # Add "array" name key
        parsed = self.add_keys(parsed)

        a = json.dumps(parsed, indent=4, sort_keys=True)
        print(a)

    def net_int(self):
        # Read the list of drives
        net_int = self.array.list_network_interfaces()

        # Convert list of arrays to json format
        jsonNet_int = json.dumps(net_int)
        parsed = json.loads(jsonNet_int)

        # Add "array" name key
        parsed = self.add_keys(parsed)

        # Make json pretty
        a = json.dumps(parsed, indent=4, sort_keys=True)
        print(a)

    def port(self):
        # Read the list of drives
        port = self.array.list_ports()

        # Convert list of arrays to json format
        jsonPort = json.dumps(port)
        parsed = json.loads(jsonPort)

        # Add "array" name key
        parsed = self.add_keys(parsed)

        # Make json pretty
        a = json.dumps(parsed, indent=4, sort_keys=True)
        print(a)

    def list_messages(self):
        # Read any active/open messages/events
        list_messages = self.array.list_messages(open="true")

        # Convert list of messages to json format
        jsonlist_messages = json.dumps(list_messages)

        # Make json pretty
        parsed = json.loads(jsonlist_messages)

        # Add "array" name key
        parsed = self.add_keys(parsed)

        # Make json pretty
        a = json.dumps(parsed, indent=4, sort_keys=True)
        print(a)

    def list_hardware(self):
        # Read the list of drives
        list_hardware = self.array.list_hardware()

        # Convert list of arrays to json format
        jsonResponse = json.dumps(list_hardware)

        # Make json pretty
        parsed = json.loads(jsonResponse)

        # Add "array" name key
        parsed = self.add_keys(parsed)

        # Insert statusID key
        for elements in parsed:
            if "ok" in elements["status"]:
                elements["statusID"] = 0
            else:
                elements["statusID"] = 1

        # Make json pretty
        listFinal = [parsed]
        a = json.dumps(listFinal, indent=4, sort_keys=True)
        print(a)

    def perf(self):
        perf_latency = self.array.get(action="monitor")
        perf_array = self.array.get(action="monitor", latency="true")
        perf = perf_array + perf_latency

        # Add "array" name key
        perf = self.add_keys(perf)

        # Make json pretty
        print(json.dumps(perf, indent=4, sort_keys=True))

    def volumes(self):
        # Get volumes size and creation time
        list_volumes = self.array.list_volumes()

        # Convert list of messages to json format
        jsonlist_volumes = json.dumps(list_volumes)

        # Make json pretty
        parsed = json.loads(jsonlist_volumes)

        new_list = []
        # Add "array" name key
        for elements in parsed:
            elements['array'] = self.host
            elements['clusterCI'] = self.clusterCI

            row ={"array": elements['array'], "name": elements['name'], "size": elements['size'],
                  "created": elements['created'], "clusterCI": elements['clusterCI']}
            new_list.append(row)
        # Make json pretty
        a = json.dumps(new_list, indent=4, sort_keys=True)
        print(a)

    def vol_host_conn(self):
        # List all volumes
        volumes = self.array.list_volumes()
        for volume in volumes:
            volname = volume["name"]
            # Check host to volume connections
            volume_hosts = self.array.list_volume_shared_connections(volume=volname)
            if not volume_hosts:
                # Volume has no connections to hosts
                volume["host_connected"] = len(volume_hosts)

            else:
                # Volume has hosts connected
                volume["host_connected"] = len(volume_hosts)

        # Convert list of messages to json format
        jsonlist_lun = json.dumps(volumes)
        # Make json pretty
        parsed = json.loads(jsonlist_lun)
        # Add "array" name key
        parsed = self.add_keys(parsed)
        # Make json pretty
        a = json.dumps(parsed, indent=4, sort_keys=True)
        print(a)

    def certs(self):
        certs = self.array.list_certificates()
        jsonlist_lun = json.dumps(certs)
        parsed = json.loads(jsonlist_lun)
        parsed = self.add_keys(parsed)
        b = json.dumps(parsed, indent=4, sort_keys=True)
        print(b)


# system arguments input
IP = str(sys.argv[1])
token = str(sys.argv[2]) #base64 input
alias = str(sys.argv[3])
clusterCI = str(sys.argv[4])
func = str(sys.argv[5])

#decode
pa = base64.b64decode(token).decode("ascii", "ignore")

attempt = 1
while attempt:
    try:
        array_name = PureAPI(IP, pa, alias, clusterCI)
        if "drives" in func:
            array_name.drives()
            array_name.array.invalidate_cookie()
        elif "hardware" in func:
            array_name.list_hardware()
            array_name.array.invalidate_cookie()
        elif "messages" in func:
            array_name.list_messages()
            array_name.array.invalidate_cookie()
        elif "port" in func:
            array_name.port()
            array_name.array.invalidate_cookie()
        elif "interface" in func:
            array_name.net_int()
            array_name.array.invalidate_cookie()
        elif "capacity" in func:
            array_name.capacity()
            array_name.array.invalidate_cookie()
        elif "perf" in func:
            array_name.perf()
            array_name.array.invalidate_cookie()
        elif "volumes" in func:
            array_name.volumes()
            array_name.array.invalidate_cookie()
        elif "int_stats" in func:
            array_name.int_stats()
            array_name.array.invalidate_cookie()
        elif "vol_host_conn" in func:
            array_name.vol_host_conn()
            array_name.array.invalidate_cookie()
        elif "certs" in func:
            array_name.certs()
            array_name.array.invalidate_cookie()
        else:
            print("Incorrect function name")
        break
    except purestorage.purestorage.PureHTTPError as response:
        y = json.loads(response.text)
        print("Error occurred: ", (y['msg']))
        break
    except purestorage.purestorage.PureError as response:
        print("Bad hostname/IP: \n",  response)
        break
