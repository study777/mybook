example hosts

```
[OSEv3:children]
masters
nodes
etcd

[OSEv3:vars]
ansible_ssh_user=root
deployment_type=openshift-enterprise
openshift_release=v3.5

# Enable cockpit
osm_use_cockpit=true
#
# Set cockpit plugins
osm_cockpit_plugins=['cockpit-kubernetes']

oreg_url=registry.example.com:5000/openshift3/ose-${component}:${version}
openshift_docker_additional_registries=registry.example.com:5000 ,172.30.0.1/16:5000
openshift_docker_insecure_registries=registry.example.com:5000 ,172.30.0.1/16:5000
openshift_examples_modify_imagestreams=true

openshift_master_identity_providers=[{'name': 'htpasswd_auth', 'login': 'true', 'challenge':'true', 'kind': 'HTPasswdPasswordIdentityProvider', 'filename': '/etc/origin/master/htpasswd'}]

openshift_hosted_registry_selector="infra=true"

openshift_hosted_router_selector="router=true"

openshift_master_cluster_method=native
openshift_master_cluster_hostname=master.example.com
openshift_master_cluster_public_hostname=master.example.com

# default subdomain to use for exposed routes
openshift_master_default_subdomain=apps.example.com

openshift_node_kubelet_args={'pods-per-core': ['10'], 'max-pods':['60']}

openshift_metrics_install_metrics=true
openshift_hosted_metrics_deploy=true
openshift_hosted_metrics_public_url=https://hawkular-metrics.apps.example.com/hawkular/metrics
openshift_metrics_image_prefix=registry.example.com:5000/openshift3/
openshift_metrics_image_version=3.5.0

openshift_hosted_logging_deploy=true
openshift_logging_image_prefix=registry.example.com:5000/openshift3/
openshift_logging_image_version=3.5.0

# enable ntp on masters to ensure proper failover
openshift_clock_enabled=false

#Skip checking docker images & memory 
#openshift_disable_check=docker_image_availability,memory_availability

# host group for masters
[masters]
master.example.com

# host group for etcd
[etcd]
master.example.com



# host group for nodes, includes region info
[nodes]
master.example.com
node1.example.com openshift_node_labels="{'region': 'primary', 'zone': 'default'}"
node2.example.com openshift_node_labels="{'region': 'primary', 'zone': 'default'}"
infra.example.com openshift_node_labels="{'region': 'infra', 'zone': 'default','infra': 'true'}"
route.example.com openshift_node_labels="{'region': 'infra', 'zone': 'default','router': 'true'}"
```