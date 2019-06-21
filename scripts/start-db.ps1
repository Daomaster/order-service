# spawn a local mysql database without consistent storage
docker run --name test-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=test123 -e MYSQL_DATABASE=order-service -d mysql:5.7