#!/bin/bash

mage generateCLICommandDocs

mkdir -p /tmp/unified/

git diff docs/cmd/*.md > /tmp/unified/cmd_docs.diff

diff_size=$(wc -c /tmp/unified/cmd_docs.diff | awk '{print $1}')

if [ $diff_size == 0 ]; then   
    echo "OK - No difference"; 
    exit 0
else
    echo "ERROR - 'mage generateCLICommandDocs' caused a diff to appear!"
    echo "Please run 'mage generateCLICommandDocs' and check in changes to docs/cmd/*.md"
    exit 1
fi

rm /tmp/unified/cmd_docs.diff
