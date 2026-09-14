set folderMatches to every folder playlist whose name is folderName
if (count of folderMatches) > 1 then
	error "Multiple playlist folders named " & folderName
end if
if (count of folderMatches) is 0 then
	set destinationFolder to make new folder playlist with properties {name:folderName}
else
	set destinationFolder to item 1 of folderMatches
end if

if destinationID is not "" then
	set currentParent to missing value
	try
		set currentParent to parent of destinationPlaylist
	on error messageText number errorNumber
		if errorNumber is not -1728 then error messageText number errorNumber
	end try
	if currentParent is missing value then
		move destinationPlaylist to destinationFolder
	else if persistent ID of currentParent is not persistent ID of destinationFolder then
		move destinationPlaylist to destinationFolder
	end if
end if
