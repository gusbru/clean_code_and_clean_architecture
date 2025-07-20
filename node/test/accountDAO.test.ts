import crypto from "crypto";
import { AccountDAODatabase } from "../src/DAO/accountDAO"

test("Should persist an account", async () => {
    const accountDAO = new AccountDAODatabase();
    const account = {
        accountId: crypto.randomUUID(),
        name: "Gustavo B",
        email: `gustavo-${crypto.randomUUID()}@example.com`,
        password: "Test1234",
        document: "11144477735"
    };
    await accountDAO.save(account);
    const savedAccount = await accountDAO.getById(account.accountId);
    expect(savedAccount.account_id).toEqual(account.accountId);
    expect(savedAccount.name).toBe(account.name);
    expect(savedAccount.email).toBe(account.email);
    expect(savedAccount.document).toBe(account.document);
    expect(savedAccount.password).toBe(account.password);
});