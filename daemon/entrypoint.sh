#!/bin/bash

zerotier-one -d
sleep 1
/bin/ztdaemon --networkID $NETWORK_ID --token $ZT_TOKEN --name $NODE_NAME