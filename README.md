## Buildpack User Documentation

### Building the Buildpack

To build this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

    ```bash
    source .envrc
    ```

    To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager

    ```bash
    ./scripts/install_tools.sh
    ```

1. Build the buildpack

    ```bash
    buildpack-packager build -stack cflinuxfs3 -cached
    ```

1. Use in Cloud Foundry

    Upload the buildpack to your Cloud Foundry.

    ```bash
    cf create-buildpack etcdproxy_buildpack etcdproxy_buildpack-*.zip 1
    cf v3-create-app app-using-etcd
    cf v3-apply-manifest -f fixtures/static/manifest.cfdev.yml
    cf v3-push app-using-etcd -p fixtures/static
    ```

    As buildpack that delivers a sidecar, you'll need an explicit `manifest.yml` that describes the startup of the sidecar. For example, running an app within CFDev, and proxying to an etcd server on your host machine (`host.cfdev.sh:2379`):

    ```yaml
    applications:
    - name: app-using-etcd
      instances: 1
      memory: 128M
      disk_quota: 512M
      stack: cflinuxfs3
      buildpacks:
      - etcdproxy_buildpack
      - staticfile_buildpack
      sidecars:
      - name: etcdproxy
        process_types:
        - web
        command: etcd grpc-proxy start --endpoints=host.cfdev.sh:2379 --listen-addr=127.0.0.1:2379
      routes:
      - route: app-using-etcd.dev.cfdev.sh
    ```

### Testing

Buildpacks use the [Cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass) framework for running integration tests.

To test this buildpack, run the following command from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

    ```bash
    source .envrc
    ```

    To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Run unit tests

    ```bash
    ./scripts/unit.sh
    ```

1. Run integration tests

    ```bash
    ./scripts/integration.sh
    ```

    More information can be found on Github [cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass).

### Reporting Issues

Open an issue on this project

## Disclaimer

This buildpack is experimental and not yet intended for production use.
