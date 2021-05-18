## Usage
```bash
git clone https://github.com/cqbqdd11519/replace-description.git
cd replace-description

go run go run cmd/replace-description/main.go \
--input test/cicd.tmax.io_integrationconfigs.yaml \
--output test/output.yaml \
--schemaFile test/translation.xlsx \
--schemaSheet IntegrationConfig
```
