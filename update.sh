# This program is intended to update homebrew, npm, mas, flutter, and
# upload the configs for my dot files repo. I run this manually every
# morning and every night.

Command() {
    COLUMNS=$(tput cols)
    printf %"$COLUMNS"s | tr " " "━"
    echo "$1" | fmt -c -w $COLUMNS
    printf %"$COLUMNS"s | tr " " "━"
    $1
}

Command brew\ update
Command brew\ upgrade
Command brew\ cleanup\ -s
Command mas\ upgrade
Command npm\ upgrade\ -g
Command fgh\ update
Command fgh\ clean
Command rustup\ update
cd /Users/mattgleich/github/Matt-Gleich/public/shell/dots/bin/fetch
Command poetry\ run\ python3\ main.py
