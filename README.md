# Kubernite
![kubernite](https://github.com/andile-innovation/kubernite/blob/master/images/kubernite.png?raw=true)

[![Build Status](https://cloud.drone.io/api/badges/andile-innovation/kubernite/status.svg)](https://cloud.drone.io/andile-innovation/kubernite)

Kubernite is a [Drone](https://drone.io/) Plugin for Kubernetes written in golang using the official [client-go](https://github.com/kubernetes/client-go) kubernetes api client library.

##Warning
![warning](https://github.com/andile-innovation/kubernite/blob/master/images/warningSign.png?raw=true)

The updated deployment manifest file is not yet pushed to source control. The following extra drone step needs to be added 
to get around this. 

```yaml
  - name: push infrastucture
    image: alpine
    volumes:
      - name: infrastructure_volume
        path: /projects/infrastructure
    commands:
      - apk add --no-cache git
      - cd /projects/infrastructure
      - git push
    when:
      event:
        exclude:
          - pull_request
```

## Motivation
Kubernite was built out of a desire to achieve complete automation of the deployment stage of the development cycle of an application running in a kubernetes cluster.

This plugin was developed in an environment where 2 situations needed to be catered for:
1. When pushing version tags to an application's branch:
   1. A docker image tagged with the version tag is to be built and pushed to an image repository
   2. An existing kubernetes deployment is to be redeployed with the new image tag triggering a rolling update
2. When pushing commits to an application's branch:
   1. A docker image tagged 'latest' is to be built and pushed to an image repository
   2. An existing kubernetes deployment is to be redeployed with the image tag 'latest' triggering a rolling update

Various established drone plugins already exist to cater for [i.] in each of these situations (see [drone-docker](https://github.com/drone-plugins/drone-docker)). As such this plugin was written to deal strictly with [ii.]. Kubernite implements the following functionality:
1. Trigger the redeployment of an *existing* deployment in a kubernetes cluster with consistent and traceable [kubernetes.io/change-cause](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#checking-rollout-history-of-a-deployment) annotations
2. Update the deployment manifest file and commit the changes to source control to keep the state of the deployment in the cluster in sync with the manifest files which describe it. (This is to keep with [Declarative Management of Kubernetes Objects](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/))
   - see [object management techniques](https://kubernetes.io/docs/concepts/overview/working-with-objects/object-management/) 
## Usage
### Elementary Example
```yaml
  - name: deploy
    image: tbcloud/kubernite:<version>
    settings:
        kubernetes_server:
          from_secret: kubernetes_server
        kubernetes_cert_data:
          from_secret: kubernetes_cert_data
        kubernetes_client_cert_data:
          from_secret: kubernetes_client_cert_data
        kubernetes_client_key_data:
          from_secret: kubernetes_client_key_data
        deployment_file_path: src/deployments/kubernetes/Deployment.yaml
```
See [additional examples](#additional-examples).
See [plugin settings](#plugin-settings).
See [drone secrets](https://readme.drone.io/configure/secrets/).
### Plugin Settings 
|Setting|Description|
|---|---|
|kubernetes_server|URL of kubernetes server. Can be found in the kube config at key 'cluster.server'. The kube config can typcially be found at **$USER/.kube/config** or by running **kubectl config view**.|
|kubernetes_cert_data|Root certificate used to verify the certificate presented by the API server when [transport security](https://kubernetes.io/docs/reference/access-authn-authz/controlling-access/#transport-security) is being established. Can be found in the kube config at key 'cluster.certificate-authority-data'. The kube config can typcially be found at **$USER/.kube/config**.|
|kubernetes_client_cert_data|Public client certificate data for client X509 certificate. Used in authentication process. Can be found in kube config at key 'user.client-certificate-data'. See [authenticating with X509 Client Certs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs), and  [generating certificates](https://kubernetes.io/docs/concepts/cluster-administration/certificates/). Can be found in the kube config at key 'user.client-key-data'. The kube config can typcially be found at **$USER/.kube/config**. Note that if you are using a hosted service such as [Digital Ocean KaaS](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/#download-the-configuration-file) you may need to download your config file from them to get access to this data.|
|kubernetes_client_key_data|Private key data for client X509 certificate. Used in authentication process. Can be found in kube config at key 'user.client-key-data'. See [authenticating with X509 Client Certs](https://kubernetes.io/docs/reference/access-authn-authz/authentication/#x509-client-certs), and  [generating certificates](https://kubernetes.io/docs/concepts/cluster-administration/certificates/). Can be found in the kube config at key 'user.client-key-data'. The kube config can typcially be found at **$USER/.kube/config**. Note that if you are using a hosted service such as [Digital Ocean KaaS](https://www.digitalocean.com/docs/kubernetes/how-to/connect-to-cluster/#download-the-configuration-file) you may need to download your config file from them to get access to this data.|
|deployment_file_path|Path to [deployment manifest](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#writing-a-deployment-spec) .yaml or .yml file which describes the deployment to be redeployed by kubernite.|
|deployment_tag_repository_path|[**optional** - default is **/drone/src**] Path to root of repository from which tag/commit information is drawn to update the kubernetes.io/change-cause annotations in the deployment file. Defaults to default drone working directory (i.e. /drone/src) which is typically the root of the repository which has triggered the deployment.|
|deployment_image_name|[**optional** if pod template contains only 1 image, **required** if pod template contains more than 1 image] The name of the image whose tag should be updated.|
|dry_run|[**optional** - default is **false**] If set, no deployment takes place and the updated deployment file which would be applied to the cluster is printed out in json format.|
|deployment_file_repository_path|[**optional** only if commit_deployment is set to **false** - no default] Path to root of repository to which deployment file with updated kubernetes.io/change-cause annotations will be committed and pushed if settings.commit_deployment is set.|
|commit_deployment|[**optional** - default is **false**] If set, deployment file with updated kubernetes.io/change-cause annotations will be committed and pushed to repository with it's root at settings.deployment_file_repository_path.|
## Working Principle
A redeployment of an existing deployment is triggered when the pod template part of the deployment's .spec section is changed and the associated resource is updated.
Kubernite leverages this behaviour to trigger a redeployment each time it is run by updating annotations in the metadata of the template and/or an image tag.

This behaviour and the logic around it is illustrated in the following diagram. **Please Note the warning related to the commit_deployment setting.**

![working principle](https://github.com/andile-innovation/kubernite/blob/master/images/work_flow.png?raw=true)
## Additional Examples
### Handle 'Tag' and 'Other' events differently
This example demonstrates the following:
1. specifying a workspace structure which changes the location into which drone clones the repository that triggers the pipeline
2. declaring a volume to share files between build stages
3. building a node project (this particular example considers a react project)
4. using [drone-docker](https://github.com/drone-plugins/drone-docker) to build and push an image to dockerhub
5. using docker:git to clone a kubernetes infrastructure repository
6. for tag events: redeploy and commit and push updated deployment file to infrastructure repository
7. for other events: redeploy 
```yaml
kind: pipeline
name: default

# [1.] specify work space structure
# results in the following:
# ├── repositories
#       └── foo <-- drone clones repository here, this is default working directory
workspace:
  base: /repositories
  path: foo

# [2.] volume to share files between steps. this is necessary as by default only files in the 
# base of the workspace are maintained between steps (i.e. /repositories/foo).
volumes:
  - name: infrastructure_volume
    temp: {}

steps:
  # [3.] build foo project
  - name: build foo
    image: node
    commands:
      - yarn install
      - yarn build
  
  # [4.] build and deploy foo image
  - name: build & deploy image
    image: plugins/docker
    settings:
      repo: fooOwner/foo
      auto_tag: true
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password
      dockerfile: /projects/foo/Dockerfile
      context: /projects/foo/build

  # [5.] clone infrastructure repository
  # results in the following:
  # ├── projects
  #       └── foo <-- drone clones repository here, this is default working directory
  #       └── infrastructure <-- this stage creates this from cloning
  - name: clone infrastucture
    image: docker:git
    volumes:
      - name: infrastructure_volume
        path: /projects/infrastructure
    commands:
      - cd /projects
      - git clone https://github.com/fooOwner/infrastructure.git

  # [6.] this stage will run only with tag events 
  - name: deploy on tag
    image: tbcloud/kubernite:<version>
    volumes:
      - name: infrastructure_volume
        path: /projects/infrastructure
    settings:
        kubernetes_server:
          from_secret: kubernetes_server
        kubernetes_cert_data:
          from_secret: kubernetes_cert_data
        kubernetes_client_cert_data:
          from_secret: kubernetes_client_cert_data
        kubernetes_client_key_data:
          from_secret: kubernetes_client_key_data
        deployment_file_path: /projects/infrastructure/Deployment.yaml
        commit_deployment: true
        deployment_file_repository_path: /projects/infrastructure
    when:
      event:
        - tag

  # [7.] this stage will run only with tag events
  - name: deploy on other
    image: tbcloud/kubernite:<version>
    volumes:
      - name: infrastructure_volume
        path: /projects/infrastructure
    settings:
        kubernetes_server:
          from_secret: kubernetes_server
        kubernetes_cert_data:
          from_secret: kubernetes_cert_data
        kubernetes_client_cert_data:
          from_secret: kubernetes_client_cert_data
        kubernetes_client_key_data:
          from_secret: kubernetes_client_key_data
        deployment_file_path: /projects/infrastructure/Deployment.yaml
    when:
      event:
        exclude:
          - tag
```
See [drone triggers](https://docker-runner.docs.drone.io/configuration/trigger/)
## FAQ
- Why/what kind of tags are used?
## Credits
- [drone-kubernetes](https://github.com/honestbee/drone-kubernetes) by the [honestbee](https://github.com/honestbee)
## TODO
- add support for other client authentication methods