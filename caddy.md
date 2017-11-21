setenforce  0

systemctl  stop firewalld





wget  https://github.com/mholt/caddy/releases/download/v0.10.10/caddy_linux_amd64.tar.gz





tar zxv caddy_v0.10.10linux_amd64personal.tar.gz



curl https://getcaddy.com | bash -s personal



vim Caddyfile

*:8080

markdown / {

  ext .md

}
