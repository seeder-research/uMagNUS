#! /bin/bash 

rm -rf *.out http\:
killall mumax3cl-httpfsd
mumax3cl-httpfsd -l :35377 &
sleep 1s
mumax3cl $(for f in *.mx3; do echo -n ' ' http://localhost:35377/$f; done) || exit 1
killall mumax3cl-httpfsd
if [ -e http\: ]; then
		echo "http:" exists
		exit 2
fi;
