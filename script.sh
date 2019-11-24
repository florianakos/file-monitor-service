#!/bin/bash

TARGET=monitored_folder
PROCESSED=processed_files

inotifywait -m -e create -e moved_to --format "%f" $TARGET \
| while read FILENAME
    do
		mv "$TARGET/$FILENAME" "$PROCESSED/$FILENAME"

		ORIGSIZE=$(stat "$PROCESSED/$FILENAME" --print="%s")
 
		gzip -9 -f "$PROCESSED/$FILENAME"

		PROCSIZE=$(stat "${PROCESSED}/${FILENAME}.gz" --print='%s')

		PERCENT=$(echo "scale=4 ; ($ORIGSIZE - $PROCSIZE) / $ORIGSIZE * 100" | bc)
        
        LOG_MSG=$(echo "$(date "+%Y-%m-%d %H:%M:%S"): New file detected ($FILENAME). Original size: $ORIGSIZE --- Zipped size: $PROCSIZE --- Compression: $PERCENT%")
        echo $LOG_MSG | tee -a service_history.log
        #echo $LOG_MSG >> service_history.log
            
		sqlite3 stats.db "INSERT INTO data (time, file, orig_size, comp_size, comp_rate) VALUES ($(date +%s), '${PROCESSED}/${FILENAME}.gz', $ORIGSIZE, $PROCSIZE, $PERCENT);"

    done