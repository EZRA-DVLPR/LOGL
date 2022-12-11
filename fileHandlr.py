# this file handles all functions for files
# including excel files and txt files

from openpyxl import Workbook
from openpyxl import load_workbook

# creates/opens the file to be written to
def wbChecker():

    # Try to open the workbook
    # else create it
    try:
        wb = load_workbook('test_wb.xlsx')
        ws = wb["Main Sheet"]
        print("Excel Sheet opened!")
        
    except:
        # Create a workbook
        wb = Workbook()

        # Change name of the initial worksheet
        ws = wb.active
        ws.title = "Main Sheet"

        # Make a new file 
        wb.save(filename = 'text_wb.xlsx')
        print("Excel Sheet created!")

    finally:
        print("Now writing to Excel Sheet...")
        #call another function to handle data input

# looks at the workbook and takes the data from the web and inputs it into the excel sheet
def addToWB(wb):
    #add entries to the workbook
    print("")

'''
# Once done print the names of all sheets of excel file
print(wb.sheetnames)

c = ws['A4']

#use CELL.value for obtaining the value within the cell
# c.value = NONE

#use this for single cell writing
ws['A4'] = 4
# c.value = 4

#use this for iteration
d = ws.cell(row = 3, column = 2,value=2)
wb.save(filename = 'text_wb.xlsx')

# opens a file for reading
with open('file_path') as f:
    f.read()
'''