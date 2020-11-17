# get-credhub-secrets action

This action fetches secrets from a credhub instance.

## Inputs

### `api`

**Required** Set the CredHub API target where commands are sent.

### `get`

**Required** Newline-separated list of secrets to fetch. Secrets must be of the format 
path/to/secret or path/to/secret.key

### `username`

**Required** Authentication username.

### `password`

**Required** Authentication password.

### `ca`

**Required** Trusted CA certificate (x509).

### `insecureSkipTLSValidation`

**Required** Disable TLS validation (not recommended). Must be a string.
**Default** "false"

## Outputs

Each secret is prefixed with an output name. The secret's resolved access value
will be available at that output in future build steps.

For example:

```yaml
steps:
- id: secrets
  uses: ciriarte/github-actions/get-credhub-secrets@main
  with:
    get: |-
      token:path/to/database-token
```

will be available in future steps as the output "token":

```yaml
- id: publish
  uses: foo/bar@main
  env:
    TOKEN: ${{ steps.secrets.outputs.token }}
```

## Example usage

```yaml
on: [push]

jobs:
  get_credhub_secrets:
    runs-on: ubuntu-latest
    name: get-credhub-secrets [Unit Test]
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Fetching Secrets
        uses: ./get-credhub-secrets
        id: credhub
        with:
          api: "https://credhub.example.com:8844"
          username: "some_user"
          password: "some_password"
          ca: |-
            -----BEGIN CERTIFICATE-----
            ... (elided)
            -----END CERTIFICATE-----
          get: |-
            HITCHHIKER:/concourse/main/the_ultimate_question_of_life_the_universe_and_everything
            BIRTH_CERTIFICATE: /concourse/main/birth.certificate
      - name: Listing the secrets
        run: |-
          echo "The answer ${{ steps.credhub.outputs.HITCHHIKER }}"
          echo "A CA cert ${{ steps.credhub.outputs.BIRTH_CERTIFICATE }}"
```

## License

See [LICENSE](LICENSE)
