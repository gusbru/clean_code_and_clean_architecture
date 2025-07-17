import { isNameValid } from "../src/validateName";

test.each([
    { name: "" },
    { name: "Gustavo" },
    { name: "Gustavo B C" },
])("Should not create an account with invalid name", async ({ name }) => {
    expect(isNameValid(name)).toBe(false);
});