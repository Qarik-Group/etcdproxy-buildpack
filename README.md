# Buildpack to add Etcd proxy as sidecar to Cloud Foundry applications

During deployment staging you will see the latest etcd version being installed into the application droplet:

```plain
$ cf v3-push
...
[STG/0] OUT -----> Etcdproxy Buildpack version 1.0.0
[STG/0] OUT -----> Supplying etcd
[STG/0] OUT        Using etcd version 3.3.13
[STG/0] OUT -----> Installing etcd 3.3.13
[STG/0] OUT        Copy [/tmp/buildpacks/3d4ba95f26291321f5c6633264c1c6bb/dependencies/02866fcde1388c50c101a96ed2d1cfdb/etcd-v3.3.13-linux-amd64.ta
```

```plain
$ cf logs app-with-etcd
...
[APP/PROC/WEB/SIDECAR/ETCDPROXY/0] ERR | etcdmain: listening for grpc-proxy client requests on 127.0.0.1:2379
```

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
    cf v3-apply-manifest -f fixtures/rubyapp/manifest.cfdev.yml
    cf v3-push app-using-etcd -p fixtures/rubyapp
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

1. Run integration tests

    ```bash
    ./scripts/integration.sh
    ```

    More information can be found on Github [cutlass](https://github.com/cloudfoundry/libbuildpack/cutlass).

### Reporting Issues

Open an issue on this project
