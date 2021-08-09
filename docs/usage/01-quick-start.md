**[[Back To Overview]](../README.md)**

---

# Quick Start

## Connect `monoctl` to a m8 control plane

To connect your local `monoctl` you can use the `config init` command like so:

`monoctl config init -u api.<yourdomain>:443`

This will create a `monoconfig` file at the well known location `$HOME/.monoskope/config`.
The location can globally overridden for all `monoctl` command with the flag `--monoconfig` where the explicit file path can be specified.

## Authentication

Authentication is simple after you've initialized `monoctl`.
Every command which requires authentication automatically starts the authentication flow with the configured m8 control plane when you issue a command.

During the authentication flow `monoctl` will open a browser window and redirect you to the identity provider which has been [configured](../deployment/02-identity-provider-setup.md) with the m8 control plane.
Depending on how the identity provider is configured and if you're already authenticated in that browser with it, there might be a login prompt and a user consent request or not.

Users in Monoskope have to be explicitly created by an operator.
So if you have no user in the configured m8 control plane, you will not be able to log in.
