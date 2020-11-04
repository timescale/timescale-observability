# tobs - The Observability Stack for Kubernetes

Tobs is a tool that aims to make it as easy as possible to install a full observability
stack into a Kubernetes cluster. Currently this stack includes:
 * [Prometheus](https://github.com/prometheus/prometheus) to collect metrics
 * [Grafana](https://github.com/grafana/grafana) to visualize what's going on
 * [Promscale](https://github.com/timescale/promscale) ([design doc][design-doc]) to store metrics for the long-term and allow analysis with both PromQL and SQL.

 We plan to expand this stack over time and welcome contributions.

Tobs provides a CLI tool to make deployment and operations easier. We also provide
Helm charts that can be used directly or as sub-charts for other projects.

# 🔥 Quick start

__Dependencies__: [Helm](https://helm.sh/docs/intro/install/)

## Using the tobs CLI tool

The [CLI tool](/cli/README.md) provides the most seamless experience for interacting with tobs.

Getting started with the CLI tool is a two-step process: First you install the CLI tool locally, then you use the CLI tool to install the tobs stack into your Kubernetes cluster.

### Installing the CLI tool

To download and install tobs, run the following in your terminal, then follow the on-screen instructions.

```bash
curl --proto '=https' --tlsv1.2 -sSLf  https://tsdb.co/install-tobs-sh |sh
```

Alternatively, you can download the CLI directly via [our releases page](/releases)

### Using the tobs CLI tool to deploy the stack into your Kubernetes cluster

After setting up tobs run the following to install the tobs helm charts into your Kubernetes cluster

```bash
tobs install
```

This will deploy all of the tobs components into your cluster and provide instructions as to next steps.

### Getting started by viewing your metrics in Grafana
To see your Grafana dashboards after installation run

```bash
tobs grafana change-password <new_password>
tobs grafana port-forward
```
Then, point your browser to http://127.0.0.1:8080/ and login with the `admin` username.


More details about the CLI tool can be found [here](/cli/README.md)

# 🛠Alternative deployment methods

## Using the Helm charts without the CLI tool

Users sometimes want to use our Helm charts as sub-charts for other project or integrate them into their infrastructure without using our CLI tool. This is a supported use-case and instructions on using the Helm charts can be found [here](/chart/README.md).

# ✏️ Contributing

We welcome contributions to tobs, which is
licensed and released under the open-source Apache License, Version 2.  The
same [Contributor's
Agreement](//github.com/timescale/timescaledb/blob/master/CONTRIBUTING.md)
applies as in TimescaleDB; please sign the [Contributor License
Agreement](https://cla-assistant.io/timescale/tobs) (CLA) if
you're a new contributor.


[design-doc]: https://tsdb.co/prom-design-doc
[timescaledb-helm-cleanup]: https://github.com/timescale/timescaledb-kubernetes/blob/master/charts/timescaledb-single/admin-guide.md#optional-delete-the-s3-backups
[timescaledb-helm-repo]: https://github.com/timescale/timescaledb-kubernetes/tree/master/charts/timescaledb-single
[promscale-repo]: https://github.com/timescale/promscale
[promscale-helm]: https://github.com/timescale/promscale/tree/master/helm-chart
[prometheus-helm-hub]: https://prometheus-community.github.io/helm-charts
[prometheus-remote-tune]: https://prometheus.io/docs/practices/remote_write/
[grafana-helm-hub]: https://grafana.github.io/helm-charts
