# 📋 Sistema de Microsserviços de Pedidos e Pagamentos (Order & Payment)

Este projeto implementa uma arquitetura de microsserviços distribuídos utilizando a **Arquitetura Hexagonal** (Ports and Adapters) na linguagem Go. O sistema é composto por dois serviços principais que se comunicam internamente via **gRPC** e utilizam bancos de dados **MySQL** isolados através do Docker.

---

## 🏗️ Fluxo de Comunicação do Sistema

[Cliente (grpcurl/Postman)]
│
▼ (gRPC na Porta 3000)
┌──────────────┐
│  Order Svc   │ ───► [Banco MySQL: order]
└──────────────┘
│
▼ (gRPC na Porta 3001 - PAYMENT_SERVICE_URL)
┌──────────────┐
│ Payment Svc  │ ───► [Banco MySQL: payment]
└──────────────┘

1. O **Cliente** faz uma requisição de criação de pedido para o serviço **Order** (`localhost:3000`).
2. O serviço **Order** registra as informações no seu banco de dados local (`order`).
3. O serviço **Order** realiza uma chamada gRPC interna para o serviço **Payment** (`localhost:3001`) enviando os dados de cobrança.
4. O serviço **Payment** processa e persiste a transação em seu respectivo banco de dados (`payment`).
5. Se tudo der certo, o **Order** finaliza o fluxo e retorna o ID da compra para o cliente.

---

## 🛠️ Pré-requisitos

Antes de iniciar, certifique-se de ter instalado em sua máquina:
* **Go** (versão 1.18 ou superior)
* **Docker** e **Docker Compose**
* **grpcurl** (para testes de rota via terminal)

---

## 🚀 Como Executar o Sistema

Siga o passo a passo executando os comandos em terminais separados.

### 1. Subir o Banco de Dados (MySQL via Docker)
O sistema precisa de duas bases de dados (`order` e `payment`). Certifique-se de estar na pasta raiz do projeto onde está o arquivo `init.sql` e execute:

```bash
docker run -p 3306:3306 -e MYSQL_ROOT_PASSWORD=minhasenha -v "$(pwd)/init.sql:/docker-entrypoint-initdb.d/init.sql" mysql
```

2. Iniciar o Microsserviço de Pagamento (Payment Service)
Abra um segundo terminal, navegue até a pasta do serviço de pagamento e execute o comando abaixo para injetar as configurações e ligar o servidor gRPC na porta 3001:

```bash
cd payment
DB_DRIVER=mysql DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3306)/payment" APPLICATION_PORT=3001 ENV=development go run cmd/main.go
```
💡 Nota: O terminal ficará travado aguardando conexões. Isso significa que o servidor de pagamentos está online.

3. Iniciar o Microsserviço de Pedidos (Order Service)
Abra um terceiro terminal, navegue até a pasta do serviço de pedidos e informe a URL do serviço de pagamentos através da variável PAYMENT_SERVICE_URL. Ligue o servidor gRPC na porta 3000:
```bash
cd order
DB_DRIVER=mysql DATA_SOURCE_URL="root:minhasenha@tcp(127.0.0.1:3306)/order" APPLICATION_PORT=3000 ENV=development PAYMENT_SERVICE_URL="localhost:3001" go run cmd/main.go
```

## 🧪 Como Testar a Integração (gRPC)
Como o sistema utiliza gRPC, você não conseguirá testar pelo navegador. Abra um quarto terminal e dispare a requisição utilizando o grpcurl:

```bash
grpcurl -d '{"costumer_id": 123, "order_items": [{"product_code": "prod_1", "quantity": 4, "unit_price": 12.5}], "total_price": 50.0}' -plaintext localhost:3000 order.Order/Create
```

## 🔍 O que valida o sucesso do teste?
No terminal de envio: Você receberá um objeto JSON contendo o ID do pedido gerado com sucesso.

No terminal do Order: Verá os logs de criação de pedido e a chamada gRPC sendo enviada para o Payment.

No terminal do Payment: Verá os logs piscando instantaneamente, confirmando que recebeu a ordem de cobrança vinda do serviço de Order.

📁 Estrutura Organizacional (Arquitetura Hexagonal)
Ambos os projetos seguem rigidamente a divisão de responsabilidades dos Ports and Adapters:

cmd/main.go: Inicializa a aplicação, injeta as dependências e sob os adaptadores de servidor.

config/: Centraliza e isola a leitura das variáveis de ambiente de forma segura.

internal/application/core/: O "coração" da aplicação. Contém as regras de negócio puras (Entidades e Casos de Uso/APIs). Não possui dependência de frameworks ou bancos de dados.

internal/ports/: Contratos e interfaces que definem como o mundo externo pode entrar (Inbound Ports) e como o sistema se comunica com o mundo externo (Outbound Ports).

internal/adapters/: Implementações tecnológicas. Aqui ficam os drivers do banco de dados (MySQL/GORM) e os servidores/clientes de rede (gRPC).
