docker run -d \
    --env MONGODB_PORT=27017 \
    --env MAX_BACKUPS=5 \
    --volume host.folder:/backup \
    --name mongodb \
    tutum/mongodb

docker run -d --env MAX_BACKUPS=5 --link mongodb:mongodb --name backup -v host.folder:/backup tutum/mongodb-backup