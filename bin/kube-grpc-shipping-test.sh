kubectl exec -it deploy/warehouse2 -- ls /var/
WH1=$($(kubectl get pod -o name |  sed 's/pod\///' | grep warehouse1))
WH2=$($(kubectl get pod -o name |  sed 's/pod\///' | grep warehouse2))
kubectl cp default/$WH1:/var/shipping_send_log /tmp/warehouse2_send.txt
kubectl cp default/$WH2:/var/shipping_receiv_log /tmp/warehouse1_receiv.txt

#diff /tmp/warehouse2_send.txt /tmp/warehouse1_receiv.txt

#diff warehouse2_send.txt warehouse1_receiv.txt
