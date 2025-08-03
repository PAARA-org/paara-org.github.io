# PAARA.org Website

This repository contains the PAARA.org website which is hosted on Github Pages.

Club members are encouraged to help file Issues for the problems encountered on the website, but at the same time we are more than happy to receive Pull Requests with fixes or improvements.

This page acts mainly as a HowTo for other contributors. The old readme is now at [ARCHITECT.md](ARCHITECT.md).

## How to add a new meeting?

1. Modify the meeting skeleton to add: 
   * topic
   * presenter
   * details
   * Zoom meeting details
   * example at <https://github.com/PAARA-org/paara-org.github.io/blob/master/meetings/2025/202506.md?plain=1>

2. Upload an image (optional) and render it inline
3. Modify `_includes/meeting-short.md` to include the same details (this will be rendered on the first page)
   * example at <https://github.com/PAARA-org/paara-org.github.io/blob/41e09eb62eccab82d7b528a0111cd3cf2c46dceb/_includes/meeting-short.md?plain=1>

### How to add Zoom meeting details

Add zoom meeting link and phone dial number to the `_includes/meeting-short.md` as in the example below::

```markdown
## Next club meeting
* **Date**: `6 June 2025`
* **Topic**: `The Legacy of Trans-Pacific Radio`
* **Presenter**: `Richard Dillman (W6AWO)`
* **Zoom Meeting**:
   * <https://us02web.zoom.us/j/81572645379>
   * <+16699006833,,81572645379#> US (San Jose)
```

## What to do after a monthly meeting ends? Rollover to a new month!

Run the `./script/rollover-to-next-month.sh` script from the top repository location which will generate an output similar to this:

```
% ./script/rollover-to-next-month.sh
[ Step 1: Detecting next monthly meeting filename ]
  File meetings.md exists.
  Current meeting is for year(2025) month(08). Next meeting is year(2025) month(09).
[ Step 2: Deleting the current symlink ]
  Running: 'rm meetings.md' ...DONE
[ Step 3: Creating the new meeting file if necessary ]
  The file 'meetings/2025/202509.md' does not exist or is empty. Initializing from meetings/template.md ... DONE
[ Step 4: Creating the new symlink ]
  Running: 'ln -s meetings/2025/202509.md meetings.md' ...DONE
[ Step 5: Updating the _includes/meeting-short.md file ]
  Successfully updated _includes/meeting-short.md.
[ Step 6: Updating past meeting history ]
  Running ./script/parse-past-meetings.sh ...DONE
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

## How to add the YouTube to a meeting

1. Open the YouTube video in your browser
2. Click the `Share` -> `Embed` and copy the **iframe** code:

```html
<iframe width="560" height="315" src="https://www.youtube.com/embed/5QrTYcUS4hU?si=ImhnyzMrJMUEigxn" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
```

3. Add this code into the particular monthly meeting file, just above the `## Presentation materials` line

4. Regenerate the past meetings index, which will add the video tag to this meeting as well:

```bash
./script/parse-past-meetings.sh > _includes/meetings-template.md
```
