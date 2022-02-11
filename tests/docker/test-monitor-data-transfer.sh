docker exec -it warehouse2 cat /var/supplywatch/warehouse/Data > /tmp/warehouse2_data_send.txt
docker exec -it warehouse1 cat /var/supplywatch/warehouse/Data > /tmp/warehouse1_data_send.txt
docker exec -it monitor cat /var/supplywatch/monitor/products-warehouse1DataExchange > /tmp/warehouse1_data_receiv.txt
docker exec -it monitor cat /var/supplywatch/monitor/products-warehouse2DataExchange > /tmp/warehouse2_data_receiv.txt

diff /tmp/warehouse2_data_send.txt /tmp/warehouse2_data_receiv.txt
diff /tmp/warehouse1_data_send.txt /tmp/warehouse1_data_receiv.txt
