# backups

Moves downloaded backups into the exports folder in iCloud. After each successful
move, checks off the matching checklist item in the open **Backups** task in
Things 3's **Today** list.

Enable Things URLs in **Things → Settings → General → Enable Things URLs → Manage**
and save its authorization token in macOS Keychain. In Keychain Access, create a
new password item in your login keychain with these values:

- Keychain Item Name: `go.mattglei.ch.scripts.backups.things`
- Account Name: `things-auth-token`
- Password: the Things authorization token

The command reads this token automatically. `THINGS_AUTH_TOKEN`, if set, overrides
the saved token. The token is never stored in the repository.
macOS may also ask you to allow your terminal to control Things.

Checklist titles match the backup names, except **Yamaha N800A** matches
**Yamaha R-N800A** and **Uniden R8** matches **Uniden**. Add a **Goodnotes** checklist
item if you want that backup checked off too.

The command finds today's task each time, so recurring tasks work without a
hardcoded ID. It reads the checklist from Things' local SQLite database in
read-only mode and submits the updated checklist through the
[Things URL API](https://culturedcode.com/things/support/articles/2803573/).
The API replaces the checklist; the command preserves its titles, order, and
completed/canceled states. Avoid editing the checklist while backups is running.
This depends on Things' internal database schema.

Already checked items require no update. Missing or duplicate items, missing
configuration, and Things errors are reported without stopping subsequent backups.
The parent Backups task remains open. Completion is confirmed by reading the
checklist back after the update.

Run with `go run ./backups` or install with `go install ./backups`.
