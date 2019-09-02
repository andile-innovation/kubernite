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
    image: tbcloud/kubernite:<version>
    settings:
        # url for accessing the kubernetes server
        kubernetes_server:
          from_secret: kubernetes_server
        # cluster ca certificate (for TLS)
        kubernetes_cert_data:
          from_secret: kubernetes_cert_data
        # client ca certificate (for TLS)
        kubernetes_client_cert_data:
          from_secret: kubernetes_client_cert_data
        # for cluster access
        kubernetes_client_key_data:
          from_secret: kubernetes_client_key_data
        # path to deployment manifest file
        deployment_file_path: src/deployments/kubernetes/Deployment.yaml
```
### Plugin Settings 
|Setting|Description|
|---|---|
|kubernetes_server|URL of kubernetes server. Can be found in the kube config at key 'cluster.server'. The kube config can typcially be found at **$USER/.kube/config** or by running **kubectl config view**.|
|kubernetes_cert_data|Root certificate used to verify the certificate presented by the API server when [transport security](https://kubernetes.io/docs/reference/access-authn-authz/controlling-access/#transport-security) is being established. Can be found in the kube config at key 'cluster.certificate-authority-data'. The kube config can typcially be found at **$USER/.kube/config**.|
|kubernetes_client_cert_data|Public client certificate data for client X509 certificate. Used in authentication process. Can be found in kube config at key 'user.client-certificate-data'. See [authenticating with X509 Client Certs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs), and  [generating certificates](https://kubernetes.io/docs/concepts/cluster-administration/certificates/). Can be found in the kube config at key 'user.client-key-data'. The kube config can typcially be found at **$USER/.kube/config**. Note that if you are using a hosted service such as [Digital Ocean KaaS](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/#download-the-configuration-file) you may need to download your config file from them to get access to this data.|
|kubernetes_client_key_data|Private key data for client X509 certificate. Used in authentication process. Can be found in kube config at key 'user.client-key-data'. See [authenticating with X509 Client Certs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs), and  [generating certificates](https://kubernetes.io/docs/concepts/cluster-administration/certificates/). Can be found in the kube config at key 'user.client-key-data'. The kube config can typcially be found at **$USER/.kube/config**. Note that if you are using a hosted service such as [Digital Ocean KaaS](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/#download-the-configuration-file) you may need to download your config file from them to get access to this data.|
|deployment_file_path|Path to [deployment manifest](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#writing-a-deployment-spec) .yaml or .yml file which describes the deployment to be redeployed by kubernite.|
|deployment_tag_repository_path|[**optional**]|
|dry_run|[**optional**]|
|commit_deployment|[**optional**]|
|deployment_file_repository_path|[**optional**]|
## Working Principle
![working principle](https://github.com/andile-innovation/kubernite/blob/master/images/work_flow.png?raw=true)
## FAQ
## Credits
- [drone-kubernetes](https://github.com/honestbee/drone-kubernetes) by the [honestbee](https://github.com/honestbee)
## TODO
- add support for other client authentication methods