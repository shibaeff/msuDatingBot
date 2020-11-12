kill -9  $(ps | grep main | awk '{print $1}')
git stash
git pull
rm ./out.txt
chmod +x *sh
