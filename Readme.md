# Clean Code and Clean Architecture - Trading Platform

A multi-language implementation of a trading platform API demonstrating clean code principles and clean architecture patterns. This project includes both Node.js/TypeScript and Go implementations of the same API.

## Project Overview

This trading platform provides basic account management and asset trading functionality including:
- User account creation with validation
- Asset deposits and withdrawals
- Account balance management
- Support for multiple cryptocurrencies (BTC, USD)

## Architecture

The project follows clean architecture principles with:
- **Domain Layer**: Business logic and validation rules
- **Application Layer**: Use cases and business workflows
- **Infrastructure Layer**: Database access and external dependencies
- **Presentation Layer**: HTTP API endpoints

## Project Structure

```
├── node/                    # Node.js/TypeScript implementation
│   ├── src/
│   │   └── main.ts         # Main application entry point
│   ├── test/               # Integration tests
│   ├── docker/             # Docker configuration
│   └── database/           # Database schema
├── go/                     # Go implementation
│   ├── cmd/api/
│   │   └── main.go         # Main application entry point
│   ├── internal/types/     # Type definitions
│   ├── tests/              # Integration tests
│   ├── docker/             # Docker configuration
│   └── database/           # Database schema
└── .vscode/                # VS Code configuration
```

## API Endpoints

### Authentication
- `POST /signup` - Create a new user account

### Account Management
- `GET /accounts/{accountId}` - Get account details and assets

### Trading Operations
- `POST /deposit` - Deposit assets to an account
- `POST /withdraw` - Withdraw assets from an account

## Getting Started

### Prerequisites
- Docker and Docker Compose
- Node.js 18+ (for Node.js version)
- Go 1.21+ (for Go version)

### Running the Node.js Version

1. Navigate to the node directory:
```bash
cd node
```

2. Start the services:
```bash
npm run compose:up:build
```

3. Run tests:
```bash
npm run test
```

### Running the Go Version

1. Navigate to the go directory:
```bash
cd go
```

2. Start the services:
```bash
make docker-up
```

3. Run tests:
```bash
make test
```

## API Usage Examples

### Create Account
```bash
curl -X POST http://localhost:3000/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "SecurePass123",
    "document": "11144477735"
  }'
```

### Deposit Assets
```bash
curl -X POST http://localhost:3000/deposit \
  -H "Content-Type: application/json" \
  -d '{
    "accountId": "your-account-id",
    "assetId": "BTC",
    "quantity": 10
  }'
```

### Withdraw Assets
```bash
curl -X POST http://localhost:3000/withdraw \
  -H "Content-Type: application/json" \
  -d '{
    "accountId": "your-account-id",
    "assetId": "BTC",
    "quantity": 5
  }'
```

### Get Account Details
```bash
curl -X GET http://localhost:3000/accounts/{accountId}
```

## Validation Rules

### Account Creation
- **Name**: Must contain first and last name (exactly 2 words)
- **Email**: Must be a valid email format
- **Password**: Minimum 8 characters with at least one uppercase letter, one lowercase letter, and one number
- **Document**: Must be a valid 11-digit Brazilian CPF format

### Asset Operations
- **Asset ID**: Must be one of the supported assets (BTC, USD)
- **Quantity**: Must be a positive number
- **Account ID**: Must be a valid UUID and existing account

## Database Schema

The application uses PostgreSQL with the following main tables:

- `ccca.account` - User account information
- `ccca.account_asset` - Asset balances per account
- `ccca.order` - Trading orders (for future implementation)
- `ccca.trade` - Trade history (for future implementation)

## Development

### Node.js Commands
```bash
# Start development server
npm run start:dev

# Run tests
npm run test

# Docker operations
npm run compose:up
npm run compose:down
npm run compose:logs
```

### Go Commands
```bash
# Build and run
make run

# Development mode
make dev

# Run tests
make test

# Docker operations
make docker-up
make docker-down
make docker-logs
```

## Testing

Both implementations include comprehensive integration tests covering:
- Account creation with various validation scenarios
- Asset deposit and withdrawal operations
- Error handling and edge cases
- Database integration

Tests use Docker containers to ensure consistent testing environments.

## Clean Code Principles Applied

1. **Meaningful Names**: Functions and variables have descriptive names
2. **Single Responsibility**: Each function has a single, well-defined purpose
3. **Error Handling**: Proper error handling with meaningful error messages
4. **Validation**: Input validation with clear business rules
5. **Separation of Concerns**: Database, business logic, and API layers are separated
6. **Testing**: Comprehensive test coverage for all endpoints

## Future Enhancements

- [ ] Order placement and matching engine
- [ ] Real-time price feeds
- [ ] Advanced trading features
- [ ] Authentication and authorization
- [ ] Rate limiting
- [ ] Monitoring and metrics
- [ ] API documentation with OpenAPI/Swagger

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

This project is for educational purposes demonstrating clean code and clean architecture principles.
