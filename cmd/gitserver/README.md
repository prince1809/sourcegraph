# gitserver

gitserver exposes an "exec" API over HTTP for running git commands against
clones of repositories. gitserver also exposes APIs for the management of
clones.

The management of clones comprises most of the complexity in gitserver since:
