#!/usr/bin/env python3
"""
Generate newsletter.md content by querying the PAARA-org/PAARAgraphs GitHub repository.
This script uses only Python standard library (no external dependencies).

Usage:
    export GITHUB_TOKEN=your_token_here
    ./generate-newsletter.py > newsletter.md
"""

import re
import json
import os
import subprocess
import sys
import time
import urllib.request
from pathlib import Path
from datetime import datetime

# GitHub API configuration
GITHUB_API_URL = "https://api.github.com/repos/PAARA-org/PAARAgraphs/contents"
GITHUB_RAW_BASE = "https://github.com/PAARA-org/PAARAgraphs/blob/main"

# Month configuration
MONTHS = [
    ("JAN", "01", "jan"),
    ("FEB", "02", "feb"),
    ("MAR", "03", "mar"),
    ("APR", "04", "apr"),
    ("MAY", "05", "may"),
    ("JUN", "06", "jun"),
    ("JUL", "07", "jul"),
    ("AUG", "08", "aug"),
    ("SEP", "09", "sep"),
    ("OCT", "10", "oct"),
    ("NOV", "11", "nov"),
    ("DEC", "12", "dec"),
]

# Year range (from earliest to latest)
START_YEAR = 1993
END_YEAR = 2025


def fetch_all_years(token: str) -> list[int]:
    req = urllib.request.Request(GITHUB_API_URL)
    req.add_header("Accept", "application/vnd.github.v3+json")
    req.add_header("User-Agent", "PAARA-Newsletter-Generator/1.0")
    req.add_header("Authorization", f"token {token}")

    YEAR_RE = re.compile(r"\d{4}$")

    with urllib.request.urlopen(req, timeout=10) as response:
        data = json.loads(response.read().decode("utf-8"))

        return sorted(int(r["name"]) for r in data if YEAR_RE.match(r["name"]))


def fetch_year_contents(token: str, year: int) -> set[str]:
    """
    Fetch the list of files in a given year's directory from GitHub.
    Returns a set of filenames that exist.
    """
    url = f"{GITHUB_API_URL}/{year}"

    try:
        req = urllib.request.Request(url)
        req.add_header("Accept", "application/vnd.github.v3+json")
        req.add_header("User-Agent", "PAARA-Newsletter-Generator/1.0")
        req.add_header("Authorization", f"token {token}")

        with urllib.request.urlopen(req, timeout=10) as response:
            data = json.loads(response.read().decode("utf-8"))

            # Extract filenames
            filenames = {item["name"] for item in data if item["type"] == "file"}
            return filenames
    except urllib.error.HTTPError as e:
        if e.code == 404:
            # Year directory doesn't exist
            return set()

        raise


def generate_newsletter_cell(
    current_year_month: str,
    year: int,
    month_abbr: str,
    month_num: str,
    month_img: str,
    available_files: set[str],
) -> str:
    """
    Generate a single table cell for a newsletter.
    Returns either a linked image or "-" if the newsletter doesn't exist.
    """
    # Determine the year suffix (last 2 digits)
    year_suffix = str(year)[-2:]

    # Expected filename
    filename = f"graph{month_num}{year_suffix}.pdf"

    if f"{year}{month_num}" >= current_year_month:
        return ""

    if filename in available_files:
        # Newsletter exists - create the linked image
        img_url = f"/images/newsletter/{month_img}.png"
        pdf_url = f"{GITHUB_RAW_BASE}/{year}/{filename}?raw=true"

        return (
            f'[<img src="{img_url}" width="25" height="25" alt="{month_abbr}">]'
            f'(https://docs.google.com/viewer?url={pdf_url}){{:target="_blank"}}'
        )

    # Newsletter doesn't exist
    return "-"


def generate_table_row(
    current_year_month: str, year: int, newsletters_by_year: dict[int, set[str]]
) -> str:
    """
    Generate a complete table row for a given year.
    """
    available_files = newsletters_by_year.get(year, set())

    cells = [f" | {year} |"]

    for month_abbr, month_num, month_img in MONTHS:
        cell = generate_newsletter_cell(
            current_year_month, year, month_abbr, month_num, month_img, available_files
        )
        cells.append(f" {cell} |")

    return "".join(cells)


def main():
    """Main entry point."""
    # Fetch available newsletters for each year
    github_token = os.environ.get("GITHUB_TOKEN")
    if github_token:
        print(
            "Fetching newsletter data from GitHub (using GITHUB_TOKEN)...",
            file=sys.stderr,
        )
    else:
        print(f"Usage: GITHUB_TOKEN=<token> {sys.argv[0]}", file=sys.stderr)
        sys.exit(1)

    newsletters_by_year: dict[int, set[str]] = {}

    years = fetch_all_years(github_token)

    print("Generating from {} to {}...".format(years[0], years[-1]), file=sys.stderr)

    for year in years:
        print(f"  Checking year {year}...", file=sys.stderr)
        newsletters_by_year[year] = fetch_year_contents(github_token, year)

    print("Generating markdown...", file=sys.stderr)

    # Build the markdown content
    lines = [
        "# PAARAgraphs Newsletter",
        "",
        "PAARAgraphs issues older than 6 months are added to the archive for all to read.",
        "",
        "[PAARA members](/membership.html) receive current issues in print or email format.",
        "",
        "## Archive",
        "",
    ]

    # Add table header
    header = " |      | JAN | FEB | MAR | APR | MAY | JUN | JUL | AUG | SEP | OCT | NOV | DEC |"
    separator = " |------|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|-----|"
    lines.append(header)
    lines.append(separator)

    current_year_month = datetime.now().strftime("%Y%m")

    # Add table rows (from newest to oldest)
    for year in reversed(years):
        row = generate_table_row(current_year_month, year, newsletters_by_year)
        lines.append(row)

    lines.extend(
        [
            "",
            "## Staff",
            "",
            "| Position | Name | Call Sign | Email |",
            "| --- | --- | --- | --- |",
            "| Editorial Board     | Bob Van Tuyl      | K6RWY | <rrvt@swde.com> |",
            "| Editorial Board     | Jim Thielemann    | K6SV  | <THIELEM@pacbell.net> |",
            "| Editorial Board     | Bob Ridenour      | KN6YGN| <kn6ygn@paara.org> |",
            "| Editorial Board     | Doug Teter        | KG6LWE| <kg6lwe@paara.org> |",
            "| Editor - Primary    | Bob Van Tuyl      | K6RWY | <rrvt@swde.com> |",
            "| Editor - Backup     | Jim Thielemann    | K6SV  | <THIELEM@pacbell.net> |",
            "| Advertising         | Walt Gyger        | W6WGY | <walt@tradewindsaviation.com> |",
            "| Member Profiles     | __VACANT__        |       | |",
            "| Technical Tips      | Ric Hulett        | N6AJS | <n6ajs@paara.org> |",
            "| Photographer        | __VACANT__        |       | |",
        ]
    )

    print("\n".join(lines))
    print("Done!", file=sys.stderr)


if __name__ == "__main__":
    main()
