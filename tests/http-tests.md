# HTTP Tests

To test our HTTP Endpoints we do the following:

* Make sure there are no logs for today yet
    * logs are saved unter `/var/supplywatch/log/warehouselog-{TIMESTAMP}`    
    * a file with todays day already exists, delete it
* Configure the number of packets you want the sensors to send in `config.yml`
    * to make the verification easy, choose a low number
* Start the application
* wait for all sensors to finish
* go to http://localhost:8000/sensorhistory?date=12-01-2021
* compare the http response with the log file that has been created for today - they should match
