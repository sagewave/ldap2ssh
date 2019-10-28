package main

// SignFlags flags used for the `sign` command
type SignFlags struct {
	Account string
	Key     string
	Token   string
	Force   bool
	Outfile string
}

// ConfigureFlags flags used for the `configure` command
type ConfigureFlags struct {
	Account       string
	User          string
	VaultAddress  string
	VaultEndpoint string
	DefaultKey    string
}
