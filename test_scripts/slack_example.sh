#!/bin/sh
# slack_example.sh - demonstrates per-plugin config vs per-plugin secrets

# receives three kinds of data, from three different places:
# BUILD_*       build data, same for every plugin, from BuildResult
# CHANNEL       per-plugin config, a literal value from confg.json
# SLACK_TOKEN   per-plugin secret, injected from a host env var by name

echo "channel: ${CHANNEL:-<unset>}" >&2
echo "build:: $BUILD_ID on $BUILD_BRANCH ($BUILD_STATUS)" >&2

# never print a seccret, show only enough to prove it arraived
if [ -n "$SLACK_TOKEN" ]; then
    echo "token: present, ${#SLACK_TOKEN} chars, starts ${SLACK_TOKEN%"${SLACK_TOKEN#????}"}..." >&2
else
    echo "token: MISSING" >&2
    exit 1
fi

# prove the allowlist holds: a var the host has but never granted us.
echo "leaked: EXAMPLE_SLACK_TOKEN=${EXAMPLE_SLACK_TOKEN:-<not visible, correct>}" >&2

echo "would post to $CHANNEL" >&2

exit 0

