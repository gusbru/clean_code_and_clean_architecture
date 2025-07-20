export function isNameValid(name: string): boolean {
    const regex = /^[A-Za-z]+ [A-Za-z]+$/;
    return regex.test(name);
}