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
RAFFLE_MAIN="_includes/raffle.md"
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
NEXT_RAFFLE_FILE=`printf "_includes/raffle/%d%02d.md" "$YEAR" "$MONTH"`
NEXT_RAFFLE_INCLUDE=`printf "raffle/%d%02d.md" "$YEAR" "$MONTH"`

echo "[ Step 2: Deleting the current $MEETINGS symlink ]"
echo -n "  Running: 'rm $MEETINGS' ..."
if rm $MEETINGS
then
	echo "DONE"
else
	echo "FAILED"
	exit
fi

echo "[ Step 2.1: Deleting the current $RAFFLE_MAIN file ]"
echo -n "  Running: 'rm $RAFFLE_MAIN' ..."
if rm $RAFFLE_MAIN
then
        echo "DONE"
else
        echo "FAILED"
        exit
fi


echo "[ Step 3: Creating the new meeting file if necessary ]"
DAY=`cal $MONTH $YEAR | awk '/Fr/{getline;if(NF==1){getline;}printf("%d\n",$(NF-1));}'`
MONTH_YEAR=`cal $MONTH $YEAR | head -1 | egrep -o "[a-zA-Z]+\s[0-9]+"`
MONTH_STRING=`cal $MONTH $YEAR | awk '{print $1}'`
if [ -s "$NEXT_MEETING_FILE" ]; then
  echo "  The file '$NEXT_MEETING_FILE' exists and is not empty."
else
	echo """# ${MONTH_STRING} Monthly Meeting

* **Date**: \`$DAY $MONTH_YEAR\`
* **Time**: \`07:00 PM Pacific Time\`
* **Topic**: \`To be announced\`
* **Presenter**: \`To be announced\`
{% include zoom-details.md %}

## Details

## Presentation materials

## Raffle

{% include ${NEXT_RAFFLE_INCLUDE} %}

{% include meetings-template.md %}
""" > $NEXT_MEETING_FILE

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

echo """## Next club meeting
* **Date**: \`$DAY $MONTH_YEAR\`
* **Topic**: \`TBA\`
* **Presenter**: \`TBA\`
""" > $SHORT_INCLUDE

if [[ $? -eq 0 ]];
then
	echo "  Successfully updated $SHORT_INCLUDE."
else
	echo "  Failed to update $SHORT_INCLUDE. Exiting!"
	exit
fi

echo "[ Step 6: Create new raffle include file ]"
DAY=`cal $MONTH $YEAR | awk '/Fr/{getline;if(NF==1){getline;}printf("%d\n",$(NF-1));}'`
MONTH_YEAR=`cal $MONTH $YEAR | head -1 | egrep -o "[a-zA-Z]+\s[0-9]+"`

mkdir -p meetings/raffle/${YEAR}/
for i in {1..5} ; do cp meetings/raffle/raffle_prize.png meetings/raffle/${YEAR}/${YEAR}${MONTH}-${i}.png ; done

echo """
<details>
  <summary><b>Click the ticket to see the raffle prizes! <img src=\"/images/raffle-ticket.png\" alt=\"raffle-ticket\" width=\"90\"></b></summary>
  <table>
    <tr>
        <th>5th prize</th>
        <th>4th prize</th>
        <th>3rd prize</th>
        <th>2nd prize</th>
        <th>1st prize</th>
    </tr>
    <tr>
        <td><img src=\"/meetings/raffle/${YEAR}/${YEAR}${MONTH}-5.png\" alt=\"image\"></td>
        <td><img src=\"/meetings/raffle/${YEAR}/${YEAR}${MONTH}-4.png\" alt=\"image\"></td>
        <td><img src=\"/meetings/raffle/${YEAR}/${YEAR}${MONTH}-3.png\" alt=\"image\"></td>
        <td><img src=\"/meetings/raffle/${YEAR}/${YEAR}${MONTH}-2.png\" alt=\"image\"></td>
        <td><img src=\"/meetings/raffle/${YEAR}/${YEAR}${MONTH}-1.png\" alt=\"image\"></td>
    </tr>
    <tr>
        <td>5th Prize description</td>
        <td>4th Prize description</td>
        <td>3nd Prize description</td>
        <td>2nd Prize description</td>
        <td>1st Prize description</td>
    </tr>
  </table>
</details>
""" > ${NEXT_RAFFLE_FILE}


echo "[ Step 7: Create a copy of $NEXT_RAFFLE_FILE into $RAFFLE_MAIN ]"
echo -n "  Running: 'cp $NEXT_RAFFLE_FILE $RAFFLE_MAIN' ..."
if cp $NEXT_RAFFLE_FILE $RAFFLE_MAIN
then
        echo "DONE"
else
        echo "FAILED"
        exit
fi



echo "[ Step 8: Updating past meeting history ]"
echo -n "  Running $PARSE_PAST_MEETINGS ..."

if $PARSE_PAST_MEETINGS > $TEMPLATE_INCLUDE
then
	echo "DONE"
else
	echo "FAILED. Exiting!"
	exit
fi
