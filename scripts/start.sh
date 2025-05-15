#!/bin/sh
set -ex

/app/wait_for_db.sh

max_retries=5
count=0
until goose -dir /app/migrations mysql "$DB_USER:$DB_PASSWORD@tcp($DB_HOST:$DB_PORT)/$DB_NAME?parseTime=true" up; do
    count=$((count+1))
    if [ $count -ge $max_retries ]; then
        echo "Migration failed after $max_retries attempts"
        exit 1
    fi
    echo "Migration failed, retrying in 3 seconds..."
    sleep 3
done

exec /app/main