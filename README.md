# <a href="https://sourcegraph.com"><img  alt="Sourcegraph" src="https://storage.googleapis.com/sourcegraph-assets/sourcegraph-logo.png" height="32px"/></a>

[![build](https://badge.buildkite.com/00bbe6fa9986c78b8e8591cffeb0b0f2e8c4bb610d7e339ff6.svg?branch=master)](https://buildkite.com/prince1809/sourcegraph)[![apache license](https://img.shields.io/badge/license-Apache-blue.svg)](LICENSE)

[Sourcegraph](https://about.sourcegraph.com/) is a fast, open-source, fully featured code search and navigation engine.

![Screenshot](https://user-images.githubusercontent.com/1646931/46309383-09ba9800-c571-11e8-8ee4-1a2ec32072f2.png)

**Features**

- Fast global code search with a hybrid backend that combines a trigram index with in-memory streaming
- Code intelligence for many language via the [Language Server Protocol](https://langserver.com)
- Enhances Github, Gitlab, Phabricator, and other code hosts and code review tools via the [Sourcegraph browser extension](https://docs.sourcegraph.com/integration/browser_extension)
- Integration with third-party developer tools via the [Sourcegraph extension API](https://docs.sourcegraph.com/extensions)

## Try it yourself

- Try out the public instance on any open-source repository at [sourcegraph.com](https://sourcegraph.com/github.com/golang/go/-/blob/src/net/http/httptest.go#L41:6&tab=references).
- Install the free and open-source [browser extension](http://chrome.google.com/webstore/detail/sourcegraph/dgjhfomjieaadpoljlnidmbgkdffpack?hl=en).
- Visit [about.sourcegraph.com](https://about.sourcegraph.com) for more information about product features.

## Development

### Prerequisites

- Git
- Go (1.11 or later)
- Docker
- PostgreSQL (verison 9)
- Node.js (version 8 or 10)
- Redis
- Yarn
- Nginx

For a detailed guide to installing prerequisites, see [themse instructions](doc/dev/local_development.md#step-2-install-dependencies).

### Installation

> Prebuilt Docker images are the fastest way to use Sourcegraph. See the [quickstart installation guide](https://docs.sourcegraph.com/#quickstart).
