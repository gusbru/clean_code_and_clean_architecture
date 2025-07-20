export function isDocumentValid(document: string): boolean {
    if (!document) return false;
    const digits = cleanDocument(document);
    if (digits.length !== 11) return false;
    if (isDocumentAllDigitsTheSame(digits)) return false;
    const firstDigit = calculateDigit(digits, 10);
    const secondDigit = calculateDigit(digits, 11);
    return extractDocumentDigits(digits) === `${firstDigit}${secondDigit}`
}

function cleanDocument(document: string): string {
    return document.replace(/\D/g, "");
}

function isDocumentAllDigitsTheSame(document: string): boolean {
    const [firstDigit] = document;
    return [...document].every(digit => digit === firstDigit)
}

function calculateDigit(document: string, factor: number): number {
    let total = 0;
    for (const digit of document) {
        if (factor > 1) total += parseInt(digit) * factor--;
    }
    const rest = total % 11;
    return (rest < 2) ? 0 : 11 - rest;
}

function extractDocumentDigits(document: string): string {
    return document.slice(9);
}