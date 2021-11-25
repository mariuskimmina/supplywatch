# Sensor UPD test

To test the reliability with which our sensors send data to the warehouse we can configure 
a fix number of packets that the sensors are going to send. The warehouse tracks how many packets
it receives from each sensor - we can compare the configured number of packets with the number of 
received packets.

## Example test conifg

```
warehouse:
  listenIP: "0.0.0.0"
  udpPort: 4444
  tcpPort: 8000

sensorWarehouse:
  udpPort: 4444 # should match warehouse udpPort
  delay: 500 # time between packets in milliseconds
  numOfPackets: 100 # 0 means infinite
```

with this config, each sensor will send 100 packets in 50 seconds.
The `/tmp/logcount` file on the warehouse container is going to tell us how many of these
have actually arrived.

## Checking arrvied packets at the warehouse

```
cat /tmp/logcount
{"SensorID":"ac7e40d1-4dd1-11ec-8759-0242ac190003","Counter":100}
{"SensorID":"ac91af5b-4dd1-11ec-a9fd-0242ac190004","Counter":100}
```
