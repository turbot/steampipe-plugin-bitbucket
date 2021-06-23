connection "bitbucket" {
  plugin = "bitbucket"

  # username = "YOUR_USERNAME"

  # The Bitbucket plugin uses a app password to authenticate
  # to the bitbucket APIs  (it looks like `bLHDMmvlk7L9wGerIfMg`).
  # You must create a App Password (https://developer.atlassian.com/bitbucket/api/2/reference/meta/authentication#app-pw)
  # and assign the following scopes:
  #  - `repo` (all)
  #  - `read:org`
  #  - `read:user`
  #  - `user:email`
  # password  = "YOUR_TOKEN_HERE"
}
