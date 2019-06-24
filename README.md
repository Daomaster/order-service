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

### Notes
-page query will start at `0` on the route `GET /orders`