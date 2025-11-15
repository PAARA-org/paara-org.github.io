#!/bin/bash

echo """## Notes on using zoom

---

**NOTE**: \`Please make sure you're running the latest version of ZOOM in order to connect.\`

**Zoom Tips:**
1. (Jim, K6SV) Download [Zoom PDF Guide](/meetings/ZoomGuide.pdf)
2. (Rob, KC6TYD) [How To Change Your Zoom Name Before Entering](http://www.youtube.com/watch?v=TJisWaSTYFg)

---

## Club Meetings

The Palo Alto Amateur Radio Association meets on the 1st Friday of the month at 7:00 pm (19:00 J) (except when the first Friday lands on a holiday and typically moved to the 2nd Friday of the month).

* **Location**: Room H6, Cubberley Community Center 4000 Middlefield Road, Palo Alto, CA 94303
* **Telephone**: (650) 329-2418

![cubberley-600.jpg](/meetings/cubberley-600.jpg)

After the meeting, many members gather for beer, pizza, and eyeball QSO at **Mountain Mike's Pizza**, located at **3918 Middlefield Rd**, in Palo Alto. It's literally next door to Cubberly, in the Charleston Shopping Center.

## Board Meetings

The PAARA Board meets on the 3rd Wednesday of every month at 7:00 pm by video conference.  We conduct most of our club business there in order to allow more time for guest speakers and socializing at regular club meetings.  Any PAARA member is welcome to attend.

A few meetings throughout the year have pre-defined agenda items. These include:
* **September Board Meeting**: The Nominating Committee is appointed by the President to decide on a slate of candidates to be recommended to the membership.
* **October General Meeting**: The Nominating Committee presents its nominations to the membership at the regular club meeting.  Nominations from club members opens at the beginning of the meeting, and closes at the end of the meeting.

## YouTube recordings

Around 2021 (during the COVID pandemic), we started recording our monthly presentations, which you can find either under the individual meetings below, or in our [YouTube channel](https://youtube.com/@paarapresentations-c3w?si=0nIFJeYS5yRwfU9A){:target=\"_blank\"}.

<iframe style=\"max-width: 480px\" width=\"100%\" height=\"360\" src=\"https://www.youtube.com/embed/videoseries?list=UUmIe9q2LiRcDk0swxNGfw6A\" frameborder="0" allowfullscreen></iframe>

## Past Meetings

"""

CURRENT_DATE="$(date +'%Y%m')"

# List all meeting years
find meetings/ -type d | sort -nr | grep -E "[0-9]$" | while read meet_year; do
  year=$(echo $meet_year | cut -d "/" -f2)
  if ! [[ $year =~ ^20[0-9]{2}$ ]]; then
    continue
  fi

  echo "### $year"
  echo ""
  find $meet_year -type f | grep -E ".md$" | sort -nr | while read meet_month; do
    if [ "$(basename $meet_month .md)" -le "$CURRENT_DATE" ]; then
      meet_month_html=$(echo $meet_month | sed 's/\.md/.html/')
      month=$(echo $meet_month | cut -d "/" -f 3 | cut -d "." -f 1)
      presenter=$(grep "\*\*Presenter\*\*:" $meet_month | cut -d ":" -f 2- | sed 's/`/**/g')
      topic=$(grep "\*\*Topic\*\*:" $meet_month | cut -d ":" -f 2-)
      has_youtube=$(grep "www.youtube.com" $meet_month)
      echo -n "* [$month](/$meet_month_html) "
      if [[ -n "$has_youtube" ]]; then
        echo "$presenter :$topic (<mark>with video</mark>)"
      else
        echo "$presenter :$topic"
      fi
    fi
  done
  echo ""
done
