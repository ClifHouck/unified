#!/bin/bash

mage generateStreamHandlers

mkdir -p /tmp/unified/

git diff client/protect_*_stream_handler.go > /tmp/unified/stream_handlers.diff

diff_size=$(wc -c stream_handlers.diff | awk '{print $1}')

if [ $diff_size == 0 ]; then   
    echo "OK - No difference"; 
    exit 0
else
    echo "ERROR - 'mage generateStreamHandlers' caused a diff to appear! "
    echo "Please check in changes to client/protect_*_stream_handler.go"
    exit 1
fi

rm /tmp/unified/stream_handlers.diff
