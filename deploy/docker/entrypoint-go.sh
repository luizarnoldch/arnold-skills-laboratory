#!/bin/sh
set -e

if [ -x ./db_tool ]; then
	./db_tool up
fi

exec "$@"
