import IAccountService from "./accountService";
import IAssetService from "./assetService";

export interface IAccountAssetService {
    accountService: IAccountService;
    getAccountWithAssets(accountId: string): Promise<any>;
    createAccountWithInitialAssets(account: any, initialAssets?: any[]): Promise<{ accountId: string }>;
    deposit(depositRequest: any): Promise<any>;
    withdraw(withdrawRequest: any): Promise<any>;
}

export class AccountAssetService implements IAccountAssetService {
    constructor(
        public accountService: IAccountService,
        private assetService: IAssetService
    ) {}

    async getAccountWithAssets(accountId: string) {
        const account = await this.accountService.getAccountById(accountId);
        const assets = await this.assetService.getAssetsByAccountId(accountId);
        return {
            ...account,
            assets
        };
    }

    async createAccountWithInitialAssets(account: any, initialAssets: any[] = []) {
        const { accountId } = await this.accountService.signup(account);
        
        for (const asset of initialAssets) {
            await this.assetService.deposit({
                accountId,
                assetId: asset.assetId,
                quantity: asset.quantity
            });
        }
        
        return { accountId };
    }

    async deposit(depositRequest: any) {
        await this.accountService.getAccountById(depositRequest.accountId);
        return this.assetService.deposit(depositRequest);
    }

    async withdraw(withdrawRequest: any) {
        await this.accountService.getAccountById(withdrawRequest.accountId);
        return this.assetService.withdraw(withdrawRequest);
    }
}