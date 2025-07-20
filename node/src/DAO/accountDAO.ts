import pgPromise from "pg-promise";


export default interface IAccountDAO {
    save (account: any): Promise<void>;
    getById (accountId: string): Promise<any>;
    getByEmail(email: string): Promise<any>;
}

export class AccountDAODatabase implements IAccountDAO {
    async save(account: any): Promise<void> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        await db.query("insert into ccca.account (account_id, name, email, document, password) values ($1, $2, $3, $4, $5)", [account.accountId, account.name, account.email, account.document, account.password]);
        await db.$pool.end();
    }
    
    async getById(accountId: string): Promise<any> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        const account = await db.oneOrNone("select * from ccca.account where account_id = $1", [accountId]);
        await db.$pool.end();
        return account
    }

    async getByEmail(email: string): Promise<any> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        const account = await db.oneOrNone("select * from ccca.account where email = $1", [email]);
        await db.$pool.end();
        return account;
    }
}

export class AccountDAOMemory implements IAccountDAO {
    private accounts: Map<string, any> = new Map();
    private emailIndex: Map<string, string> = new Map();

    async save(account: any): Promise<void> {
        this.accounts.set(account.accountId, account);
        this.emailIndex.set(account.email, account.accountId);
    }

    async getById(accountId: string): Promise<any> {
        return this.accounts.get(accountId);
    }

    async getByEmail(email: string): Promise<any> {
        const accountId = this.emailIndex.get(email);
        if (accountId) {
            return this.accounts.get(accountId);
        }
        return null;
    }

}