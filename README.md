# Kubernite
![kubernite](https://github.com/andile-innovation/kubernite/blob/master/images/kubernite.png?raw=true)

[![Build Status](https://cloud.drone.io/api/badges/andile-innovation/kubernite/status.svg)](https://cloud.drone.io/andile-innovation/kubernite)

Kubernite is a [Drone](https://drone.io/) Plugin for Kubernetes written in golang using the official [client-go](https://github.com/kubernetes/client-go) kubernetes api client library.
## Motivation
Kubernite was built out of a desire to achieve complete automation of the deployment stage of the development cycle of an application running in a kubernetes cluster.
To that end the plugin was developed with the following functionality:
1. trigger the redeployment of an *existing* deployment in a kubernetes cluster with consistent and traceable [kubernetes.io/change-cause](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#checking-rollout-history-of-a-deployment) annotations
2. update the deployment manifest file and commit the changes to source control to keep the state of the deployment in the cluster in sync with the manifest files which describe it. (This is to keep with [Declarative Management of Kubernetes Objects](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/))
   - see [object management techniques](https://kubernetes.io/docs/concepts/overview/working-with-objects/object-management/) 
## Usage
### Elementary Example
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
## Working Principle
![working principle](https://github.com/andile-innovation/kubernite/blob/master/images/work_flow.png?raw=true)
## FAQ
## Credits
- [drone-kubernetes](https://github.com/honestbee/drone-kubernetes) by the [honestbee](https://github.com/honestbee)