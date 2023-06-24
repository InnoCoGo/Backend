#!/bin/bash

while true; do
	response=$(curl -s "http://localhost:8000/ping")
	if [ "$response" == "pong" ]; then
		echo "Received 'pong', proceeding with database migration."
		break
	else 
		echo "Wait for 'pong' response ..."
		sleep 5
	fi
done

make migrate
