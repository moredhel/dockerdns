This is a project inspired by a talk given at NixCon 2018.


Wouldn't it be great to have nice, memorable and stable names for all your containering needs?
This is just that project, a dns server which monitors the creation of docker containers and creates DNS entries.

Let's try it out! (note: you may have to turn off/kill dnsmasq with a `kill dnsmasq`)

``` shell
docker-compose up -d
```

And now, let's see if we can hit the DNS server's health-check endpoint.


``` shell
# it seems this may not work if curl is not compiled with the required flags
curl -i `dig @127.0.0.1 +short core.dns.local`:8080/health
```

# TODO

- write a nixos module for CoreDNS so I can run it natively rather than having to have a container running. (#cloudnative?)
- it would be nice if there was a way to have the DNS server ignored all domains except *.local (or other I guess)
