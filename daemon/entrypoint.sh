#!/bin/bash

zerotier-one -d
/bin/ztdaemon --networkID $NETWORK_ID --token $ZT_TOKEN --name $NODE_NAME
