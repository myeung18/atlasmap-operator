#!/bin/sh

if [[ -z ${1} ]]; then
    CATALOG_NS="operator-lifecycle-manager"
else
    CATALOG_NS=${1}
fi

CSV=`cat deploy/olm-catalog/atlasmap-operator/0.0.1/atlasmap-operator.v0.0.1.clusterserviceversion.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
CRD=`cat deploy/crds/atlasmap_v1alpha1_atlasmap_crd.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`
PKG=`cat deploy/olm-catalog/atlasmap-operator/0.0.1/atlasmap.package.yaml | sed -e 's/^/          /' | sed '0,/ /{s/          /        - /}'`

cat << EOF > deploy/olm-catalog/atlasmap-operator/0.0.1/catalog-source.yaml
apiVersion: v1
kind: List
items:
  - apiVersion: v1
    kind: ConfigMap
    metadata:
      name: atlasmap-resources
      namespace: ${CATALOG_NS}
    data:
      clusterServiceVersions: |
${CSV}
      customResourceDefinitions: |
${CRD}
      packages: >
${PKG}

  - apiVersion: operators.coreos.com/v1alpha1
    kind: CatalogSource
    metadata:
      name: atlasmap-resources
      namespace: ${CATALOG_NS}
    spec:
      configMap: atlasmap-resources
      displayName: AtlasMap Operators
      publisher: Red Hat
      sourceType: internal
    status:
      configMapReference:
        name: atlasmap-resources
        namespace: ${CATALOG_NS}
EOF
