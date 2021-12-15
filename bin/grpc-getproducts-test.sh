docker exec -it warehouse_2 cat /var/all_products_send_log > warehouse2_products_send.txt
docker exec -it warehouse_2 cat /var/all_products_receiv_log > warehouse2_products_receiv.txt
docker exec -it warehouse_1 cat /var/all_products_send_log > warehouse1_products_send.txt
docker exec -it warehouse_1 cat /var/all_products_receiv_log > warehouse1_products_receiv.txt

diff warehouse1_products_send.txt warehouse2_products_receiv.txt

diff warehouse2_products_send.txt warehouse1_products_receiv.txt
