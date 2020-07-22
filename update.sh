# This program is intended to update homebrew, npm, and
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
Command brew\ cask\ upgrade
Command brew\ cleanup
Command npm\ upgrade\ -g
cd /Users/matthewgleich/Documents/GitHub/Personal/Bash/Dot-Files/
Command sh\ push-latest.sh
