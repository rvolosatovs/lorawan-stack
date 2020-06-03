---
title: Configure ttnctl
weight: 15
description: Install and configure CLI tools
---

You will need to use the latest version of `ttnctl`, the CLI for the public community network (v2).
Follow the instructions from [the TTN docs][1]. An overview is given below:

Download `ttnctl` [for your operating system][2].

Update to the latest version:

```bash
$ ttnctl selfupdate
```

Go to [https://accounts.thethingsnetwork.com][3] and click [ttnctl access code][4].

Use the returned code to login from the CLI with:

```bash
$ ttnctl user login "t9XPTwJl6shYSJSJxQ1QdATbs4u32D4Ib813-fO9Xlk"
```

[1]: https://www.thethingsnetwork.org/docs/network/cli/quick-start.html
[2]: https://www.thethingsnetwork.org/docs/network/cli/quick-start.html#installation
[3]: https://accounts.thethingsnetwork.com
[4]: https://account.thethingsnetwork.org/users/authorize?client_id=ttnctl&redirect_uri=/oauth/callback/ttnctl&response_type=code

> **NOTE**: The commands above assume that ttnctl can be found in your PATH. If that
> is not the case, you need to specify the path to `ttnctl` instead, e.g. `./ttnctl`.
