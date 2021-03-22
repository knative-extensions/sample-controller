# Schema Tool

Schema is a seed of a CLI tool that a downstream can use to use reflection and go file inspection to 
generate a base version of OpenAPI for a CRD. The resulting schema will be used by Kubernetes to
provide results from `kubectl explain <type>` calls, and type validation.  

## Integration steps

### Demo

Register the type you wish to expose in the CLI,

```go
registry.Register(&example.AddressableService{})
```

Run the `schema dump` command for the <Kind> you wish to have schema for,

```
cd ./schema
go run ./ dump AddressableService | pbcopy
```

Paste this inside the CRD for LoremIpsum,

```yaml
...
      schema:
        openAPIV3Schema:
          <**paste**>
      additionalPrinterColumns:
...
```

### Downstream

Start with [example.go](https://github.com/knative/hack/blob/main/schema/example.go), copy this into the downstream and modify which 
kinds are registered via `registry.Register`. You can register more than one kind at a time. (TODO: support versions in the CLI.)