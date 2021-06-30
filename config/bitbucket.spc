connection "bitbucket" {
  plugin = "bitbucket"

  # username = "YOUR_BITBUCKET_USERNAME"

  # The Bitbucket plugin uses a app password to authenticate
  # to the bitbucket APIs  (it looks like `bLHDMmvlk7L9wGerIfMg`).
  # You must create a App Password (https://bitbucket.org/account/settings/app-passwords/)
  # and assign the following scopes(https://developer.atlassian.com/bitbucket/api/2/reference/meta/authentication#scopes-bbc):
  # - workspace:read
  # - account:read
  # - read:user
  # - snippet:read
  # - wiki:read
  # - repository:read
  # - user:read
  # - pullrequest:read
  # - issue:read
  # password  = "YOUR_TOKEN_HERE"
}
