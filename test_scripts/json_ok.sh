#!/bin/sh
# json_ok.sh

input=$(cat)

echo "got: $input" >&2      #stderr - stout is the protocol
# >&2 redirects the echo to stderr, echo "x" writes to stdout. if we didnt include
# >&2, echo would write to stdout and stdout buffer would contain
# got: input
# {"ok": true, "message": "posted"}
# host cant unmarshall this, host reports plugin as failed but script runs correctly
echo '{"ok": true, "message": "posted"}'
exit 0
