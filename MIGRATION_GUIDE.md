# 📚 Guia de Migrations

## 🚀 Setup Inicial

### 1. Instalar a ferramenta migrate CLI
```bash
make install-migrate
```

Ou manualmente:
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

**Importante:** Certifique-se que `~/go/bin` está no seu `PATH`:
```bash
export PATH=$PATH:~/go/bin
```

## 📝 Como usar

### Criar uma nova migration
```bash
make migration name=create_users_table
```

Isso criará 2 arquivos:
- `000001_create_users_table.up.sql` - Para aplicar a migration
- `000001_create_users_table.down.sql` - Para reverter a migration

### Executar migrations (UP)
```bash
make migrate-up
```

### Reverter migrations (DOWN)
```bash
make migrate-down
```

### Ver status das migrations
```bash
make migrate-status
```

## 📋 Exemplo de Migration

**000001_create_users_table.up.sql:**
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**000001_create_users_table.down.sql:**
```sql
DROP TABLE IF EXISTS users;
```

## 🔄 Workflow Completo

```bash
# 1. Subir o banco de dados
make docker-up

# 2. Criar uma migration
make migration name=create_users_table

# 3. Editar os arquivos .up.sql e .down.sql criados

# 4. Executar a migration
make migrate-up

# 5. Se precisar reverter
make migrate-down
```

## ⚠️ Boas Práticas

1. **Sempre** crie o arquivo `.down.sql` para poder reverter
2. **Nunca** modifique uma migration já aplicada em produção
3. **Teste** as migrations localmente antes de aplicar em produção
4. **Use** transações quando possível (`BEGIN; ... COMMIT;`)
5. **Nomeie** as migrations de forma descritiva

## 🐛 Troubleshooting

### Erro: "Dirty database version"
Isso acontece quando uma migration falha no meio. Para resolver:
```bash
# Conecte ao banco e force a versão correta
docker exec -it app-db psql -U pay -d app
# No psql:
UPDATE schema_migrations SET dirty = false WHERE version = X;
```

### Migration já foi aplicada
Se você modificou uma migration que já foi aplicada:
1. Reverta: `make migrate-down`
2. Modifique os arquivos
3. Aplique novamente: `make migrate-up`
