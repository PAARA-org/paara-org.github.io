# PAARA.org Website

This repository contains the PAARA.org website which is hosted on Github Pages.

Club members are encouraged to help file Issues for the problems encountered on the website, but at the same time we are more than happy to receive Pull Requests with fixes or improvements.

This page acts mainly as a HowTo for other contributors. The old readme is now at [ARCHITECT.md](ARCHITECT.md).

## What to do after a monthly meeting ends?

Examples are for changes done after the May 2025 Monthly meeting

1. Update the symlink for /meetings.md from `meetings/2025/202505.md` (May) to `meetings/2025/202506.md` (June)

```bash
$ rm -f meetings.md
$ ln -s meetings/2025/202506.md meetings.md
```

This will make the June meeting to be loaded when someone looks at the Monthly Meetings page.

2. Update the short meeting description (`_includes/meeting-short.md`), which is included on the main page:

* from

```
## Next club meeting
* **Date**: `2 May 2025`
* **Topic**: `The software defined radio, and why it belongs in your shack`
* **Presenter**: `Carlos Felix, K9OL`

For more information, visit the [meetings page](/meetings.html).
```

* to

```
## Next club meeting
* **Date**: `6 June 2025`
* **Topic**: `TBA`
* **Presenter**: `TBA`

For more information, visit the [meetings page](/meetings.html).
```

## What to do when we know the topic and presenter?

Update `meetings/2025/202506.md` and `_includes/meeting-short.md` with the topic and presenter details.

