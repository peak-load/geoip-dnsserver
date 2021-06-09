# geoip-dnsserver
Simple DNS server in Golang, which answers with client IP address and its GeoIP data (Maxmind DB is used). 

I wrote this service inspired by presentation at RIPE 75 meeting https://ripe75.ripe.net/presentations/20-A-curious-case-of-broken-DNS-responses-RIPE-75.pdf

Where service maxmind.test-ipv6.com was used, I tied to reach @jfesler asking for help / source code of that service with no success, so I ended up with writting this service.  

Service reply with GeoIP data on every DNS TXT query which begins with geoip 

**Prerequisites:**
1) You configure your DNS zone to use IP address of publicly available service running on port 53.   
2) Configure your firewall rules.
3) GeoIP2-ISP.mmdb database in the same directory.

**Building**
```
git clone https://github.com/peak-load/geoip-dnsserver
cd geoip-dnsserver
go get -v .
go build . 
```

**Running test server & testing**
```
$ ./geoip-dnsserver
$ dig @127.0.0.1 -t txt +short geoip

# Query from your host
$ dig +short -t txt geoip.mydomain.com
"ip='83.169.xx.xx',asn='3209',asn_organization='Vodafone GmbH',isp='Vodafone Germany Cable',organization='Vodafone Germany Cable'"

# Query from your host over Google Public DNS server
$ dig @8.8.8.8 +short -t txt geoip.mydomain.com
"ip='172.253.199.3',asn='15169',asn_organization='GOOGLE',isp='Google',organization='Google'"
```
