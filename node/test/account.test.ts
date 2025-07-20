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

test.each([
    { email: `gustavo-${crypto.randomUUID()}@example` },
    { email: `gustavo-${crypto.randomUUID()}@.com` },
    { email: "@example.com" },
])("Should not create an account with invalid email", async ({ email }) => {
    // Given
    const input = {
        name: "Gustavo B",
        email,
        password: "Test1234",
        document: "11144477735"
    }
    // When
    try {
        await axios.post("http://backend:3000/signup", input);
        // If the request succeeds, we should fail the test
        throw new Error("Should not create account")
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(422);
        expect(error.response.data.error).toBe("Invalid Email format.");
    }
});

test("Should not create an account with duplicated email", async () => {
    // Given
    const input = {
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Test1234",
        document: "11144477735"
    }
    await axios.post("http://backend:3000/signup", input);
    // When
    try {
        await axios.post("http://backend:3000/signup", input);
        // If the request succeeds, we should fail the test
        throw new Error("Should not create account")
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(422);
        expect(error.response.data.error).toBe("Duplicated email");
    }
});

test.each([
    { document: "" },
    { document: "123" },
    { document: "11111111111"},
    { document: "1234567890a" },
    { document: "123456789012" },
    { document: "1234567890" }
])("Should not create account for invalid document", async ({ document }) => {
    // Given
    const input = {
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Test1234",
        document
    }
    try {
        // When
        await axios.post("http://backend:3000/signup", input);
        // If the request succeeds, we should fail the test
        throw new Error("Should not create account")
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(422);
        expect(error.response.data.error).toBe("Invalid document format.");
    }

});

test.each([
    { password: "" },
    { password: "1234567" },
    { password: "123456789012as90123ds678901234fsd89012" },
    { password: "abcd1234" },
    { password: "1234ab!@" },
    { password: "a" }
])("Should not create an account with invalid password", async ({ password }) => {
    // Given
    const input = {
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password,
        document: "11144477735"
    }
    try {
        // When
        await axios.post("http://backend:3000/signup", input);
        // If the request succeeds, we should fail the test
        throw new Error("Should not create account")
    } catch (error: any) {
        // Then
        expect(error.response.status).toBe(422);
        expect(error.response.data.error).toBe("Invalid password format.");
    }
});