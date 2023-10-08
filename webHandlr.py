#this file will handle all web connections and functions related to the web

from bs4 import BeautifulSoup
import requests

#initiates web scraping for entire gamelist
#returns the web-scraped data in an array:
#   [[GAME NAME, [MAIN STORY, COMPLETIONIST]], [GAME NAME, ['no data']], ...]
def beginWebScrape(gamelist):
    print('Webscraping has begun!')
    print()

    data_all = []

    #for each game in list, perform search and update data_all with info
    for i in range (len(gamelist)):
        print('Working on game ' + str(i+1) + ' of ' + str(len(gamelist)))
        print()
        
        #add a new entry in data_all with the title of the game
        data_all.append([gamelist[i]])

        #generates query based on game name
        #then extracts hltb URL via the generated query
        gamelist[i] = queryGenerator(gamelist[i])
        hltbURL = googleSearch(gamelist[i])
        
        #if valid google connection was established then appends extracted data
        if (not (hltbURL == 'ERROR ON GOOGLE CONNECTION')):
            #extracts stats on game and appends them to data_all[i]

            #hours: [main, completionist]
            # OR
            #hours: ['no data']
            data_all[i].append(hltbExtract(hltbURL))

    print('Webscraping process has ended!')
    print()
    return data_all

#gets the url of the hltb query for the specified game
#returns the hltb url as a string
def googleSearch(gamename):

    #make the google search with the given name
    googleurl = 'https://www.google.com/search?q=hltb+' + gamename
    try:
        googlereq = requests.get(googleurl)
        if (googlereq.status_code == '200'):
            raise InterruptedError('No valid connection')

    except:
        #error with connection to Google
        print('ERROR CONNECTING FOR ' + gamename)
        print()
        return 'ERROR ON GOOGLE CONNECTION'

    else:
        #make a soup then search for the first hit on google
        soup = BeautifulSoup(googlereq.text, 'lxml')
        match = soup.find('div', class_='egMi0 kCrYT')

        #obtain and return first link that matches
        link = match.find('a')
        newlink = link.get('href')
        newurl = newlink[7:newlink.find('&')]
        return newurl

#connects to hltb and extracts the gamedata
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
        #error connecting to HLTB
        print('ERROR CONNECTING to HLTB ')
        print()
        return 'ERROR ON HLTB CONNECTION'

    else:
        #make a soup and connect to HLTB website
        soup = BeautifulSoup(HLTBreq.text, 'lxml')

        #hours will be the extracted data if it exists
        #otherwise it will be ['no data']
        hours = []

        #should be the parent div of each section (Main Story, Main Story +, ...)
        #this is subject to change by the developers at HLTB and requires attention for each run
        match = soup.find('div', class_='GameStats_game_times__KHrRY')

        if (not match):
            #this game is not a standard single-player experience, and thus no data will not be grabbed
            hours = ['no data']
            print('no matching div, check manually')
        else:
            #only grabs the numbers and not the hours

            try:
                #check for main story
                if (len(match.find_all('h4')) > 0) and (match.find_all('h4')[0].text == 'Main Story'):
                    hours.append(float(hourDecConvert(match.find_all('h5')[0].text[:match.find_all('h5')[0].text.find(' ')])))
            
                #check for completionist
                if (len(match.find_all('h4')) > 2) and (match.find_all('h4')[2].text == 'Completionist'):
                    hours.append(float(hourDecConvert(match.find_all('h5')[2].text[:match.find_all('h5')[2].text.find(' ')])))

                #if no such fields exist then there is no relevant data
                if hours == []:
                    hours = ['no data']
                    print('no hours information available despite match, check manually')

            except:
                hours = ['no data']
                print('no info, check manually')
        
        return hours

#generates a query to use in google search based on given game name
#returns a string that is query-safe with no common special characters
def queryGenerator(gamename):
    #first replace all commonly used non-alphanumeric characters that have different formats
    #into a readable format for query
    gamename = gamename.replace('\'','%27').replace('&','%26').replace(':','%3A').replace(',','%2C').replace('+','%2B')

    #turns any spaces into '+'
    gamename = gamename.replace(' ', '+')
    
    return gamename

#returns given string with decimal point instead of fraction
def hourDecConvert(gamehours):
    gamehours = gamehours.replace('½','.5')
    return gamehours
