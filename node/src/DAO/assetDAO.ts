import pgPromise from "pg-promise";

export default interface IAssetDAO {
    save(asset: any): Promise<void>;
    getById(accountId: string, assetId: string): Promise<any>;
    getByAccountId(accountId: string): Promise<any[]>;
    updateQuantity(accountId: string, assetId: string, quantity: number): Promise<void>;
}

export class AssetDAODatabase implements IAssetDAO {
    async save(asset: any): Promise<void> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        await db.query("insert into ccca.account_asset (asset_id, account_id, quantity) values ($1, $2, $3)", [asset.assetId, asset.accountId, asset.quantity]);
        await db.$pool.end();
    }
    
    async getById(accountId: string, assetId: string): Promise<any> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        const asset = await db.oneOrNone("select * from ccca.account_asset where asset_id = $1 and account_id = $2", [assetId, accountId]);
        await db.$pool.end();
        return asset;
    }

    async getByAccountId(accountId: string): Promise<any[]> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        const assets = await db.query("select * from ccca.account_asset where account_id = $1", [accountId]);
        await db.$pool.end();
        return assets;
    }

    async updateQuantity(accountId: string, assetId: string, quantity: number): Promise<void> {
        const db = pgPromise()({
            host: "db",
            port: 5432,
            database: "app",
            user: "postgres",
            password: "postgres"
        });
        await db.query("update ccca.account_asset set quantity = quantity + $1 where asset_id = $2 and account_id = $3", [quantity, assetId, accountId]);
        await db.$pool.end();
    }
}

export class AssetDAOMemory implements IAssetDAO {
    private assets: Map<string, any> = new Map();

    async save(asset: any): Promise<void> {
        this.assets.set(`${asset.assetId}-${asset.accountId}`, asset);
    }

    async getById(accountId: string, assetId: string): Promise<any> {
        return this.assets.get(`${assetId}-${accountId}`);
    }

    async getByAccountId(accountId: string): Promise<any[]> {
        return Array.from(this.assets.values()).filter(asset => asset.accountId === accountId);
    }

    async updateQuantity(accountId: string, assetId: string, quantity: number): Promise<void> {
        const asset = this.assets.get(`${assetId}-${accountId}`);
        if (asset) {
            asset.quantity = (asset.quantity || 0) + quantity;
            this.assets.set(`${assetId}-${accountId}`, asset);
        }
    }
}