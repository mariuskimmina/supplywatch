docker exec -it warehouse2 cat /var/supplywatch/grpc/client/sendlog > /tmp/warehouse2_send.txt
docker exec -it warehouse2 cat /var/supplywatch/grpc/server/receivLog > /tmp/warehouse2_receiv.txt
docker exec -it warehouse1 cat /var/supplywatch/grpc/client/sendlog > /tmp/warehouse1_send.txt
docker exec -it warehouse1 cat /var/supplywatch/grpc/server/receivLog > /tmp/warehouse1_receiv.txt

diff /tmp/warehouse1_send.txt /tmp/warehouse2_receiv.txt
diff /tmp/warehouse2_send.txt /tmp/warehouse1_receiv.txt
