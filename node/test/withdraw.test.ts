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
    expect(response.data).toEqual({
        accountId,
        assetId: "BTC",
        quantity: 10,
    });
});