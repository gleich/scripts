# This program is intended to update homebrew, npm, mas, flutter, and
# upload the configs for my dot files repo

Command() {
    COLUMNS=$(tput cols)
    printf %"$COLUMNS"s | tr " " "━"
    echo "$1" | fmt -c -w $COLUMNS
    printf %"$COLUMNS"s | tr " " "━"
    $1
}

Command brew\ update
Command brew\ upgrade
Command brew\ upgrade\ --cask
Command brew\ cleanup
Command mas\ upgrade
Command npm\ upgrade\ -g
Command fgh\ update
Command fgh\ clean
Command flutter\ upgrade
cd /Users/mattgleich/github/Matt-Gleich/public/Shell/dots/bin
Command sh\ push-latest.sh
