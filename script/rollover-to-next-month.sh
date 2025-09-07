#!/bin/bash

##
## Author: kn6yuh
## Date:   20250802
## Description: This script can be used to perform all the necessary changes
##              for the montly rollover.
##

MEETINGS="meetings.md"
TEMPLATE="meetings/template.md"
SHORT_INCLUDE="_includes/meeting-short.md"
TEMPLATE_INCLUDE="_includes/meetings-template.md"
PARSE_PAST_MEETINGS="./script/parse-past-meetings.sh"

echo "[ Step 1: Detecting next monthly meeting filename ]"
if ls $MEETINGS >/dev/null
then
        echo "  File $MEETINGS exists."
else
        echo "  File $MEETINGS doesn't exist. Exiting."
        exit
fi

YEAR=`ls -al $MEETINGS | egrep -o "meetings/[0-9]{4}/[0-9]{6}\.md" | cut -d / -f 2`
MONTH=`ls -al $MEETINGS | egrep -o "meetings/[0-9]{4}/[0-9]{6}\.md" | cut -d / -f 3 | egrep -o "[0-9]{2}\.md" | cut -d "." -f1`

echo -n "  Current meeting is for year($YEAR) month($MONTH). "

if [ $MONTH -eq "12" ]; then
	YEAR=$((10#$YEAR+1))
	MONTH=01
else
	MONTH=$((10#$MONTH+1))
fi

printf "Next meeting is year(%d) month(%02d).\n" "$YEAR" "$MONTH"

NEXT_MEETING_FILE=`printf "meetings/%d/%d%02d.md" "$YEAR" "$YEAR" "$MONTH"`

echo "[ Step 2: Deleting the current symlink ]"
echo -n "  Running: 'rm $MEETINGS' ..."
if rm $MEETINGS
then
	echo "DONE"
else
	echo "FAILED"
	exit
fi

echo "[ Step 3: Creating the new meeting file if necessary ]"
if [ -s "$NEXT_MEETING_FILE" ]; then
  echo "  The file '$NEXT_MEETING_FILE' exists and is not empty."
else
  echo -n "  The file '$NEXT_MEETING_FILE' does not exist or is empty. "
  echo -n "Initializing from $TEMPLATE ... "
  if cp $TEMPLATE $NEXT_MEETING_FILE
  then
	  echo "DONE"
  else
	  echo "FAILED"
	  exit
  fi
fi

echo "[ Step 4: Creating the new symlink ]"
echo -n "  Running: 'ln -s $NEXT_MEETING_FILE $MEETINGS' ..."
if ln -s $NEXT_MEETING_FILE $MEETINGS
then
	echo "DONE"
else
	echo "FAILED"
	exit
fi

echo "[ Step 5: Updating the $SHORT_INCLUDE file ]"
DAY=`cal $MONTH $YEAR | awk '/Fr/{getline;if(NF==1){getline;}printf("%d\n",$(NF-1));}'`
MONTH_YEAR=`cal $MONTH $YEAR | head -1 | egrep -o "[a-zA-Z]+\s[0-9]+"`

echo """## Next club meeting
* **Date**: \`$DAY $MONTH_YEAR\`
* **Topic**: \`TBA\`
* **Presenter**: \`TBA\`
* **Zoom Meeting**:
   * <https://us02web.zoom.us/j/83692257191>
   * +16699006833,,83692257191# US (San Jose)

For more information, visit the [meetings page](/meetings.html).
""" > $SHORT_INCLUDE

if [[ $? -eq 0 ]];
then
	echo "  Successfully updated $SHORT_INCLUDE."
else
	echo "  Failed to update $SHORT_INCLUDE. Exiting!"
	exit
fi

echo "[ Step 6: Updating past meeting history ]"
echo -n "  Running $PARSE_PAST_MEETINGS ..."

if $PARSE_PAST_MEETINGS > $TEMPLATE_INCLUDE
then
	echo "DONE"
else
	echo "FAILED. Exiting!"
	exit
fi
