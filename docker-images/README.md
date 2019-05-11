# Sourcegraph derivative Docker images

The directory contains Sourcegraph docker images which are derivative of an existing Docker image, but with better defaults for our use cases. For example:

- `sourcegraph/alpine` handler setting up a `sourcegraph` user account, installing common packages.
- `sourcegraph/postgres` is `postgres` but with some Sourcegraph defaults.

If you are looking for our non-derivative Docker images, see e.g. `/cmd/.../Dockerfile` and `/enterprise/.../frontend/Dockerfile` instead.

### Building 

These images are not yet built on CI. To build one, you must sign in to our Docker hub and run `make <image name>` in this repository. For example:

```Makefile
make alpine
```

Before running the above command, you should have your changes reviewed and merged into `master`.

### Known issues

Many of our derivative images have not yet been moved here
