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
* tooling via make file

# Cross Compilation

Done via gox: github.com/mitchellh/gox

Needs packages:
* github.com/konsorten/go-windows-terminal-sequences

```
export VERSION=0.2
export NAME=ldap2ssh

gox -ldflags "-X main.Version=${VERSION}" \
    -osarch="darwin/amd64" \
    -osarch="linux/amd64" \
    -osarch="windows/amd64" \
    -output "build/{{.Dir}}_${VERSION}_{{.OS}}_{{.Arch}}/$NAME"\
    .

cd build/
for dir in $(ls .); do tar czvf "$dir.tar.gz" "$dir"; done
```