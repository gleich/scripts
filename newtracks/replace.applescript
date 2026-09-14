if snapshot is not expectedSnapshot then
	error "Playlists changed while reading them; rerun to retry"
end if

set desiredTracks to {}
repeat with trackID in desiredIDs
	set end of desiredTracks to first track of sourcePlaylist whose persistent ID is (contents of trackID)
end repeat

if destinationID is "" then
	set destinationPlaylist to make new user playlist with properties {name:destinationName}
	move destinationPlaylist to destinationFolder
end if
delete every track of destinationPlaylist
repeat with selectedTrack in desiredTracks
	duplicate (contents of selectedTrack) to destinationPlaylist
end repeat

set actualIDs to persistent ID of every track of destinationPlaylist
set AppleScript's text item delimiters to ","
set actualText to actualIDs as text
set AppleScript's text item delimiters to ""
return actualText
