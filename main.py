import fileHandlr
import webHandlr

#read a txt file and insert it into an array (gamelist)
gamelist = fileHandlr.fileReader()

#begin web scraping with the gamelist
#save the returned data into data_all
data_all = webHandlr.beginWebScrape(gamelist)

#open the excel workbook and add it to the workbook
fileHandlr.addToWB(fileHandlr.wbChecker(), data_all, 0)

#color code the data
fileHandlr.colorCoder(fileHandlr.wbChecker(), 0)

