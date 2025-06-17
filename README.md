# README

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

```
__/\\\___________________/\\\\\__________/\\\\\\\\\\\\__/\\\_____________        
 _\/\\\_________________/\\\///\\\______/\\\//////////__\/\\\_____________       
  _\/\\\_______________/\\\/__\///\\\___/\\\_____________\/\\\_____________      
   _\/\\\______________/\\\______\//\\\_\/\\\____/\\\\\\\_\/\\\_____________     
    _\/\\\_____________\/\\\_______\/\\\_\/\\\___\/////\\\_\/\\\_____________    
     _\/\\\_____________\//\\\______/\\\__\/\\\_______\/\\\_\/\\\_____________   
      _\/\\\______________\///\\\__/\\\____\/\\\_______\/\\\_\/\\\_____________  
       _\/\\\\\\\\\\\\\\\____\///\\\\\/_____\//\\\\\\\\\\\\/__\/\\\\\\\\\\\\\\\_ 
        _\///////////////_______\/////________\////////////____\///////////////__
```
-----------------------------------------------------------------------------------------------------------------------------------------------------------------

LOGL 

AKA

Library
Of
Game
Lengths

# Description


This program sets out to hold a local database of games and associated time completion data.

I wanted to create a piece of software that is platform independent that allows you to keep track of how long it takes to complete a game in your library.
This should include any game from any library (Steam, GOG, etc.) as well as user input games.
The data should be retrieved from [HLTB](https://howlongtobeat.com/) and other websites that may offer the similar information. 

There are two main reasons for making this program:
1) Help get rid of the Gaming Backlog
2) Learning GoLang

## Features

- Customizable
	- Default Light/Dark Themes
    - Able to create any theme you like (See: [YouTube guide on customization](https://www.youtube.com/playlist?list=PL_gNvZlhoitBNANmcZFgoQpT1FjZiBs7I))
    - Search Source Selection
- Import from the following gaming services:
    - Epic Games
    - GOG
    - PSN
    - Steam
- Easy to use
- Efficient data storage and retrieval
- Small application size (< 100 MB)
- FOSS (Free and Open-Source Software)
- Cross-Platform 
	- MacOS
    - Linux 
    - Windows
- Exports to:
	- CSV 
	- SQL 
	- MD 
- Imports from:
    - CSV
    - SQL (destructive)
    - TXT
    - Game Name Search
    - Manual Entry
- Data Sourced from:
    - [HLTB](https://howlongtobeat.com/)
    - [Completionator](https://completionator.com/)
- Uses GoLang
- Chippy the Chipmunk will be extremely happy

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

## Why make this program


I often find myself wanting to knock another game off of the gaming backlog, but see the massive list and feel unsure as to which to play next.
One key factor that weights into whether I will select a particular game will be the time it takes to complete it. 
For example, to complete a game in 10 hours would be more desirable for a single week experience vs an 80+ hour game which would probably last several weeks to months. 
This information would make selecting a game to remove from the backlog much easier.

I don't want to always have to go to HLTB and search for my game, then compare it to others that I then have to search as well. 
This process could take several minutes at a time, depending on how many games I want to compare. 
I want a single point of information of times and games that I can hold that is catered to me as an individual based on my list of games.

### TLDR; 

I want an individualized stored set of data that contains relevant information to the games that I have/own.

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

## Ethical Conundrum


This program uses a technique called [Web Scraping](https://en.wikipedia.org/wiki/Web_scraping) which is morally gray.
When visiting websites, typically they expect a real user and prevent bots/programs from accessing without proper/prior permissions.
This program automates the process for scraping the time data for the particular games.
As such, it is possible to view this program as malicious in usage towards websites such as HLTB.

If you feel this is unethical, then I highly encourage you to **not use the program**.
This program is not illegal, nor was it created with bad intentions.

The access for API for HLTB in particular has been requested for several years ([See: this forum thread](https://howlongtobeat.com/forum/thread/807/1)).
If the processes are able to be done via API as opposed to scraped, then I will happily consider adjusting the inner workings of the program for HLTB and others.

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

# Use Guide Suggestions

- Use a VPN or other IP masking device to avoid getting IP blacklisted if you are importing many games at once.
- Don't open multiple instances of the program at once. 
I have no clue what will happen and don't intend to develop for such an edge case unless necessary. 
This program was created with a single window to be used in mind.
- I am not responsible for any consequences should you misuse this software.
- If you have any suggestions or find any bugs feel free to submit a ticket [here](https://github.com/EZRA-DVLPR/LOGL/issues) with the proper tags.

## YOUTUBE TUTORIAL SERIES

[Link](https://www.youtube.com/playlist?list=PL_gNvZlhoitBNANmcZFgoQpT1FjZiBs7I)

## ONLINE PDF MANUAL


[Link](https://github.com/EZRA-DVLPR/LOGL/blob/main/docs/PDF/Manual.pdf)

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

# INSTALLATION GUIDE

## Windows:

1) Download the `W-LOGL.zip` file
2) Unzip the folder (Winrar/7Zip is good)
3) Drag and drop `LOGL.exe` onto whatever location in your file system you like
4) Enjoy!

## MacOS

1) Download the `D-LOGL.zip` file
2) Unzip the folder (The Unarchiver is good)
3) Drag and drop `LOGL.app` into the `Applications` folder
4) Enjoy!

## Linux

1) Download the `LOGL.tar.xz` file
2) Unzip the contents (I used `tar -xf`)
3) Install with make `sudo make install`
4) Enjoy!

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

# LEGAL

-----------------------------------------------------------------------------------------------------------------------------------------------------------------

## License  

- The **code and documentation** in this repository are licensed under the **MIT License**. See [LICENSE](LICENSE) for details.  
- The **project name and icon** are licensed under the **Creative Commons Attribution 4.0 International License**. See [LICENSE_ICON_NAME.md](LICENSE_ICON_NAME.md) for details.  

This is an independent project done by me with no further collaborators.

This project has no affiliation with Google, Bing, HowLongToBeat, nor Completionator.
