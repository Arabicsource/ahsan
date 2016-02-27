# Ahsan - A Monitoring - Download - Extracting service for Maktabah 

------------------------------------------------------------------

Ahsan is part of an ecosystem of 'services' that make up the Maktabah Application. A lot of 
duties that Ahsan will carry out for the most part will consist of checking if any new books
have been published by shamela.ws and then download it. Ahsan will be responsible for passing
the job of extracting and indexing the file's content on to the queue for another service to 
carry it out.


#### Functionality
------------------
Parsing web pages for download links since Shamela.ws has no public API for new content uploaded on 
its website. 



##### TODO 

- [ ]   Scrape shamela.ws ever interval period (set by flags)
- [ ]   Store urls in json format
- [ ]   Upon addition of new urls (of books) trigger an event for another service to poll the json file
- [ ]   Check shamela.ws for new books 
- [ ]   todo ....  




-------------------------------------------------------------------

#### AUTHORS

* aboo "shaybix" shayba



#### CONTRIBUTORS





