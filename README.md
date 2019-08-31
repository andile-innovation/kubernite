# kubernite
![kubernite](https://raw.githubusercontent.com/andile-innovation/kubernite/master/kubernite.png)

[![Build Status](https://cloud.drone.io/api/badges/andile-innovation/kubernite/status.svg)](https://cloud.drone.io/andile-innovation/kubernite)

Kubernite is a Drone Plugin for Kubernetes built using official kubernetes client-go

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