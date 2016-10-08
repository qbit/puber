#!/bin/sh

function getKey {
    echo -n "Press yubikey: "
    read YK
}
KEY="potato"

getKey
RES=$(curl -sH "Content-Type: application/json" -X POST -d "{\"user\": \"${KEY}\", \"pubKey\": \"asdfasdfasdfasdf\", \"yKey\": \"${YK}\"}" http://localhost:8081/add)

if [ "$RES" == "Added" ]; then
	echo "[success]: added key '${KEY}'"
else
	echo "[fail]: adding key '${KEY}' (${RES})"
fi

getKey
RES=$(curl -sH "Content-Type: application/json" -X POST -d "{\"user\": \"debug\", \"yKey\": \"${YK}\", \"pubkey\": \"debug key two\"}" http://localhost:8081/rm)
if [ "$RES" == "Removed 'debug key two' from 'debug'" ]; then
	echo "[success]: removing key 'debug key two' for user 'debug'"
else
	echo "[fail]: removing key '${KEY}' (${RES})"
fi

getKey
RES=$(curl -sH "Content-Type: application/json" -X POST -d "{\"user\": \"debug\", \"yKey\": \"${YK}\"}" http://localhost:8081/rm)
if [ "$RES" == "Removed" ]; then
	echo "[success]: removing all keys 'debug'"
else
	echo "[fail]: removing key '${KEY}' (${RES})"
fi
