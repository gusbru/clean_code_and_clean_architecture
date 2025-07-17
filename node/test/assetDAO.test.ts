import { AssetDAODatabase } from "../src/DAO/assetDAO"

test("Should persist an asset", async () => {
    const AssetDAO = new AssetDAODatabase();
    const asset = {
        accountId: crypto.randomUUID(),
        assetId: "BTC",
        quantity: 10,
    };
    await AssetDAO.save(asset);
    const savedAsset = await AssetDAO.getById(asset.accountId, asset.assetId);
    expect(savedAsset.asset_id).toEqual(asset.assetId);
    expect(parseInt(savedAsset.quantity)).toBe(asset.quantity);
    expect(savedAsset.account_id).toBe(asset.accountId);     
});