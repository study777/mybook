# NTP 服务器
```
yum install -y ntp  

ntpdate 1.asia.pool.ntp.org


> /etc/ntp.conf

vim  /etc/ntp.conf

driftfile /var/lib/ntp/drift
server 1.asia.pool.ntp.org
restrict default  nomodify
server  127.127.1.0     # local clock
fudge   127.127.1.0 stratum 10
includefile /etc/ntp/crypto/pw
keys /etc/ntp/keys

systemctl enable ntpd  
systemctl start ntpd  
```
# DNS

```
yum -y install bind*

> /etc/named.conf

vim  /etc/named.conf

options {
        directory "/var/named";
        allow-query     { 0.0.0.0/0; };
        listen-on port 53 { any; };
        #allow-query-cache { any; };
        #forwarders { 172.16.10.2; };
};

zone "example.com" {
        type master;
        file "example.com.zone";
};

zone "2.16.172.in-addr.arpa" {
        type master;
        file "172.16.2.zone";
};

zone "pass.com" {
        type master;
        file "pass.com";
};


vim /var/named/172.16.10.zone

$TTL 3600
@             IN        SOA   10.16.172.in-addr.arpa. admin.example.com. (
                                 20150520
                                   1H
                                   15M
                                   1W
                                   1D)
               IN        NS        dns.example.com.
170            IN        PTR        dns.example.com.
161          IN       PTR        master-161.example.com.
162          IN       PTR        node-162.example.com.
163          IN      PTR         node-163.example.com.
164          IN      PTR         node-164.example.com.



vim /var/named/example.com.zone

$TTL 3600
@              IN  SOA  dns.example.com. admin.example.com. (
                        20150520
                        1H
                        15M
                       1W
                       1D)

               IN  NS   dns.example.com.
dns            IN   A    172.16.10.170
master-161         IN   A    172.16.10.161
node-162        IN   A    172.16.2.162
node-163         IN   A    172.16.2.163
node-164        IN   A    172.16.2.164



systemctl start named

systemctl enable  named





```
