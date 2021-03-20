---
Title: Proxies
---

# proxies

Alternative to running holden on port 80 on your server, you can easily reverse proxy to the holden server or the holden docker container. The setting in `config.toml` called "prefix" will be useful if you want to reverse proxy to holden on a subpath, i.e. `domain.com/docs` etc.

For simple setups, or if your networking infrastructure is modular, we recommend either Caddy or Traefik depending on your implementation.

# examples

## caddy

Caddy Subdomain Proxy to a docker container
```
docs.cya.cx {
	reverse_proxy holden:11011
}
```

Caddy Subpath Proxy to a docker container
```
domain.com {
	root * /srv/domain.com
	file_server {
		index index.html
	}

	redir /holden {uri}/
	route /holden/* {
		reverse_proxy holden:11011
	}
}
```

## nginx

```nginx
server {
    listen 80; # or whatever 443 etc it's your setup
    server_name domain.com;

    location / {
        proxy_pass holden:11011;
    }

    #or 

    location /some/path/ {
        proxy_pass holden:11011;
    }
}
```

## Other

As a rule of thumb, just reverse proxy (something) to the holden server, and set holden's prefix as necessary. holden doesn't require any special settings.