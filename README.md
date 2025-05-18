# PAARA.org Website

This repository contains the PAARA.org website which is hosted on Github Pages.

Club members are encouraged to help file Issues for the problems encountered on the website, but at the same time we are more than happy to receive Pull Requests with fixes or improvements.

This page acts mainly as a HowTo for other contributors. The old readme is now at [ARCHITECT.md](ARCHITECT.md).

## What to do after a monthly meeting ends?

Examples are for changes done after the May 2025 Monthly meeting

1. Update the symlink for /meetings.md from `meetings/2025/202505.md` (May) to `meetings/2025/202506.md` (June)

```bash
rm -f meetings.md
ln -s meetings/2025/202506.md meetings.md
```

This will make the June meeting to be loaded when someone looks at the Monthly Meetings page.

2. Update the short meeting description (`_includes/meeting-short.md`), which is included on the main page:

* from

```md
## Next club meeting
* **Date**: `2 May 2025`
* **Topic**: `The software defined radio, and why it belongs in your shack`
* **Presenter**: `Carlos Felix, K9OL`

For more information, visit the [meetings page](/meetings.html).
```

* to

```md
## Next club meeting
* **Date**: `6 June 2025`
* **Topic**: `TBA`
* **Presenter**: `TBA`

For more information, visit the [meetings page](/meetings.html).
```

## What to do when we know the topic and presenter?

Update `meetings/2025/202506.md` and `_includes/meeting-short.md` with the topic and presenter details.

## How to add presentation materials to a meeting?

For some of the meetings, we get a PDF of PPT that we can link in the meeting page (for example <https://paara.org/meetings/2025/202504.html>).

The following steps must be followed:
1. Upload the PDF or PPT into the <https://github.com/PAARA-org/paara.org-attachments> repository
   * the files must be uploaded in meetings/`YYYY`, where `YYYY` is the meeting year (meetings/`2025` for example)
   * the file name must start with YYYY`MM`, where `MM` is the month (2025`04` for example)
   * a full filename example is `meetings/2025/202501_Instant_Tuning_For_A_Manual_Tuner.pdf` or `meetings/2025/202504_Stealth_HF_Antennas_McClary.ppt`
2. In the meeting page, under the `## Presentation materials` heading, add links to the PDF in both Download and Docs Viewer mode (example at <https://github.com/PAARA-org/paara-org.github.io/blob/77a3ad09ce3f91de06af4f51ec856ad1e66e5321/meetings/2025/202501.md?plain=1#L17>)

```md
* *PDF*: [[download](https://github.com/PAARA-org/paara.org-attachments/blob/main/events/2025/202501_Instant_Tuning_For_A_Manual_Tuner.pdf?raw=true)] [[Docs Viewer](https://docs.google.com/viewer?url=https://github.com/PAARA-org/paara.org-attachments/blob/main/events/2025/202501_Instant_Tuning_For_A_Manual_Tuner.pdf?raw=true)]
```
