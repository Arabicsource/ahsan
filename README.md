# Ahsan - A Monitoring - Download - Extracting service for Maktabah 

------------------------------------------------------------------

Ahsan is part of an ecosystem of 'microservices' that make up Maktabah.  
Ahsan responsibilities for the most part will consist of checking if there are new books and if they
have been published by shamela.ws, if so then to download it. 

On the contrary, I'm considering to push the urls to be downloaded to a messaging queue and then let another service download, 
extract and index it in ElasticSearch. 


#### Functionality
------------------
Parsing web pages for download links since Shamela.ws has no public API for new 
content uploaded on the website. 


##### TODO

- [ ]   Scrape shamela.ws category pages in an interval period (set by flags)
- [ ]   Store retrieved urls in json file
- [ ]   compare previous json file with the newly scraped urls and download the the difference between the two. 
- [ ]   Push the new urls into a messaging queue (like RabbitMQ)


##### FIXME

- [ ]   Category links scraped, however work on goroutines for scraping individual category pages.
- [ ]   create channel to pass through links of categories
- [ ]   refactor the regex code and have it as its own package to be imported.


-------------------------------------------------------------------

#### AUTHORS

* aboo "shaybix" shayba



#### CONTRIBUTORS

