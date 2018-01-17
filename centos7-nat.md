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




