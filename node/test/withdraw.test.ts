import axios from "axios";

test("Should allow withdrawal with valid data", async () => {
    // Given
    const newAccountInput = {
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Senha123",
        document: "11144477735"
    };
    const accountResponse = await axios.post("http://backend:3000/signup", newAccountInput);
    const accountId = accountResponse.data.accountId;
    const inputDeposit = {
        accountId,
        assetId: "BTC",
        quantity: 10,
    }
    await axios.post("http://backend:3000/deposit", inputDeposit);
    const input = {
        accountId,
        assetId: "BTC",
        quantity: 10,
    };
    // When
    const response = await axios.post("http://backend:3000/withdraw", input);
    // Then
    expect(response.status).toBe(200);
    const accountRequest = await axios.get(`http://backend:3000/accounts/${accountId}`);

});

test("Should not allow withdrawal with invalid accountId", async () => {
    // Given
    const input = {
        accountId: "invalid-account-id",
        assetId: "BTC",
        quantity: 10,
    };
    // When
    try {
        await axios.post("http://backend:3000/withdraw", input);
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(400);
        expect(error.response.data.error).toBe("Invalid accountId");
    }
});

test("Should not allow withdrawal with non-existent account", async () => {
    // Given
    const input = {
        accountId: "6f813af6-f151-4cbf-a423-6135909daa51",
        assetId: "BTC",
        quantity: 10,
    };
    // When
    try {
        await axios.post("http://backend:3000/withdraw", input);
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(400);
        expect(error.response.data.error).toBe("Account not found");
    }
});

test.each([
    { assetId: "" },
    { assetId: "invalid-asset-id" },
    { assetId: "EUR" },
])("Should not allow withdrawal with invalid assetId", async ({ assetId }) => {
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
        assetId,
        quantity: 10,
    };
    try {
        // When
        await axios.post("http://backend:3000/withdraw", input);
        expect(true).toBe(false);
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(400);
        expect(error.response.data.error).toBe("Invalid assetId");
    }
});