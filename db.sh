docker run -d \
    -p 27017:27017 \
    --env MONGODB_PORT=27017 \
    --env MAX_BACKUPS=5 \
    --volume ~/db:/backup \
    --name mongodb \
    mongo

docker run -d --env MAX_BACKUPS=5 --link mongodb:mongodb --name backup -v ~/db:/backup tutum/mongodb-backup