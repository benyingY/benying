## Services

Templates:
- python-service-template: FastAPI service with /healthz and pytest test
- go-service-template: net/http service with /healthz and go test

Use the generator script:
- scripts/new_service.sh <python|go> <service-name>

Config convention:
- Set APP_ENV=dev|staging to load config/<env>.env
- Environment variables override config file values
- Secrets live in env vars; see config/.env.example for placeholders
