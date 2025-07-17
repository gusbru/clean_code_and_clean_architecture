import { isValidPassword } from "../src/validatePassword";

test.each([
    { password: "" },
    { password: "1234567" },
    { password: "123456789012as90123ds678901234fsd89012" },
    { password: "abcd1234" },
    { password: "1234ab!@" },
    { password: "a" }
])("Should not create an account with invalid password", async ({ password }) => {
    const isValid = isValidPassword(password);
    expect(isValid).toBe(false);
});