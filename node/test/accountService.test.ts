import crypto from "crypto"
import axios from "axios";
import IAccountService, { AccountService } from "../src/services/accountService";
import { AccountDAOMemory } from "../src/DAO/accountDAO";

let accountService: IAccountService;

beforeEach(() => {
    const accountDAO = new AccountDAOMemory();
    accountService = new AccountService(accountDAO);
})

test.each([
    { name: "" },
    { name: "Gustavo" },
    { name: "Gustavo B C" },
])("Should not create an account with invalid name", async ({ name }) => {
    const input = {
        name,
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Test1234",
        document: "11144477735"
    }
    await expect(() => accountService.signup(input)).rejects.toThrow("Invalid name format. Name must contain first and last name.");
});

test.each([
    { email: `gustavo-${crypto.randomUUID()}@example` },
    { email: `gustavo-${crypto.randomUUID()}@.com` },
    { email: "@example.com" },
])("Should not create an account with invalid email", async ({ email }) => {
    const input = {
        name: "Gustavo B",
        email,
        password: "Test1234",
        document: "11144477735"
    }
    await expect(() => accountService.signup(input)).rejects.toThrow("Invalid Email format.");
});

test("Should not create an account with duplicated email", async () => {
    // Given
    const input = {
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Test1234",
        document: "11144477735"
    }
    await accountService.signup(input);
    await expect(() => accountService.signup(input)).rejects.toThrow("Duplicated email")
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
    await expect(() => accountService.signup(input)).rejects.toThrow("Invalid document format.");

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
    await expect(() => accountService.signup(input)).rejects.toThrow("Invalid password format.");
});