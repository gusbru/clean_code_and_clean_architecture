import IAssetDAO from "../DAO/assetDAO";
import { isValidAssetId } from "../validateAsset";
import { isValidQuantity } from "../validateQuantity";

export default interface IAssetService {
    deposit(depositRequest: any): Promise<{ message: string }>;
    withdraw(withdrawRequest: any): Promise<{ message: string }>;
    getAssetsByAccountId(accountId: string): Promise<any[]>;
    getAssetsById(accountId: string, assetId: string): Promise<any[]>;
}

export class AssetService implements IAssetService {
    constructor(private assetDAO: IAssetDAO) {}

    async deposit(depositRequest: any) {
        if (!isValidQuantity(depositRequest.quantity)) throw new Error("Invalid quantity");
        if (!isValidAssetId(depositRequest.assetId)) throw new Error("Invalid assetId");

        const assetExists = await this.assetDAO.getById(depositRequest.accountId, depositRequest.assetId);
        if (assetExists) {
            await this.assetDAO.updateQuantity(depositRequest.accountId, depositRequest.assetId, depositRequest.quantity);
        } else {
            await this.assetDAO.save(depositRequest);
        }
        return { message: "Deposit successful" };
    }

    async withdraw(withdrawRequest: any) {
        if (!isValidAssetId(withdrawRequest.assetId)) throw new Error("Invalid assetId");
        if (!isValidQuantity(withdrawRequest.quantity)) throw new Error("Invalid quantity");
        const accountAsset = await this.assetDAO.getById(withdrawRequest.accountId, withdrawRequest.assetId);
        if (!accountAsset) throw new Error("Account or asset not found");
        if (accountAsset.quantity < withdrawRequest.quantity) throw new Error("Insufficient asset quantity");

        await this.assetDAO.updateQuantity(withdrawRequest.accountId, withdrawRequest.assetId, -withdrawRequest.quantity);
        return { message: "Withdraw successful" };
    }

    async getAssetsByAccountId(accountId: string): Promise<any[]> {
        return this.assetDAO.getByAccountId(accountId);
    }

    async getAssetsById(accountId: string, assetId: string): Promise<any[]> {
        return this.assetDAO.getById(accountId, assetId);
    }
}