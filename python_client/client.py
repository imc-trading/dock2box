# -*- coding: utf-8 -*-
'''
Dock2Box REST Client
'''

import json
import requests

class Client:
    url = "http://localhost:8080/v1"

    def __init__(self, url = "http://localhost:8080/v1"):
        self.url = url

    class Subnet:
        def All(self):
            resp = requests.get(Client.url + "/subnets?envelope=false")
            objlist = []
            for e in resp.json():
                objlist.append(Subnet(e["subnet"], e["mask"], e["gw"], e["siteId"]))
            return objlist

        def Get(self, name):
            resp = requests.get(Client.url + "/subnets/{0}?envelope=false".format(name))
            e = resp.json()
            return Subnet(e["subnet"], e["mask"], e["gw"], e["siteId"])

class Subnet:
    def __init__(self, subnet, mask, gw, site_id):
        self.subnet = subnet
        self.mask = mask
        self.gw = gw
        self.site_id = site_id

clnt = Client()
data = clnt.Subnet().All()

for obj in data:
    print obj.mask, obj.gw, obj.site_id

obj = clnt.Subnet().Get("192.168.0.0/24")
print obj.mask, obj.gw, obj.site_id
