import fileHandlr
import webHandlr
'''
#read a txt file and insert it into an array (gamelist)
gamelist = fileHandlr.fileReader('games.txt')

#begin web scraping with the gamelist
#save the returned data into data_all
data_all = webHandlr.beginWebScrape(gamelist)

#open the excel workbook and add it to the workbook
fileHandlr.addToWB(fileHandlr.wbChecker(), 0, data_all)

#color code the data
fileHandlr.colorCoder(fileHandlr.wbChecker(), 0)

#add platform column and new games to list
fileHandlr.addNewCol(fileHandlr.wbChecker(), 0, 'newcol.txt')
fileHandlr.addNewRow(fileHandlr.wbChecker(), 0, 'newrow.txt')
'''
#sort with the new data that has been added
#alphabetically
#main story     -alphabetically, completionist, platform
#completionist  -alphabetically, main story, platform
#platform       -alphabetically, completionist, main story
fileHandlr.sorter(fileHandlr.wbChecker(), 0, 1, 2)

fileHandlr.colorCoder(fileHandlr.wbChecker(), 0)