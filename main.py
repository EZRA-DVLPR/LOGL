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

#properly format the entire sheet
fileHandlr.wsFormatter(fileHandlr.wbChecker(), 0)

#sort with the new data that has been added
fileHandlr.uncolorCoder(fileHandlr.wbChecker(), 0)

fileHandlr.sorter(fileHandlr.wbChecker(), 0, 1, 2)
