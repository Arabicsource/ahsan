# Ahsan - A Monitoring - Download - Extracting service for Maktabah 

------------------------------------------------------------------

Ahsan is part of an ecosystem of 'services' that make up the Maktabah Application. A lot of 
duties that Ahsan will carry out for the most part will consist of checking if any new books
have been published by shamela.ws and then download it. Ahsan will be responsible for passing
the job of extracting and indexing the file's content on to the queue for another service to 
carry it out.

As Ahsan currently stands it may be broken part into further smaller parts. The web crawling part 
being put in a seperate package, and functioning as an imported package independent of Ahsan and other
Maktabah services. However, this may very well be delayed for a much later period in the project.


#### Functionality
------------------
Parsing web pages for download links since Shamela.ws has no public API for new 
content uploaded on its website. 



##### TODO 

- [ ]   Scrape shamela.ws ever interval period (set by flags)
- [ ]   Store urls in json format
- [ ]   new urls (of books) trigger an event for another service to poll the json file
- [ ]   Check shamela.ws for new books 


##### FIXME

- [ ]   refactor the regex code and have it as its own package to be imported.

This could be similar to python's BeautifulSoup (bs4), and check it out for guide. However, 
bear in mind the purpose behind it is ease of use and not to write a bloated package that
tries to do everything. Focus on main purpose of Ahsan and write the package to enable yourself 
to pursue that objective.  

-------------------------------------------------------------------

#### AUTHORS

* aboo "shaybix" shayba



#### CONTRIBUTORS

