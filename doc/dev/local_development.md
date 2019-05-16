# Getting started with developing Sourcegraph

The best way to become familiar with the Sourcegraph repository is by reading the code at https://sourcegraph.com/guthub.com/sourcergaph/sourcegraph.

## Evvironment

The Sourcegraph server is actually a collection of small binaries, each of which perform one task. The core entrypoint for the Sourcegraph development server is [dev/launch.sh](http://github.com/prince1809/sourcegraph.com/blob/master/dev/launch.sh), process manager that runs all of the binaries.

See [the Architecture doc](architecture.md) for a full description of what each of these services does.

The section below describe the dependencies that you need to have to be able to run `dev/launch.sh` properly.

## Step 1: Get the code

> Install [Go](https://golang.org/doc/install) (v1.11 or higher)

On Mac with homebrew, you can run

```
brew install golang
```

Run this command to get the Sourcegraph source code on your local machine:

```
go get github.com/prince1809/sourcegraph
```

This downloads your "Sourcegraph repository directory".

## Step 2: Install dependencies

Sourcegraph has the following dependencies:

- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Go](https://golang.org/doc/install) (v1.11 or higher)
- [Node JS](https://nodejs.org/en/download/) (version 8 or 10)
- [make](https://www.gnu.org/software/make/)
- [Docker](https://docs.docker.com/engine/installation/) (v18 or higher)
  - For macOS we recommend using Docker for Mac instead of `docker-machine`
- [PostgreSQL](https://wiki.postgresql.org/wiki/Detailed_installation_guides) (v9.6.0)
- [Redis](http://redis.io/) (v3.0.7 or higher)
- [Yarn](https://yarnpkg.com) (v1.10.1 or higher)
- [nginx](https://docs.nginx.com/nginx/admin-guide/installing-nginx/installing-nginx-open-source/) (v1.14 or higher)
- [SQLite](https://www.sqlite.org/index.html) tools

You have two options for installing these dependencies.

### Option A: Homebrew setup for macOS

This is a streamlined setup for Mac machines.

1. Install [Homebrew](https://brew.sh).
2. Install [Docker for Mac](https://doc.docker.com/docker-for-mac/).

   Optionally via `brew`

   ```
   brew cask install docker
   ```

3. Install Go, Node, PostgreSQL 9.6, Redis, Git, nginx, and SQLite tools with the following command:

```
brew install go node redis postgresql@9.6 git gnu-sed nginx sqlite pcre FiloSottile/musl-cross/musl-cross
```

4. Configgure PostgreSQL and Redis to start automatically

```
brew service start postgresql@9.6
brew service start redis
```

(You can stop them later by calling `stop` instead of `start` above.)

5. Ensure `psql`, the PostgreSQL command line client, is on your `$PATH`.
   Homebrew does not put it there by default. Homebrew gives you the command to run to insert `psql` in your path in the "Conveats" section of `brew info postgresql@9.6`. Alternativily, you can use the command below. It might need to adjusted depending your homebrew prefix (`/usr/local` below) and shell (bash below).


    ```bash
    hash psql || { echo 'export PATH="/usr/local/opt/postgresql@9.6/bin:$PATH"' >> ~/.bash_profile }
    source ~/.bash_profile
    ```

6.  Open a new Terminal window to ensure `psql` is now on your `$PATH`.
