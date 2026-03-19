import requests
import json
import sys
 
username = sys.argv[1]
password = sys.argv[2]
authorization = str(sys.argv[3])
b_url= str(sys.argv[4])

url=b_url+'/oauth/token'

headers={
    "Content-Type": "application/x-www-form-urlencoded",
    "Authorization": authorization
    }

payload= {'grant_type': 'password', 'username': username,'password': password }
x = requests.post(url,headers=headers,data=payload)

json_data = json.loads(x.text)

print(json.dumps({'access_token':json_data['access_token']}))