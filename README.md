# ihostproxy

> A self-hosted SOCKS5 proxy with an extremely basic authentication webui

## Running

### Running with the default ports and a custom volume for saving the db data

```
docker run -it -p 1080:1080 -p 8080:8080 -v testvolume:/db -e WEBUI_USER=admin -e WEBUI_PASS=adminadmin -e DB_PATH=/db/main.db ihostproxy
```
