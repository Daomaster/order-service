# Order Service
### How to start

Insert the Google Map Api Key to docker-compose.yml file

```sh
MAP_API_KEY={api_key}
```

Change the permission of script
```sh
chmod +x start.sh
```

Run the start script `start.sh`

```sh
$ ./start.sh
```

### How to run unit test
```sh
$ go test -v ./...
```

### Notes
- page query will start at `0` on the route `GET /orders`
- the service will start after the database is started
- no need to init database, the service will auto migrate it
- if you want persistent database, just add a volume to the docker-compose