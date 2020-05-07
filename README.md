# Webserver app for progressive delivery demo

A sample webserver that serves a simple `/hello` endpoint with support for injecting artificial errors and delays, and exporting metrics to prometheus. Sample grafana dashboards also inclueded.

This demo application is used to demonstrate Continuous Delivery using [ArgoCD](https://argoproj.github.io/argo-cd/) coupled with Progressive Delivery using Canary testing and Automatic Canary Analysis (ACA) powered by (Flagger)[https://flagger.app/]

## Navigation

- **deploy/** - Deployment artifacts
- **deploy/helm** - Helm based deployment artifacts
- **deploy/kustomize**- Kustomize baed deployment artifacts
- **webserver_metrics.json** - Prometheus metrics manifest for stats exported by webserver. We use [`prometheus_json2go`](https://github.com/nitishm/prometheus-json2go) as a code generator for Golang metric object generation.
- **Remaining** - Contain source code for the application and Dockerfile for webserver container image generation.

*...more content to come...*
