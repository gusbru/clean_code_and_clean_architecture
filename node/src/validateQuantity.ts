export function isValidQuantity(quantity: number): boolean {
    const num = Number(quantity);
    return !isNaN(num) && num > 0 && isFinite(num);
}
