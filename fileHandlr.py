#this file handles all functions for files
#including excel files and txt files

from openpyxl import Workbook
from openpyxl import load_workbook
from openpyxl.styles import Font, PatternFill

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

#colors the cells of the sheet according to:
#blue < 10 hours
#green < 20 hours
#purple < 35 hours
#yellow < 50 hours
#orange < 80 hours
#red >= 80 hours
def colorCoder(wb, wsnum):
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

    #check if there is data to be color coded
    #if so then perform color coding
    if str(ws['A1'].value) == 'None' or str(ws['A2'].value) == 'None':
        print('Nothing there!')
        return
    else:
        #color defaults
        blue = PatternFill(fill_type = 'solid', start_color = '87CEEB', end_color='87CEEB')
        green = PatternFill(fill_type = 'solid', start_color = '6be575', end_color='6be575')
        purple = PatternFill(fill_type = 'solid', start_color = 'B589ed', end_color='B589ed')
        yellow = PatternFill(fill_type = 'solid', start_color = 'E6ed89', end_color='E6ed89')
        orange = PatternFill(fill_type = 'solid', start_color = 'F1b521', end_color='F1b521')
        red = PatternFill(fill_type = 'solid', start_color = 'F36767', end_color='F36767')
        print('Color coding begins now!')

        #obtain max number of rows within sheet
        i = 1
        size = 'A' + str(i)
        while (not (str(ws[size].value) == 'None')):
            i += 1
            size = 'A' + str(i)

        #begin coloring the data starting from row 2, column 2
        #i.e. starting from cell 'B2'
        for j in range(2,i):
            for k in range(2,4):
                curr = ''
                if k == 2:
                    curr = 'B'
                elif k == 3:
                    curr = 'C'
                
                #append the number to the letter as a string
                curr = curr + str(j)
                
                #if the data is not empty
                if (not ((str(ws[curr].value) == '--') or (str(ws[curr].value) == 'no data'))):

                    val = ws[curr].value

                    #determine if there is a '1/2'
                    #if not then find the location of the ' ' and make a substring from beginning of string to it
                    #if there is a '1/2' then find the location of it, and make a substring from beginning of string to it
                    #the correction will be performed in the next section
                    if (-1 == (str(ws[curr].value).find('½'))):
                        val = str(ws[curr].value)[:str(ws[curr].value).find(' ')]
                    elif (not (-1 == (str(ws[curr].value).find('½')))):
                        val = str(ws[curr].value)[:str(ws[curr].value).find('½')]
                    else:
                        val = str(ws[curr].value)

                    #in extreme cases, val may not hold useful data
                    #instead of coloring it, we skip it and notify the user of such an occurrence
                    if (len(ws[curr].value) > 6):
                        if onlyNum(val) == 1:

                            #turn val into a number (integer)
                            val = int(val)

                            #color code the cell
                            if val > 79:
                                ws[curr].fill = red
                            elif 80 > val and val > 50:
                                ws[curr].fill = orange
                            elif 51 > val and val > 35:
                                ws[curr].fill = yellow
                            elif 36 > val and val > 20:
                                ws[curr].fill = purple
                            elif 21 > val and val > 10:
                                ws[curr].fill = green
                            elif 11 > val:
                                ws[curr].fill = blue
                        else:
                            print('Error with data coloring - not numeric values')
                            print(curr)
                    else:
                        print('Error with data coloring - length is insufficient')
                        print(curr)
        print('Done color coding!') 
        #save work at the end
        wb.save(filename = 'gamelist_wb.xlsx')  
                    
#checks if a string contains only numerical characters
#i.e. 0-9                
def onlyNum(string):
    valid = 1
    for i in range(len(string)):
        for j in range (10):
            if string[i] == str(j):
                valid = 1
                break
            else:
                valid = 0
        if not (valid == 0):
            break
    return valid

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

#colorCoder(wbChecker(), 0)
