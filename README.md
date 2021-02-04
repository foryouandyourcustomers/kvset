# kvget

Command line utility to set a single secret from an Azure keyvault.

## Installation

```bash
go get github.com/foryouandyourcustomers/kvset/cmd/kvset
```

or download the latest release.


## Usage

```bash
./kvset -h
Usage of ./kvget.linux:
  -s string
        Name of the secret to retrieve (env var: SECRET)
  -v string
        Name of the keyvault (env var: VAULT)
  -a string
        Value of the secret to set (env var: VALUE)
```

Lets set the value of the secret "myawesomesecret" to "mysupervalue" from the keyvault "fyayctestvault"
```bash
# via cli flags
./kvset -s myawesomesecret -v fyayctestvault -a mysupervalue

# via env vars
SECRET=myawesomesecret VAULT=fyayctestvault VALUE=mysupervalue ./kvset
```

The utility first tries to use the login from the azure cli.
If this fails it will try to retrieve credentials from the [runtime environment](https://docs.microsoft.com/en-us/azure/developer/go/azure-sdk-authorization#use-environment-based-authentication).
