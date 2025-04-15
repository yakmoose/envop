# Envop
Envop is a simple cli tool for storing .env files in 1Password. 
It reads files in a similar order to the symfony/dotenv package, 
but does not mess around with picking up an environment from an `.env.local` file. 

You need a 1Password service account to use this tool eg:
```bash
OP_SERVICE_ACCOUNT_TOKEN="$(op service-account create tf --vault DeploymentSecrets:read_items,write_items --raw --expires-in 4h)" \
envop <options>
```

## Urls of interest
- https://1password.com
- https://developer.1password.com/docs/connect
- https://github.com/1Password/connect
- https://github.com/1password/onepassword-sdk-go
