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
Command "rustup update"
Command "cargo install-update -a"
Command "echo checking docker status"
if (docker stats --no-stream); then
    Command "docker system prune -af"
fi
cd ~/src/dots
Command "fetch"
