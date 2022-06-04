# This program is intended to update homebrew, npm, mas, and
# upload the configs for my dot files repo. I run this manually every
# morning.

Command() {
    COLUMNS=$(tput cols)
    printf %"$COLUMNS"s | tr " " "━"
    echo "$1" | fmt -c -w $COLUMNS
    printf %"$COLUMNS"s | tr " " "━"
    $1
}

Command "brew update"
Command "brew upgrade"
Command "brew cleanup -s"
cd $(fgh ls senior)
Command "optic check"
Command "rustup update"
Command "cargo install-update -a"
cd $(fgh ls dots)
Command "fetch"
