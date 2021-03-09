#!/bin/bash

zerotier-one -d
sleep 1
exec /bin/ztdaemon --networkID $NETWORK_ID --token $ZT_TOKEN --name $NODE_NAME
