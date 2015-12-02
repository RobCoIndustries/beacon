#!/bin/bash
set -euo pipefail

rm -f certificates/*.srl certificates/*.csr certificates/*.pem

export KEYMASTER="docker run --rm -v $(pwd)/certificates/:/certificates/ cloudpipe/keymaster"

${KEYMASTER} ca
${KEYMASTER} signed-keypair -n beacon -h 127.0.0.1 -s IP:127.0.0.1 -p server
${KEYMASTER} signed-keypair -n beacon-cli -h 127.0.0.1 -s IP:127.0.0.1 -p client
