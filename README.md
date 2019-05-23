# Add etcd proxy as sidecar to Cloud Foundry applications

You will know that [etcd](https://coreos.com/etcd/) is a distributed, reliable key-value store. This buildpack does not install or run an etcd cluster, instead to runs an etcd proxy as a sidecar next to your application. Your application need only connect to `localhost2379` and get the full benefits of connecting to a secure, highly-available, multi-node etcd cluster.

From [etcd proxy docs](https://coreos.com/etcd/docs/latest/op-guide/grpc_proxy.html):

> The gRPC proxy is a stateless etcd reverse proxy operating at the gRPC layer (L7). The proxy is designed to reduce the total processing load on the core etcd cluster. For horizontal scalability, it coalesces watch and lease API requests. To protect the cluster against abusive clients, it caches key range requests.
>
> The gRPC proxy supports multiple etcd server endpoints. When the proxy starts, it randomly picks one etcd server endpoint to use. This endpoint serves all requests until the proxy detects an endpoint failure. If the gRPC proxy detects an endpoint failure, it switches to a different endpoint, if available, to hide failures from its clients. Other retry policies, such as weighted round-robin, may be supported in the future.

The etcd proxy is distributed within the `etcd` binary, which is packaged and installed into your Cloud Foundry application using this supply buildpack.

You configure your Cloud Foundry application to both:

1. Use this supply buildpack, in addition to your application's primary buildpack.
1. Run `etcd proxy start` as a sidecar within your application containers.

## Example usage

There is a sample application within `fixtures/rubyapp` that can be installed. To use sidecars we need the `cf v3-xyz` commands:

```plain
cf create-buildpack etcdproxy_buildpack etcdproxy_buildpack-*.zip 1
cf v3-create-app app-using-etcd
cf v3-apply-manifest -f fixtures/rubyapp/manifest.cfdev.yml
cf v3-push app-using-etcd -p fixtures/rubyapp
```

During deployment staging of `cf v3-push` you will see the latest etcd version being installed into the application droplet:


```plain
$ cf v3-push
...
[STG/0] OUT -----> Etcdproxy Buildpack version 1.0.0
[STG/0] OUT -----> Supplying etcd
[STG/0] OUT        Using etcd version 3.3.13
[STG/0] OUT -----> Installing etcd 3.3.13
[STG/0] OUT        Copy [/tmp/buildpacks/3d4ba95f26291321f5c6633264c1c6bb/dependencies/02866fcde1388c50c101a96ed2d1cfdb/etcd-v3.3.13-linux-amd64.ta
```

The application container logs will show etcd proxy output with `ETCDPROXY/0` prefix:

```plain
$ cf logs app-with-etcd
...
[APP/PROC/WEB/SIDECAR/ETCDPROXY/0] ERR | etcdmain: listening for grpc-proxy client requests on 127.0.0.1:2379
```

Your application only needs to connect to `localhost:2379` or `127.0.0.1:2379` without any other configuration requirements.

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
