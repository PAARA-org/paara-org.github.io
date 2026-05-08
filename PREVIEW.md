# Previewing the site locally

## Prerequisites

The following packages are required (correct for Ubuntu 24.04.4 — your distro may differ):

```sh
sudo apt install build-essential ruby3.2-dev
```

`ruby3.2-dev` pulls in `ruby3.2`. `build-essential` is needed to compile native gem extensions.

## Usage

```sh
./preview.sh
```

Open http://localhost:4000. Jekyll auto-regenerates on file changes — just save and refresh.
