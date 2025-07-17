import { isDocumentValid } from "../src/validateCpf";

test.each([
    { document: "" },
    { document: "123" },
    { document: "11111111111"},
    { document: "1234567890a" },
    { document: "123456789012" },
    { document: "1234567890" }
])("Should not create account for invalid document", async ({ document }) => {
    expect(isDocumentValid(document)).toBe(false);
});
