.PHONY: build run dev migration migrate-up migrate-down docker-up docker-down clean help

help: ## Mostra todos os comandos disponíveis
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila a aplicação
	@go build -o bin/app cmd/main.go

run: ## Roda a aplicação
	@go run cmd/main.go

migration: ## Cria uma nova migration (uso: make migration name=create_users_table)
	@if [ -z "$(name)" ]; then \
		echo "❌ Erro: Especifique o nome da migration"; \
		echo "Uso: make migration name=create_users_table"; \
		exit 1; \
	fi
	@mkdir -p cmd/migrate/migrations
	@echo "📝 Criando migration: $(name)"
	@if ! command -v migrate &> /dev/null; then \
		echo "⚠️  migrate CLI não encontrado. Instalando..."; \
		go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		echo "✅ migrate CLI instalado!"; \
		export PATH="$$PATH:$$HOME/go/bin"; \
	fi
	@$$HOME/go/bin/migrate create -ext sql -dir cmd/migrate/migrations -seq $(name) || migrate create -ext sql -dir cmd/migrate/migrations -seq $(name)
	@echo "✅ Migration criada em cmd/migrate/migrations/"

migrate-up: ## Executa todas as migrations pendentes
	@go run cmd/migrate/main.go up

migrate-down: ## Reverte a última migration
	@go run cmd/migrate/main.go down

migrate-force: ## Força uma versão específica (uso: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "❌ Erro: Especifique a versão"; \
		echo "Uso: make migrate-force version=1"; \
		exit 1; \
	fi
	@go run cmd/migrate/main.go force $(version)

migrate-status: ## Mostra o status das migrations
	@go run cmd/migrate/main.go status || true

install-migrate: ## Instala a ferramenta migrate CLI
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ Migrate instalado! Certifique-se que ~/go/bin está no PATH"

clean: ## Remove arquivos compilados
	@rm -rf bin/
	@echo "✅ Limpo!"

deps: ## Baixa e atualiza dependências
	@go mod download
	@go mod tidy

test: ## Roda os testes
	@go test -v ./...
