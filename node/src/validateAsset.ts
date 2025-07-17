export function isValidAssetId(assetId: string): boolean {
    const validAssets = ["BTC", "USD"];
    return validAssets.includes(assetId);
}