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
ldap2ssh sign -a wlw
```

# Development

Improvements:
* path joining
* create empty config if it not exists
* tooling via make file
* remove JumpCloud; add general `credentials`
* parameter for `sign` to save cert to a specific file
* parameter for `sign` to force Vault token renewal