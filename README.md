# ihostproxy

> A self-hosted SOCKS5 proxy with an extremely basic authentication webui

## Todo

Web UI
[] Change password

Proxy
[] FQDN
[] IP Whitelist

## Running

### Running with the default ports and a custom volume for saving the db data

_NOTE: password cannot contain special symbols._

```
docker run -it --name ihostproxy -p 1080:1080 -p 8080:8080 -v v_ihostproxy:/db -e WEBUI_USER=admin -e WEBUI_PASS="YOUR PASSWORD HERE" -e DB_PATH=/db/main.db alexng353/ihostproxy
```
