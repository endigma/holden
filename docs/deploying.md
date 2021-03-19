# Deployment Guide

## Baremetal

clone the repo, build and run.

## Docker

`docker pull ghcr.io/endigma/holden`

run it, mount your docroot and config as youd expect. the docker image expects the config to be at `/config.toml`.

```yaml
version: "3.5"
services:
  holden:
    image: ghcr.io/endigma/holden:latest
    ports:
      - 80:11011
    container_name: holden
    environment:
      - PUID=1000
      - PGID=1000
    volumes:
      - ./config.toml:/config.toml
      - ./docroot:/docroot
    restart: unless-stopped
```

you can also mount `/assets/serve/vars.css` and `/assets/public`
