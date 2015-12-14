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
            for i, e in enumerate(resp.json()):
                objlist.append(Subnet(e["mask"], e["gw"], e["siteId"]))
            return objlist

class Subnet:
    def __init__(self, mask, gw, site_id):
        self.mask = mask
        self.gw = gw
        self.site_id = site_id

clnt = Client()
data = clnt.Subnet().All()

for obj in data:
    print obj.mask, obj.gw, obj.site_id
