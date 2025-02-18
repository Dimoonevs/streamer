HOST=13.48.46.221
HOMEDIR=/var/www/hls-streamer/
USER=dima

streamer-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/streamer-linux-amd64 ./

upload-streamer-service: streamer-linux
	rsync -rzv --progress --rsync-path="sudo rsync" \
		./bin/streamer-linux-amd64  \
		./cfg/prod.ini \
		./cfg/restart.sh \
		$(USER)@$(HOST):$(HOMEDIR)

restart-streamer-service:
	echo "sudo su && cd $(HOMEDIR) && bash restart.sh && exit" | ssh $(USER)@$(HOST) /bin/sh

run-local:
	go run main.go -config /home/dima/MyProjectNameCompany/hls-streamer/cfg/local.ini