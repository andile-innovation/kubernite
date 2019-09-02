# Kubernite
![kubernite](https://github.com/andile-innovation/kubernite/blob/master/images/kubernite.png?raw=true)

[![Build Status](https://cloud.drone.io/api/badges/andile-innovation/kubernite/status.svg)](https://cloud.drone.io/andile-innovation/kubernite)

Kubernite is a [Drone](https://drone.io/) Plugin for Kubernetes written in golang using the official [client-go](https://github.com/kubernetes/client-go) kubernetes api client library.
## Motivation
Kubernite was built out of a desire to achieve complete automation of the deployment stage of the development cycle of an application running in a kubernetes cluster.

This plugin was developed in an environment where 2 situations needed to be catered for:
1. When pushing version tags to an application's branch:
   1. A docker image tagged with the version tag is to be built and pushed to an image repository
   2. An existing kubernetes deployment is to be redeployed with the new image tag triggering a rolling update
2. When pushing commits to an application's branch:
   1. A docker image tagged 'latest' is to be built and pushed to an image repository
   2. An existing kubernetes deployment is to be redeployed with the image tag 'latest' triggering a rolling update

Various established drone plugins already exist to cater for (i.) in each of these situations (see [drone-docker](https://github.com/drone-plugins/drone-docker)). As such this plugin was written to deal strictly with (ii.). Kubernite implements the following functionality:
1. Trigger the redeployment of an *existing* deployment in a kubernetes cluster with consistent and traceable [kubernetes.io/change-cause](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#checking-rollout-history-of-a-deployment) annotations
2. Update the deployment manifest file and commit the changes to source control to keep the state of the deployment in the cluster in sync with the manifest files which describe it. (This is to keep with [Declarative Management of Kubernetes Objects](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/))
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
See [additional examples](#additional-examples).
### Plugin Settings 
|Setting|Description|
|---|---|
|kubernetes_server|URL of kubernetes server. Can be found in the kube config at key 'cluster.server'. The kube config can typcially be found at **$USER/.kube/config** or by running **kubectl config view**.|
|kubernetes_cert_data|Root certificate used to verify the certificate presented by the API server when [transport security](https://kubernetes.io/docs/reference/access-authn-authz/controlling-access/#transport-security) is being established. Can be found in the kube config at key 'cluster.certificate-authority-data'. The kube config can typcially be found at **$USER/.kube/config**.|
|kubernetes_client_cert_data|Public client certificate data for client X509 certificate. Used in authentication process. Can be found in kube config at key 'user.client-certificate-data'. See [authenticating with X509 Client Certs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs), and  [generating certificates](https://kubernetes.io/docs/concepts/cluster-administration/certificates/). Can be found in the kube config at key 'user.client-key-data'. The kube config can typcially be found at **$USER/.kube/config**. Note that if you are using a hosted service such as [Digital Ocean KaaS](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/#download-the-configuration-file) you may need to download your config file from them to get access to this data.|
|kubernetes_client_key_data|Private key data for client X509 certificate. Used in authentication process. Can be found in kube config at key 'user.client-key-data'. See [authenticating with X509 Client Certs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs), and  [generating certificates](https://kubernetes.io/docs/concepts/cluster-administration/certificates/). Can be found in the kube config at key 'user.client-key-data'. The kube config can typcially be found at **$USER/.kube/config**. Note that if you are using a hosted service such as [Digital Ocean KaaS](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/#download-the-configuration-file) you may need to download your config file from them to get access to this data.|
|deployment_file_path|Path to [deployment manifest](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#writing-a-deployment-spec) .yaml or .yml file which describes the deployment to be redeployed by kubernite.|
|deployment_tag_repository_path|[**optional** - default is **/drone/src**] Path to root of repository from which tag/commit information is drawn to update the kubernetes.io/change-cause annotations in the deployment file. Defaults to default drone working directory (i.e. /drone/src) which is typically the root of the repository which has triggered the deployment.|
|deployment_image_name|[**optional** if pod template contains only 1 image][**required** if pod template contains more than 1 image] The name of the image whose tag should be updated.|
|dry_run|[**optional** - default is **false**] If set, no deployment takes place and the updated deployment file which would be applied to the cluster is printed out in json format.|
|deployment_file_repository_path|[**optional** - no default] Path to root of repository to which deployment file with updated kubernetes.io/change-cause annotations will be committed and pushed if settings.commit_deployment is set.|
|commit_deployment|[**optional** - default is **false**] If set, deployment file with updated kubernetes.io/change-cause annotations will be committed and pushed to repository with it's root at settings.deployment_file_repository_path.|
## Working Principle
A redeployment of an existing deployment is triggered when the pod template part of the deployment's .spec section is changed and the associated resource is updated.
Kubernite leverages this behaviour to trigger a redeployment each time it is run by updating annotations in the metadata of the template and/or an image tag.

This behaviour and the logic around it is illustrated in the following diagram. **Please Note the warning related to the commit_deployment setting.**

![working principle](https://github.com/andile-innovation/kubernite/blob/master/images/work_flow.png?raw=true)
## Additional Examples
### 
## FAQ
## Credits
- [drone-kubernetes](https://github.com/honestbee/drone-kubernetes) by the [honestbee](https://github.com/honestbee)
## TODO
- add support for other client authentication methods