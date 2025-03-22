# ihostproxy

> A self-hosted SOCKS5 proxy with an extremely basic authentication webui

## Todo

Web UI
- [ ] Change password

Proxy
- [ ] FQDN
- [ ] IP Whitelist

## Running

### Running with the default ports and a custom volume for saving the db data

_NOTE: password cannot contain special symbols._

```
docker run -d --name ihostproxy -p 1080:1080 -p 8080:8080 -v v_ihostproxy:/db -e WEBUI_USER=admin -e WEBUI_PASS="YOUR PASSWORD HERE" -e DB_PATH=/db/main.db alexng353/ihostproxy
```

## Contributing

<!-- TODO: write correct build instructions -->

```bash
go install github.com/a-h/templ/cmd/templ@latest # this is a required build tool

git clone https://github.com/alexng353/ihostproxy
cd ihostproxy

templ generate
go build
```

## Attributions

Project is roughly based off of [serjs/socks5-server](https://github.com/serjs/socks5-server)
