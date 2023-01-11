# this file will handle all web connections and functions related to the web

import requests
import json

url = 'https://howlongtobeat.com/'
header = {"User-Agent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"}
cookie = {"cto_bidid" : "3fGYkV9kWkxDSWYwTyUyRmlUVUY1a2dTbTdsaG1ob3o0eHEzcXJNT0lydlc4V1dqYVpnMDJ4SzUwdndaVzJkSWlBRFJlbHhsMzM0aFZDNWlpWm9IaEJ6dWxaMCUyRnclM0QlM0Q",
    "cto_bundle" : "bByyW19leUVocDlHZDJZeDZ3NVZ2UEM5YlZ0cXpyeU8yMVclMkJlUUNBUWliSE96QTRrUExMQlVtUEZ4RXI0RXpmODhkSm1mOUdnSHFOUWtGY3J1WnNtUTE5VG9ZMUpCR1RFc21NaXVKZVVucnFkT3JMNCUyQjVKRkZTWlZEc29Qa2c4VGNZOGI",
    "_ga" : "GA1.2.137075115.1670704387",
    "_gid" : "GA1.2.762652137.1671044747",
    "_pbjs_userid_consent_data" : "3524755945110770"}

query = {"q" : "pokemon"}

dataa = {"searchType":"games","searchTerms":["elden"],"searchPage":1,"size":20,"searchOptions":{"games":{"userId":0,"platform":"","sortCategory":"popular","rangeCategory":"main","rangeTime":{"min":"null","max":"null"},"gameplay":{"perspective":"","flow":"","genre":""},"rangeYear":{"min":"","max":""},"modifier":""},"users":{"sortCategory":"postcount"},"filter":"","sort":0,"randomizer":0}}

r = requests.get(url, headers=header, params=query)    #, cookies=cookie)

print(r.url)
print("\n\n")

print(r.status_code)
print("\n\n")

r = requests.get(url + '?q=ninja', headers=header)

print(r.status_code)
print("\n\n")

print(r.history)

with open('file.txt', 'w') as f:
    f.write(r.text)

#if cant get it to work then instead search google for results and pick the one with the 
# HLTB as the webpage