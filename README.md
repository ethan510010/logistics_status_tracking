###  Logistics Status Tracking Website

* Use Go/Echo as programming language and web router library
* Use MySQL as database and redis as cache server
* complete query and fake api, ex:
  * `/query?sno=XXXXXXXX`
  * `/fake?num=5`
* Use cache aside pattern for cache strategy, and use hash as storage data structure
  * if data exists in cache, grabs the data from cache, otherwise grabs from DB and set to cache
  * cache ttl sets to 1hr
* Use systemd to start my web server
* Set the cron job for generating report that will be updated to S3 