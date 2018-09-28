#!/usr/bin/env bash

set -ex

oc apply -f hack/openshift/template.yaml
oc process gloo-ee-installation-template \
 -p APISERVER_OPENSHIFT_MASTER_IP=$(minishift ip) \
  | oc apply -f -

cat << EOF | oc apply -f -
kind: OAuthClient
apiVersion: oauth.openshift.io/v1
metadata:
 name: gloo
secret: gloo
redirectURIs:
 - "http://localhost:8080"
grantMethod: prompt
EOF

kubectl port-forward deployment/apiserver-ui 8080