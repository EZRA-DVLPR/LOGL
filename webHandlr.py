# this file will handle all web connections and functions related to the web

import requests

url = 'https://howlongtobeat.com'



r = requests.get(url)

print(r.text)
print("\n\n")

print(r.status_code)
print("\n\n")

