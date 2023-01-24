#this file will handle all web connections and functions related to the web

from bs4 import BeautifulSoup
import requests

#initiates the entire web scraping for the entire gameslist
#returns the web-scraped data in an array
def beginWebScrape(gamelist):
    print('Webscraping has begun!')
    print()

    #data_all will hold the titles and the hours array
    #hours: main, completionist
    # OR
    #hours: 'no data'
    data_all = []

    #iterate through given list
    for i in range (len(gamelist)):
        print('Working on game ' + str(i+1) + ' of ' + str(len(gamelist)))
        print()
        
        #add a new entry in data_all with the title of the game
        data_all.append([gamelist[i]])

        #fix name of game
        gamelist[i] = fixGamename(gamelist[i])
        
        #extracts hltbURL
        hltbURL = googleSearch(gamelist[i])
        
        #checks validity of connection and only proceeds with extraction if google connection is good
        if (not (hltbURL == 'ERROR ON GOOGLE CONNECTION')):
            #extracts stats on game and appends them to data_all[i]
            data_all[i].append(hltbExtract(hltbURL))

    print('Webscraping process has ended!')
    print()
    return data_all

#gets the url of the hltb for the specified game
#returns the url as a string
def googleSearch(gamename):

    #make the google search with the given name
    googleurl = 'https://www.google.com/search?q=hltb+' + gamename
    try:
        googlereq = requests.get(googleurl)
        if (googlereq.status_code == '200'):
            raise InterruptedError('No valid connection')

    except:
        #if error with connection then say so
        print('ERROR CONNECTING FOR ' + gamename)
        print()
        return 'ERROR ON GOOGLE CONNECTION'

    else:
        #make a soup then search for the first hit on google
        soup = BeautifulSoup(googlereq.text, 'lxml')
        match = soup.find('div', class_='egMi0 kCrYT')

        #get the link of the first hit and then return that link
        link = match.find('a')
        newlink = link.get('href')
        newurl = newlink[7:newlink.find('&')]
        return newurl

#connects to hltb and gets the gamedata
#default timeout to 2 seconds
#returns the extracted data (hours) for the game in an array
def hltbExtract(url):
    #assign header for connecting to hltb.com
    header = {"User-Agent" : "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15"}

    try:
        HLTBreq = requests.get(url, headers=header, timeout=2)
        if (HLTBreq.status_code == '200'):
            raise InterruptedError('No valid connection')

    except:
        #if error with connection then say so
        print('ERROR CONNECTING to HLTB ')
        print()
        return 'ERROR ON HLTB CONNECTION'

    else:
        #make a soup and then connect to HLTB website
        soup = BeautifulSoup(HLTBreq.text, 'lxml')

        #hours will be the extracted data if it exists
        #otherwise it will be ['no data']
        hours = []

        match = soup.find('div', class_='GameStats_game_times__5LFEc')

        if (not match):
            #this game is not a standard single-player experience, and thus no data will not be grabbed
            hours = ['no data']
            print('no info, check manually')
            print()
        else:
            #only grabs the numbers and not the hours


            #check for main story
            if (len(match.find_all('h4')) > 0) and (match.find_all('h4')[0].text == 'Main Story'):
                hours.append(fixGameHours(match.find_all('h5')[0].text[:match.find_all('h5')[0].text.find(' ')]))
                

            #check for completionist
            if (len(match.find_all('h4')) > 2) and (match.find_all('h4')[2].text == 'Completionist'):
                hours.append(fixGameHours(match.find_all('h5')[2].text[:match.find_all('h5')[2].text.find(' ')]))


            #if no such fields exist then there is no relevant data
            if hours == []:
                hours = ['no data']
                print('no info, check manually')
            
            #print(hours)
            #print()
        
        return hours

#corrects the gamename for special characters
#returns a string that is properly formatted for google searches
def fixGamename(gamename):
    #first replace all non-alphanumeric characters that have different formats
    gamename = gamename.replace('\'','%27')

    gamename = gamename.replace('&','%26')
    
    gamename = gamename.replace(':','%3A')

    gamename = gamename.replace(',','%2C')

    gamename = gamename.replace('+','%2B')

    #turns any spaces into '+'
    gamename = gamename.replace(' ', '+')
    
    return gamename

#corrects the hours if they have the single character containing '1/2' by replacing with '.5' instead
#returns a properly formatted string
def fixGameHours(gamehours):
    gamehours = gamehours.replace('½','.5')
    return gamehours

