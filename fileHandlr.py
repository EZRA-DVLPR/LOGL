#this file handles all functions for files
#including excel files and txt files

from openpyxl import Workbook
from openpyxl import load_workbook
from openpyxl.styles import Font

#creates/opens the file to be written to
def wbChecker():

    # Try to open the workbook
    # else create it
    try:
        wb = load_workbook('gamelist_wb.xlsx')

        #check if Wishlist worksheet exists and if not then create worksheet named wishlist
        if (not (wb.sheetnames[1])):
            wb.create_sheet("Wishlist")

        #set active worksheet to owned games
        wb.active = wb['Owned Games']

        print("Excel Sheet opened!")
        
    except:
        # Create a workbook
        wb = Workbook()

        # Change name of the initial worksheet
        ws = wb.active
        ws.title = "Owned Games"

        wb.create_sheet('Wishlist')

        wb.active = wb['Owned Games']

        # Make a new file 
        wb.save(filename = 'gamelist_wb.xlsx')
        print("Excel Sheet created!")
        print()

    finally:
        print("Now writing to Excel Sheet...")
        return wb

#looks at the workbook given, and inputs the data into the workbook
#next parameter is an integer: 0 = add to first worksheet, 1 = add to second worksheet
def addToWB(wb, data_all, wsnum):
    #open wb and 'Owned Games' from ws

    #switches wsnum for the selected worksheet for the data to be input into
    if (wsnum == 0):
        wb.active = wb['Owned Games']
        ws = wb['Owned Games']
    elif (wsnum == 1):
        wb.active = wb['Wishlist']
        ws = wb['Wishlist']
    else:
        print('invalid worksheet')
        return -1

    #adjust first 3 columns to be wider for legibility
    ws.column_dimensions['A'].width = 20
    ws.column_dimensions['B'].width = 20
    ws.column_dimensions['C'].width = 20

    #list the column names in bold
    ws.cell(row = 1, column = 1, value = 'Game Name')
    ws['A1'].font = Font(bold=True)
    ws.cell(row = 1, column = 2, value = 'Main Story')
    ws['B1'].font = Font(bold=True)
    ws.cell(row = 1, column = 3, value = 'Completionist')
    ws['C1'].font = Font(bold=True)

    #adjust row height for column names
    ws.row_dimensions[1].height = 20

    #adjust all rows height to be 20 while inputting the data into the row
    for i in range(len(data_all)):
        ws.row_dimensions[i + 2].height = 20
        
        #add name of game
        ws.cell(row = i+2, column = 1, value = data_all[i][0])

        #add data from game if it exists
        #otherwise add 'no data'

        #data_all[index][0 = gamename, 1 = array of hours][which value of hours we want: 0 = main story, 1 = completionist]
        if data_all[i][1][0] == 'no data':
            ws.cell(row = i+2, column = 2, value = data_all[i][1][0])
            ws.cell(row = i+2, column = 3, value = data_all[i][1][0])
        else:
            ws.cell(row = i+2, column = 2, value = data_all[i][1][0])
            ws.cell(row = i+2, column = 3, value = data_all[i][1][1])

    #save work at the end
    wb.save(filename = 'gamelist_wb.xlsx')

    print("Done!")

#reads a txt file and inputs the data from the txt file into a single array
def fileReader():
    lines = []
    try:
        with open('games.txt') as f:
            print('Reading the txt file!')
            for line in f:
                lines.append(line.strip())
            print('Done reading the txt file!')
            return lines
    
    except:
        print('ERROR READING TXT FILE')
        return