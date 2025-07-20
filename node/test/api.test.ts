import crypto from "crypto"
import axios from "axios";

test("Should create an account", async () => {
    // Given
    const input = {
        name: "Gustavo Brunetto",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Senha123",
        document: "11144477735"
    }
    // When
    const responseSignup = await axios.post("http://backend:3000/signup", input);
    // Then
    const outputSignup = responseSignup.data;
    expect(responseSignup.status).toBe(201);
    expect(outputSignup).toHaveProperty("accountId");
    const accountResponse = await axios.get(`http://backend:3000/accounts/${outputSignup.accountId}`);
    const getAccountData = accountResponse.data;
    expect(getAccountData.accountId).toBe(outputSignup.accountId);
    expect(getAccountData.name).toBe(input.name);
    expect(getAccountData.email).toBe(input.email);
    expect(getAccountData.document).toBe(input.document);
    expect(getAccountData.password).toBeUndefined();
});

test.each([
    { name: "" },
    { name: "Gustavo" },
    { name: "Gustavo B C" },
])("Should not create an account with invalid name", async ({ name }) => {
    // Given
    const input = {
        name,
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Test1234",
        document: "11144477735"
    }
    try {
        // When
        await axios.post("http://backend:3000/signup", input);
        // If the request succeeds, we should fail the test
        expect(true).toBe(false);
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(422);
        expect(error.response.data.error).toBe("Invalid name format. Name must contain first and last name.");
    }
});
