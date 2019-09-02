# Kubernite
![kubernite](https://github.com/andile-innovation/kubernite/blob/master/images/kubernite.png?raw=true)

[![Build Status](https://cloud.drone.io/api/badges/andile-innovation/kubernite/status.svg)](https://cloud.drone.io/andile-innovation/kubernite)

Kubernite is a [Drone](https://drone.io/) Plugin for Kubernetes written in golang using the official [client-go](https://github.com/kubernetes/client-go) kubernetes api client library.
## Motivation
## Working Principle
![working principle](https://github.com/andile-innovation/kubernite/blob/master/images/work_flow.png?raw=true)
## Example Usage

```yaml
  - name: deploy
    image: tbcloud/kubernite
    settings:
        kubernetes_server:
          from_secret: kubernetes_server
        kubernetes_cert_data:
          from_secret: kubernetes_cert_data
        kubernetes_client_cert_data:
          from_secret: kubernetes_client_cert_data
        kubernetes_client_key_data:
          from_secret: kubernetes_client_key_data
        deployment_tag_repository_path: /drone/src
        deployment_file_path: src/deployments/kubernetes/Deployment.yaml
```
## Credits
- [drone-kubernetes](https://github.com/honestbee/drone-kubernetes) by the [honestbee](https://github.com/honestbee)