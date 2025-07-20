# Account Domain Class Diagram

This diagram shows the clean architecture implementation of the Account domain following the Repository Pattern and Dependency Inversion Principle.

```mermaid
classDiagram
    class Account {
        +string AccountID
        +string Name
        +string Email
        +string Document
        +string Password
    }

    class IAccountDAO {
        <<interface>>
        +Save(account *Account) error
        +GetByID(accountID string) (*Account, error)
        +GetByEmail(email string) (*Account, error)
    }

    class AccountDAODatabase {
        -getDB() (*sql.DB, error)
        +Save(account *Account) error
        +GetByID(accountID string) (*Account, error)
        +GetByEmail(email string) (*Account, error)
    }

    class AccountService {
        -IAccountDAO accountDAO
        +CreateAccount(input AccountInput) (*Account, error)
        +GetAccount(accountID string) (*Account, error)
        +ValidateCredentials(email, password string) (*Account, error)
    }

    %% Relationships following Dependency Inversion Principle
    IAccountDAO <|.. AccountDAODatabase : implements
    AccountService --> IAccountDAO : depends on abstraction
    AccountDAODatabase --> Account : persists

    %% Clean Architecture Annotations
    class IAccountDAO {
        <<Repository Pattern>>
        <<Abstraction Layer>>
    }
    
    class AccountDAODatabase {
        <<Infrastructure Layer>>
        <<PostgreSQL Implementation>>
    }
    
    class AccountService {
        <<Application Layer>>
        <<Business Logic>>
    }
```

## Architecture Principles Demonstrated

### Dependency Inversion Principle (DIP)
- `AccountService` depends on the `IAccountDAO` interface, not concrete implementations
- Business logic remains stable while storage implementations can vary

### Repository Pattern
- `IAccountDAO` abstracts data persistence concerns
- Enables testing with `AccountDAOMemory` and production with `AccountDAODatabase`

### Single Responsibility Principle
- **Account**: Domain entity with clear data structure
- **IAccountDAO**: Data access contract definition
- **AccountService**: Business logic orchestration
- Each DAO implementation handles one storage mechanism