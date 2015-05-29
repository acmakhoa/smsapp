#!/bin/bash
 
# The binary is most likely to have the current directory name
PROCESS=${PWD##*/}
 
# We kill the previous process
killall ${PROCESS}
 
# We launch the process
go build
./${PROCESS} &
 
echo "Waiting for changes..."
inotifywait -mqr --timefmt '%d/%m/%y %H:%M' --format '%T %w %f' -e moved_to . ${GOPATH} | \
while read date time dir file; do
    ext="${file##*.}"
 
    # We only monitor go files
    if [[ "$ext" = "go" || "$ext" = "html" ]]; then
        echo "$file changed @ $time $date, rebuilding..."
        
        # We build it
        go build -v >.build_status 2>&1
 
        # If everything went fine
        if [[ "$?" == "0" ]]; then
        	# We kill the previous process
        	killall -9 ${PROCESS}
 
        	# We launch the process
        	./${PROCESS} || if [ $? != 137 ]; then notify-send -i error "Program stopped" "Program ${PROCESS} stopped !"; fi &
 
	       	# We report it
        	notify-send -i emblem-default "OK" "Compiled OK"
 
        	# Then if we have more than one line in the building process
        	# it means that we probably had to compile a dependencies. Which
        	# is why we will prefetch them for next time
        	if [[ `cat .build_status | wc -l` -gt 1 ]]; then
        		echo "Rebuilding dependencies..."
        		go get -v >.dep_build_status 2>&1 && \
        		notify-send -i emblem-default "Rebuilt dependencies" "Rebuild dependencies" || \
        		notify-send -i error "Dependencies error" "`cat .dep_build_status`"
        	fi
        else
        	# If it fail, we report it and display the error
        	notify-send -i error "Error" "`cat .build_status`"
        	cat .build_status
        fi
    fi
done
