kill -9  $(ps | grep main | awk '{print $1}')
git pull
rm ./out.txt
