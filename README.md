# README

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

## Description

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

This serves as a beginner into web-scraping in python. It is not intended to be abused or malicious.

The program searches for the lengths of completion for games on the website 'www.howlongtobeat.com' and puts that information into an excel sheet.

There are additional functions that allow for sorting, color coding, appendnig column/row (s).

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

## What running this program does:

-----------------------------------------------------------------------------------------------------------------------------------------------------------------
The file to be run is 'main.py'.

It begins by reading the txt file named 'games.txt', which then gets searched on 'www.google.com' where it looks at the first matching link to the 'www.howlongtobeat.com' URL containing the desired game's information

After obtaining the URL from HLTB, it grabs the amount of hours needed for 'Main Story' and 'Completionist', if available.

With this information inputs it into an excel file titled 'gamelist_wb.xlsx'.

From there it will color-code the data.

Further functionality includes:
    new columns being added (`Platform`)
    new rows being added (new games)

Once this information is input, the data may be sorted

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

# LEGAL

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

## License  

- The **code and documentation** in this repository are licensed under the **MIT License**. See [LICENSE](LICENSE) for details.  
- The **project name and icon** are licensed under the **Creative Commons Attribution 4.0 International License**. See [LICENSE_ICON_NAME.md](LICENSE_ICON_NAME.md) for details.  

This is an independent project done by me with no further collaborators.

This project has no affiliation with Google, Bing, HowLongToBeat, nor Completionator.
-----------------------------------------------------------------------------------------------------------------------------------------------------------------
