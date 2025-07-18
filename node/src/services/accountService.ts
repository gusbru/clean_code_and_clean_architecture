import crypto from "crypto";
import { isValidPassword } from "../validatePassword";
import { isEmailValid } from "../validateEmail";
import { isNameValid } from "../validateName";
import { isDocumentValid } from "../validateCpf";
import { isValidUUID } from "../validateUUID";
import IAccountDAO from "../DAO/accountDAO";

export default interface IAccountService {
    signup(account: any): Promise<{ accountId: string }>;
    getAccountById(accountId: string): Promise<any>;
}

export class AccountService implements IAccountService {
    constructor(private accountDAO: IAccountDAO) {}
    signup = async (account: any) => {
        if (!isNameValid(account.name)) throw new Error("Invalid name format. Name must contain first and last name.");
        if (!isEmailValid(account.email)) throw new Error("Invalid Email format.");
        if (await this.accountDAO.getByEmail(account.email)) throw new Error("Duplicated email");
        if (!isDocumentValid(account.document)) throw new Error("Invalid document format.");
        if (!isValidPassword(account.password)) throw new Error("Invalid password format.");
        account.accountId = crypto.randomUUID();
        await this.accountDAO.save(account);
        return { accountId: account.accountId };
    };

    getAccountById = async (accountId: string) => {
        if (!isValidUUID(accountId)) throw new Error("Invalid accountId format.");
        const account = await this.accountDAO.getById(accountId);
        if (!account) throw new Error("Account not found");
        return {
            accountId: account.account_id,
            name: account.name,
            email: account.email,
            document: account.document,
        };
    };
}