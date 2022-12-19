#!/bin/bash

if [ "$#" -ne 1 ]; then
    exit 1
fi

username_mac=$1

vrf_id=`/home/jetb/cli.sh sh pppoe subsc | grep -i $username_mac | awk '{print $1}'`

if [ -z "$vrf_id" ]; then
    exit 1
fi

id_len=${#vrf_id}
if [ "$id_len" -gt 7 ]; then
    exit 1
fi
#echo $vrf_id
/home/jetb/cli.sh pppoe disconnect $vrf_id
echo "$vrf_id $username_mac cutted off" >> /tmp/cut-off.log

