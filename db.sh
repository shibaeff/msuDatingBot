docker run -d \
    --env MONGODB_PORT=27017 \
    --env MAX_BACKUPS=5 \
    --volume host.folder:/backup
    tutum/mongodb-backup

docker run -d  --env MAX_BACKUPS=5 --link mongodb:mongodb -v host.folder:/backup tutum/mongodb-backup