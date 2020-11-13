kill -9  $(ps | grep main | awk '{print $1}')
git stash
git pull
chmod +x *sh
