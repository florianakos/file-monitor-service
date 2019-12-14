#!/bin/bash

# set up some variables 
TARGET=monitored_folder
PROCESSED=processed_files
DB_FILE=database/stats.db
LOG_FILE=logs/service_history.log

# start inotify servic in TARGET folder
inotifywait -m -e create -e moved_to --format "%f" $TARGET |
	# process every notification from inotify as they come (block if nothing)
	while read FILENAME; do
		# allow larger files uploads to finish before triggering move ang gzip
		sleep 5
		
		# move file to processed_files folder
		mv "$TARGET/$FILENAME" "$PROCESSED/$FILENAME"
		
		# calculate original size before compression
		ORIGSIZE=$(stat "$PROCESSED/$FILENAME" --print="%s")
	
		# compress file via gzip
		gzip -9 -f "$PROCESSED/$FILENAME"

		# calculate new size after compression
		PROCSIZE=$(stat "${PROCESSED}/${FILENAME}.gz" --print='%s')

		# get compression ration
		PERCENT=$(echo "scale=4 ; ($ORIGSIZE - $PROCSIZE) / $ORIGSIZE * 100" | bc)

		# craft and append log message to file-based logging
		LOG_MSG=$(echo "$(date "+%Y-%m-%d %H:%M:%S"): New file detected ($FILENAME). Original size: $ORIGSIZE --- Zipped size: $PROCSIZE --- Compression: $PERCENT%")
		echo "$LOG_MSG" | tee -a $LOG_FILE

		# insert stats about current iteration to sqlite db
		sqlite3 $DB_FILE "INSERT INTO data (time, file, orig_size, comp_size, comp_rate) VALUES ($(date +%s), '${PROCESSED}/${FILENAME}.gz', $ORIGSIZE, $PROCSIZE, $PERCENT);"

	done