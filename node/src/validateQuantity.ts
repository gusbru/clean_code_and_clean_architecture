export function isValidQuantity(quantity: number): boolean {
    const num = quantity;
    return !isNaN(num) && num > 0 && isFinite(num);
}
