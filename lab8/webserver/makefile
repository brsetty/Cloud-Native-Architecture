

run:
	sudo docker run -d --network item-server  -p 8050:8050 --name webserver webserver
	sudo docker run -d --network item-server --network-alias mongo -p 27017-27019:27017-27019 --name mongodb mongo:latest

remove:
	sudo docker container stop webserver mongodb
	sudo docker container rm webserver mongodb


sim:
	curl "localhost:8000/add?item=bag&price=15"
	curl "localhost:8000/add?item=bag&price=15"
	curl "localhost:8000/add?item=shoe&price=25"
	curl "localhost:8000/list"
	curl "localhost:8000/update?item=bag&price=20"
	curl "localhost:8000/list"
	curl "localhost:8000/delete?item=shoe"
	curl "localhost:8000/list"
