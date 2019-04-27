# Sourcegraph documentation

[Sourcegaph](https://sourcegraph.com) is a web-based. open-source, self-hosted code search and navigation tool for developers, used by Uber, Lyft, Yelp, and more.

## Quickstart guide

> NOTE: If you get stuck or need help, [file an issue](https://github.com/prince1809/sourcegraph/issues/new?&title=Improve+quickstart+guide)

It takes less than 5 minutes to install Sourcegraph using Docker. If you've got [Docker installed](http://docs.docker.com/engine/installations), you're ready to start the server which listens on port `7080` by default.

<!--
DO NOT CHANGE THIS TO A CODEBLOCK
We want line breaks for readabiliy, but backslashes to escape them do not work cross-platform. This uses line breaks that are rendered but not copy-pasted to the clipboard.
-->
<pre class="pre-wrap"><code>docker run<span class="virtual-br"></span> --publish 7080:7080 --publish 2633:2633 --rm<span class="virtual-br"></span> --volume ~/.sourcegraph/config:/etc/sourcegraph<span class="virtual-br"></span> --volume ~/.sourcegraph/data:/var/opt/sourcegraph<span class="virtual-br"></span> sourcegraph/server:3.3.2</code></pre>