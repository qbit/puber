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
    #	getKey
    #	RES=$(curl -sH "Content-Type: application/json" -X POST -d "{\"user\": \"${KEY}\", \"yKey\": \"${YK}\"}" http://localhost:8081/rm)
    #	if [ "$RES" == "Removed" ]; then
    #		echo "[success]: removing key '${KEY}'"
    #	else
    #		echo "[fail]: removing key '${KEY}' (${RES})"
    #	fi
else
    echo "[fail]: adding key '${KEY}' (${RES})"
fi
