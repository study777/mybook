## 单网口

yum -y install iptables-services


iptables -F


iptables -t nat -F

 
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf



iptables -t nat -A POSTROUTING -s 172.16.10.110   -j SNAT --to 172.16.10.100


iptables -A FORWARD -s 172.16.10.110  -j ACCEPT


iptables -L

## 双网口

yum -y install iptables-services

iptables -F


iptables -t nat -F

 
echo "net.ipv4.ip_forward=1" >> /etc/sysctl.conf

sysctl -p

iptables -A FORWARD -s 9.9.9.0/24 -j ACCEPT

iptables -t nat -A POSTROUTING -s 9.9.9.0/24 -j SNAT --to  172.16.101.1


### dnat

iptables -t nat  -F

iptables -t nat -A PREROUTING -d 172.16.101.249 -p tcp -m tcp --dport 3301 -j DNAT --to-destination 182.16.10.12:3306


iptables -t nat -A POSTROUTING -d 182.16.10.12  -p tcp  -m tcp --dport 3306 -j SNAT --to-source 182.16.10.11


iptables -A FORWARD -o eth0 -d 182.16.10.12  -p tcp --dport 3306 -j ACCEPT


iptables -A FORWARD -i eth0  -s 182.16.10.12 -p tcp --sport 3306 -j ACCEPT

