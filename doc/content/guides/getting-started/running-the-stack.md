---
title: "Running The Things Stack"
description: ""
weight: 4
---

You can run The Things Stack using Docker or container orchestration solutions using Docker. An example [Docker Compose configuration](https://github.com/TheThingsNetwork/lorawan-stack/tree/master/docker-compose.yml) is available to get started quickly.

With the `docker-compose.yml` file in the directory of your terminal prompt, enter the following commands to initialize the database, create the first user `admin`, create the CLI OAuth client and start The Things Stack:

```bash
$ docker-compose pull
$ docker-compose run --rm stack is-db init
$ docker-compose run --rm stack is-db create-admin-user \
  --id admin \
  --email admin@localhost
$ docker-compose run --rm stack is-db create-oauth-client \
  --id cli \
  --name "Command Line Interface" \
  --owner admin \
  --no-secret \
  --redirect-uri 'local-callback' \
  --redirect-uri 'code'
$ docker-compose run --rm stack is-db create-oauth-client \
  --id console \
  --name "Console" \
  --owner admin \
  --secret console \
  --redirect-uri 'https://localhost:8885/console/oauth/callback' \
  --redirect-uri 'http://localhost:1885/console/oauth/callback' \
  --redirect-uri '/console/oauth/callback'
$ docker-compose up
```

With The Things Stack up and running, it's time to connect gateways, create devices and work with streaming data. See [Command-line Interface]({{< relref "cli" >}}) to proceed.
