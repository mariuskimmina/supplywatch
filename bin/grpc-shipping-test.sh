docker exec -it warehouse_2 cat /var/shipping_send_log > warehouse2_send.txt
docker exec -it warehouse_2 cat /var/shipping_receiv_log > warehouse2_receiv.txt
docker exec -it warehouse_1 cat /var/shipping_send_log > warehouse1_send.txt
docker exec -it warehouse_1 cat /var/shipping_receiv_log > warehouse1_receiv.txt

diff warehouse1_send.txt warehouse2_receiv.txt

diff warehouse2_send.txt warehouse1_receiv.txt
