# 📋 Sistema de Microsserviços — Order, Payment & Shipping

Este projeto implementa uma arquitetura de microsserviços distribuídos utilizando a **Arquitetura Hexagonal** (Ports and Adapters) em Go. São **três serviços** independentes que se comunicam via **gRPC** e utilizam bancos de dados **MySQL** isolados via Docker.

---

## 🏗️ Fluxo de Comunicação

<img width="1693" height="929" alt="image" src="https://github.com/user-attachments/assets/4a31ff0c-c312-475d-9056-0cf03dadb52d" />

1. O **cliente** envia uma requisição gRPC de criação de pedido para o **Order Service** (`porta 3000`).
2. O **Order Service** valida o pedido (máx. 50 itens) e verifica se cada `product_code` existe na tabela de estoque.
3. O **Order Service** persiste o pedido no banco `order` e chama o **Payment Service** (`porta 3001`) via gRPC.
4. O **Payment Service** valida (valor máx. R$ 1000), persiste no banco `payment` e retorna sucesso.
5. Com o pagamento aprovado, o **Order Service** chama o **Shipping Service** (`porta 3002`) via gRPC para calcular o prazo de entrega.
6. O **Shipping Service** calcula os dias de entrega com base na quantidade total de itens e persiste no banco `shipping`.
7. O **Order Service** atualiza o status do pedido para `Paid` e retorna o `order_id` ao cliente.

> **Regra de frete:** prazo mínimo de **1 dia** + **1 dia extra a cada 5 unidades** (ex: 10 itens = 3 dias, 20 itens = 5 dias).

> **Validações do Order:** máx. 50 itens no pedido | todos os `product_code` devem existir no estoque.

> **Validação do Payment:** valor total máx. R$ 1000.

---

## 🛠️ Pré-requisitos

- **Go** (versão 1.22 ou superior)
- **Docker**
- **grpcurl** — para disparar chamadas gRPC pelo terminal

```bash
# Instalar grpcurl (Linux)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

---

## � Executar com Docker Compose (recomendado)

Para subir todos os serviços de uma vez, a partir da pasta raiz do projeto:

```bash
docker compose up --build
```

Aguarde todos os containers iniciarem. Você verá os logs dos três serviços no terminal.

Para parar:

```bash
docker compose down
```

---

## 🚀 Executar Manualmente (sem Docker Compose)

> Abra **5 terminais** separados e execute cada passo em um deles, na ordem abaixo.

---

### Terminal 1 — Banco de Dados (MySQL via Docker)

Na **pasta raiz** do projeto (onde está o `init.sql`), execute:

```bash
docker run -d \
  --name mysql-micro \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=minhasenha \
  -v "$(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql" \
  mysql:8
```

Aguarde alguns segundos e verifique se os bancos foram criados:

```bash
docker exec mysql-micro mysql -uroot -pminhasenha -e "SHOW DATABASES;"
```

Você deve ver `order`, `payment` e `shipping` na lista.

---

### Terminal 2 — Payment Service (porta 3001)

```bash
cd payment
```

```bash
DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3306)/payment" \
APPLICATION_PORT=3001 \
ENV=development \
go run cmd/main.go
```

✅ Esperado: `starting payment service on port 3001 ...`

> O terminal ficará bloqueado aguardando conexões. Isso é o comportamento correto.

---

### Terminal 3 — Shipping Service (porta 3002)

```bash
cd shipping
```

```bash
DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3306)/shipping" \
APPLICATION_PORT=3002 \
ENV=development \
go run cmd/main.go
```

✅ Esperado: `gRPC shipping server running on port 3002`

---

### Terminal 4 — Order Service (porta 3000)

```bash
cd order
```

```bash
DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3306)/order" \
APPLICATION_PORT=3000 \
ENV=development \
PAYMENT_SERVICE_URL="localhost:3001" \
SHIPPING_SERVICE_URL="localhost:3002" \
go run cmd/main.go
```

✅ Esperado: `gRPC server running on port 3000`

---

### Terminal 5 — Teste da Integração com grpcurl

#### Itens disponíveis no estoque (pré-cadastrados pelo `init.sql`)

| product_code | Descrição |
|---|---|
| `prod_1` | Camiseta |
| `prod_2` | Calça Jeans |
| `prod_3` | Tênis Esportivo |
| `prod_4` | Mochila |
| `prod_5` | Relógio |

> Usar um `product_code` fora desta lista retorna erro `NOT_FOUND`.

Com os três serviços rodando, dispare uma requisição de criação de pedido:

```bash
grpcurl \
  -d '{"costumer_id": 123, "order_items": [{"product_code": "prod_1", "quantity": 4, "unit_price": 12.5}], "total_price": 50.0}' \
  -plaintext localhost:3000 \
  order.Order/Create
```

Resposta esperada:
```json
{
  "orderId": 1
}
```

#### Teste direto no Shipping Service

```bash
grpcurl \
  -d '{"order_id": 1, "order_items": [{"product_code": "prod_1", "quantity": 4}]}' \
  -plaintext localhost:3002 \
  shipping.Shipping/CalculateShipping
```

Resposta esperada:
```json
{
  "deliveryDays": 3
}
```

---

## 🔍 O que observar nos logs

| Terminal | Log esperado |
|---|---|
| **Order** | `gRPC server running on port 3000` → após pedido: `Shipping calculated for order X: 3 delivery days` |
| **Payment** | `Creating payment...` ao receber a cobrança do Order |
| **Shipping** | `CalculateShipping called for order_id=1` ao receber a chamada do Order |

---

## 📁 Estrutura do Projeto (Arquitetura Hexagonal)

Todos os serviços seguem o mesmo padrão de organização:

```
<service>/
├── cmd/main.go                        # Entry point — inicializa e injeta dependências
├── config/config.go                   # Leitura de variáveis de ambiente
└── internal/
    ├── application/core/
    │   ├── domain/                    # Entidades e regras de negócio puras
    │   └── api/                       # Casos de uso (Application Services)
    ├── ports/
    │   ├── api.go                     # Inbound port (contrato da aplicação)
    │   ├── db.go                      # Outbound port (contrato do banco)
    │   └── ...                        # Outros outbound ports (payment, shipping)
    └── adapters/
        ├── grpc/                      # Servidor gRPC (entrada) e clientes gRPC (saída)
        └── db/                        # Implementação MySQL via GORM
```

### Serviços e portas

| Serviço | Porta | Banco MySQL |
|---|---|---|
| Order Service | 3000 | `order` |
| Payment Service | 3001 | `payment` |
| Shipping Service | 3002 | `shipping` |
