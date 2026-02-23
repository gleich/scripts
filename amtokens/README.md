# Apple Token Refresh Script

Every 6 months the Apple Developer Token expires and so with it the Apple Music User token. This script and HTML file serve as a way to refresh these tokens.

1. Go to [Apple Developer Page for Keys](https://developer.apple.com/account/resources/authkeys/list) and revoke the old key.
2. Create a new key with MusicKit enabled under media services.
3. Save the downloaded key into this directory as a `key.p8` file.
4. Run the script and save the tokens
