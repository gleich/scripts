set sourceMatches to every playlist whose name is sourceName
if (count of sourceMatches) is not 1 then
	error "Expected exactly one source playlist named " & sourceName
end if
set sourcePlaylist to item 1 of sourceMatches
if special kind of sourcePlaylist is folder then
	error "Source playlist cannot be a folder: " & sourceName
end if
set sourceID to persistent ID of sourcePlaylist
set sourceIDs to persistent ID of every track of sourcePlaylist

set destinationMatches to every playlist whose name is destinationName
if (count of destinationMatches) > 1 then
	error "Multiple playlists named " & destinationName
end if
set destinationID to ""
set destinationIDs to {}
if (count of destinationMatches) is 1 then
	set destinationPlaylist to item 1 of destinationMatches
	if class of destinationPlaylist is not user playlist then
		error "Destination must be an ordinary user playlist: " & destinationName
	end if
	if smart of destinationPlaylist or genius of destinationPlaylist or special kind of destinationPlaylist is not none then
		error "Destination must be an ordinary user playlist: " & destinationName
	end if
	set destinationID to persistent ID of destinationPlaylist
	set destinationIDs to persistent ID of every track of destinationPlaylist
end if

set AppleScript's text item delimiters to ","
set sourceText to sourceIDs as text
set destinationText to destinationIDs as text
set AppleScript's text item delimiters to ""
set snapshot to sourceID & "|" & sourceText & "|" & destinationID & "|" & destinationText
