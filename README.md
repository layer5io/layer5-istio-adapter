<p style="text-align:center;" align="center">
  <img align="center" src="https://raw.githubusercontent.com/layer5io/layer5/master/src/assets/images/layer5/layer5-tagline/png/layer5-tag-community-white-bg.png" width="45%" /></p>

<div align="center">

![GitHub contributors](https://img.shields.io/github/contributors/layer5io/layer5.svg)
![GitHub](https://img.shields.io/github/license/layer5io/layer5.svg)
![GitHub issues by-label](https://img.shields.io/github/issues/layer5io/layer5/help%20wanted.svg?color=%23DDDD00)
[![Slack](https://img.shields.io/badge/Slack-@layer5.svg?logo=slack)](http://slack.layer5.io)
![Twitter Follow](https://img.shields.io/twitter/follow/layer5.svg?label=Follow&style=social)

</div>

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
## <a name="contributing"></a><a name="community"></a> Community and Contributing
Please do! Contributions, updates, [discrepancy reports](/../../issues) and [pull requests](/../../pulls) are welcome. This project is community-built and welcomes collaboration. Contributors are expected to adhere to the [Contributor Covenant](http://contributor-covenant.org) code of conduct.

<p>
This project is community-built and welcomes collaboration! See the <a href="https://docs.google.com/document/d/17OPtDE_rdnPQxmk2Kauhm3GwXF1R5dZ3Cj8qZLKdo5E/edit">Layer5 Community Welcome Guide</a>.
</p>
<p>
  <ul>
    <li style="list-style-type: circle;"><em><strong>Join</strong></em> <a href="https://drive.google.com/open?id=1c07UO9dS7_tFD-ClCWHIrEzRnzUJoFQ10EzfJTpS7FY">weekly community meeting</a> on <a href="https://calendar.google.com/calendar/b/1?cid=bGF5ZXI1LmlvX2VoMmFhOWRwZjFnNDBlbHZvYzc2MmpucGhzQGdyb3VwLmNhbGVuZGFyLmdvb2dsZS5jb20">Fridays from 10am - 11am Central</a>.</li>
    <ul>
        <li><em><strong>Watch</strong></em> community <a href="https://www.youtube.com/playlist?list=PL3A-A6hPO2IMPPqVjuzgqNU5xwnFFn3n0">meeting recordings</a>.</li>
    </ul>
    <li style="list-style-type: circle;"><em><strong>Access</strong></em> the <a href="https://drive.google.com/drive/u/4/folders/0ABH8aabN4WAKUk9PVA">community drive</a>.</li>
      <ul>
        <li><em><strong>Comment</strong></em> on the <a href="https://drive.google.com/open?id=15Gv1kTP8QHaMrDadfmgVmmS3hQn6YKFVag38Cr_JBEI">design document.</a></li>
      </ul>
    </ul>
</p>

*Not sure where to start?* Grab an open issue with the [help-wanted label](https://github.com/issues?utf8=✓&q=is%3Aopen+is%3Aissue+archived%3Afalse+org%3Alayer5io+label%3A%22help+wanted%22+).

### About

[Layer5.io](https://layer5.io) is a community-maintained repository of information pertaining to the technology ecosystem surrounding service meshes, API gateways, edge proxies, ingress and egress controllers - - microservice management in cloud native environments.

#### License

This repository and site are available as open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).
