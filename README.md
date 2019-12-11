# Layer5 Istio Adapter

## Deploying the adapter to a Kubernetes cluster with Istio
Deploy the adapter itself and its service:
```
kubectl -n istio-system apply k8s.yaml
```
Now navigate into the `testdata` folder:
```
cd testdata
```
Apply the attributes and template yaml:
```
kubectl -n istio-system apply -f attributes.yaml -f template.yaml
```
Now lets apply the layer5 adater definition:
```
kubectl -n istio-system apply -f layer5.yaml
```
Finally, lets deploy the operator config:
```
kubectl -n istio-system apply -f sample_operator_cfg_cluster.yaml
```

## Setup dev environment
```
mkdir -p $GOPATH/src/istio.io/ && \
cd $GOPATH/src/istio.io/  && \
git clone https://github.com/istio/istio
```

```
export MIXER_REPO=$GOPATH/src/istio.io/istio/mixer

export ISTIO=$GOPATH/src/istio.io
```

```
cd $MIXER_REPO/adapter
git clone https://github.com/layer5io/layer5-istio-adapter layer5
cd layer5
```

To build:
```
make build
```

To run locally:
```
make run
```

To run mixer server locally:
```
make mixs
```

To run mixer client:
```
make mixc
```

If you make changes to the proto file, then we have to generate the associated Go code:
```
make gen
```