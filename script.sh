#!/bin/bash

TARGET=monitored_folder
PROCESSED=processed_files

inotifywait -m -e create -e moved_to --format "%f" $TARGET \
        | while read FILENAME
                do
			mv "$TARGET/$FILENAME" "$PROCESSED/$FILENAME"

			ORIGSIZE=$(stat "$PROCESSED/$FILENAME" --print="%s")
                        echo "Detected $FILENAME, moving and zipping"
			echo "Orig size: $ORIGSIZE"

			gzip -9 -f "$PROCESSED/$FILENAME"
                        
			#PROCSIZE="$(stat "${PROCESSED}/${FILENAME}.gz" --print="$s")"
			PROCSIZE=$(stat "${PROCESSED}/${FILENAME}.gz" --print='%s')
			echo "New size: $PROCSIZE"
			
			PERCENT=$(echo "scale=4 ; ($ORIGSIZE - $PROCSIZE) / $ORIGSIZE * 100" | bc)
			echo "Compression: $PERCENT%"
\
			sqlite3 stats.db "INSERT INTO data (time, file, orig_size, comp_size, comp_rate) VALUES ($(date +%s), '${PROCESSED}/${FILENAME}.gz', $ORIGSIZE, $PROCSIZE, $PERCENT);"



                done
