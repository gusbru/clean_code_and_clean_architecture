import axios from "axios";

test("Should allow deposit with valid data", async () => {
    // Given
    const newAccountInput = {
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Senha123",
        document: "11144477735"
    };
    const accountResponse = await axios.post("http://backend:3000/signup", newAccountInput);
    const accountId = accountResponse.data.accountId;
    const input = {
        accountId,
        assetId: "BTC",
        quantity: 1,
    };
    const response = await axios.post("http://backend:3000/deposit", input);
    expect(response.status).toBe(201);
});

test.each([
    { accountId: "" },
    { accountId: "invalid-account-id" },
])("Should not allow deposit with invalid accountId", async ({ accountId }) => {
    // Given
    const input = {
        accountId,
        assetId: "BTC",
        quantity: 10,
    };
    try {
        // When
        await axios.post("http://backend:3000/deposit", input);
        expect(true).toBe(false);
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(422);
        expect(error.response.data.error).toBe("Invalid accountId format.");
    }
});