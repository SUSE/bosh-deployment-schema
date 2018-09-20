# bosh-deployment-schema

## BOSH Deployment Manifest OAS v3 Spec

The [OpenAPI v3.0 spec](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.0.md#schemaObject) (OAS) describes the [BOSH deployment manifests](https://bosh.io/docs/manifest-v2/).

It can be used for validation, as well as documentation and code generation.

## Example Validation CLI

     go build
     ./bosh-deployment-schema openapi.yml manifest.yml

## Issues

* Some fields allow multiple types, like an integer or a range. The spec lists those fields as `string`. Thus using an integer in a deployment manifest will result in an validation error. YAML unmarshals unquoted numbers to integers.
  Example: using `instances: "1"` instead of `instances: 1` passes validation.
* For compatibility reasons some fields are allowed or might be missing.
  Example: `instance_groups: { azs: [] }`.
* `create-env` manifests, might not validate because of IaaS dependent sections, like `disk_pools`.
* OpenAPI v3.0 uses [JSON schema](http://json-schema.org/specification.html). JSON schema [does not allow any sibling nodes](https://github.com/OAI/OpenAPI-Specification/issues/1514) next to a `$ref`. The spec uses `allOf` to work around this limitation and add descriptions directly to the node. This might not work with all tooling.
