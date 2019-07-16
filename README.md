# watchinator
Send a file to a UDP address any time it updates

## Example

Send `some_file.txt` to IPv6 link-local anycast on `en0` any time it changes

    ./watchinator some_file.txt [ff02::1%en0]:6291
