#!/bin/bash

# List all meeting years
find meetings/ -type d | sort -nr | egrep "[0-9]$" | while read meet_year ; do
	year=`echo $meet_year | cut -d "/" -f2`
	echo "* $year"
	find $meet_year -type f | egrep ".md$" | sort -nr | while read meet_month ; do
		month=`echo $meet_month | cut -d "/" -f 3 | cut -d "." -f 1`
		presenter=`grep "\*\*Presenter\*\*:" $meet_month | cut -d ":" -f 2-`
		topic=`grep "\*\*Topic\*\*:" $meet_month | cut -d ":" -f 2-`
		has_youtube=`grep "www.youtube.com" $meet_month`
		if [[ -n "$has_youtube" ]]; then
			echo "   * [$month]($meet_month) - [with YouTube recording]"
		else
			echo "   * [$month]($meet_month)"
		fi
		echo "     * $presenter :$topic"
	done
done
