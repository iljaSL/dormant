![workflow badge](https://github.com/iljaSL/dormant/actions/workflows/ci.yml/badge.svg)

Dormant is a go.mod dependencies analyzing tool. Find out which of your used dependencies are actively, inactively or sporadically maintained.

# Table of content


- [Table of content](#table-of-content)
- [Usage](#usage)
- [Overview](#overview)
- [Changing Default Values](#changing-default-values)
- [Upcoming Features](#upcoming-features)
- [How To Contribute](#how-to-contribute)
- [Bugs](#bugs)
- [License](#license)

# Usage

CLI Overview

```
dormant --help
```

Inspect Dependencies inside a `go.mod` file:
```
dormant inspect go.mod
```
<p align="center">
  <img src="./assets/dormant.gif">
</p>

# Overview

Dormant is using GitHub's REST API, in particular this endpoint here, https://docs.github.com/en/rest/reference/commits, in order to retrieve the information needed.
GitHub does not require an authentication for this endpoint, but it comes with a rate limit in how many times you can call the API with your IP address, the rate limit allows you to make up to 10 requests per minute.

The authenticated requests feature will come in the near future.
This feature will also allow to inspect Dependencies on GitLab, which Dormant is not supporting at the moment.

# Changing Default Values

By default Dormant is set to determine an inactive dependency which has
not been updated for more than 6 months and an sporadic status for a dependency
which has been updated in a period between 4 and 6 months. Everything under 4 months
is actively maintained according to Dormant's default settings.
This default values can bee changed by creating a file called `.dormant.yaml`
inside the Home Directory. 

```
 ~ cat .dormant.yaml
inactivityDuration: 12
sporadicDuration: 6
```
Dormant is automatically parsing the Home Directory for the `.dormant.yaml` file,
a special command to use the env file is not needed.

# Upcoming Features

* Authenticated requests for GitHub
* Option for analyzing only the Direct Dependencies
* Dependencies Health Percentage
* Fancy HTML Report

# How To Contribute

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Added some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

# Bugs

If you experience any problems, please let me know by creating a new issue.

# License
Dormant is released under the MIT license. See [LICENSE](./LICENSE)
