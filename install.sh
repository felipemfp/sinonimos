#!/bin/bash

get_latest_release() {
  # https://gist.github.com/lukechilds/a83e1d7127b78fef38c2914c4ececc3c
  curl --silent "https://api.github.com/repos/$1/releases/latest" | # Get latest release from GitHub api
    grep '"tag_name":' |                                            # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                    # Pluck JSON value
}

LATEST=$(get_latest_release felipemfp/sinonimos-cli)

sudo curl --silent -L "https://github.com/felipemfp/sinonimos-cli/releases/download/$LATEST/sinonimos" -o /usr/local/bin/sinonimos

sudo chmod +x /usr/local/bin/sinonimos

sinonimos -v