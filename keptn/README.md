# Webserver quality gates

## Create the project

keptn create project webserver --shipyard=./shipyard.yaml

## Onboard service to project

keptn onboard service webserver --project=webserver

## Add resource SLI

keptn add-resource --project=webserver --service=webserver --resource=prometheus/sli.yaml  --resourceUri=prometheus/sli.yaml --all-stages

## Add resource SLO

keptn add-resource --project=webserver --service=webserver --resource=slo.yaml --resourceUri=slo.yaml --all-stages

## Configure prometheus sli provider

keptn configure monitoring prometheus --project=webserver --service=webserver

## Start testing

### Run load tester

hey -z 1h -n -1 http://$(kubectl get svc -n orkestra webserver -ojsonpath='{.status.loadBalancer.ingress[0].ip}')/hello

### Trigger evaluation

keptn trigger evaluation --project=webserver --service=webserver --timeframe=5m