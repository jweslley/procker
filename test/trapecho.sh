#!/bin/sh

trap "echo bye" SIGTERM

sleep $1
echo -n $2
