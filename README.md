# ldap2ssh

A command line tool for generation SSH certificates using Vault and LDAP.

# Usage

```
# configure new account by adding it to ~/.ldap2ssh
ldap2ssh configure \
    --account myProject \
    --vault-address 'https://vault.example.com' \
    --vault-endpoint 'ssh-client-signer/sign/ca' \
    --default-key 'id_rsa.pub' \
    --ldap-user 'billy.bob'

# create the certificate and save it to ~/.ssh/{key_name}-cert.pub
ldap2ssh sign -a myProject
```

# Development

To start developing run `make mod` to download all dependencies.

To create a new release export the GitHub access token `export GITHUB_RELEASE_ACCESS_TOKEN="xxx"` and run:

```
make clean
make build
make dist
make release
```