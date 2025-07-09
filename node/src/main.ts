import express, { Request, Response } from "express";
import crypto from "crypto";
import pgPromise  from "pg-promise";

const app = express();
app.use(express.json());
// Connect to DB using pg-promise
const db = pgPromise()({
    host: "db",
    port: 5432,
    database: "app",
    user: "postgres",
    password: "postgres"
});

const isNameValid = (name: string): boolean => {
    return name.split(" ").length == 2;
}

const isEmailValid = (email: string): boolean => {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
}

const isEmailDuplicated = async (email: string): Promise<boolean> => {
    const accounts = await db.query("select * from ccca.users where email = $1", [email]);
    return accounts.length !== 0;
}

function isDocumentValid(document: string): boolean {
    if (!document) return false;
    const digits = cleanDocument(document);
    if (digits.length !== 11) return false;
    if (isDocumentAllDigitsTheSame(digits)) return false;
    const firstDigit = calculateDigit(digits, 10);
    const secondDigit = calculateDigit(digits, 11);
    return extractDocumentDigits(digits) === `${firstDigit}${secondDigit}`
}

function cleanDocument(document: string): string {
    return document.replace(/\D/g, "");
}

function isDocumentAllDigitsTheSame(document: string): boolean {
    const [firstDigit] = document;
    return [...document].every(digit => digit === firstDigit)
}

function calculateDigit(document: string, factor: number): number {
    let total = 0;
    for (const digit of document) {
        if (factor > 1) total += parseInt(digit) * factor--;
    }
    const rest = total%11;
    return (rest < 2) ? 0 : 11 - rest;
}

function extractDocumentDigits(document: string): string {
    return document.slice(9);
}

function isValidPassword(password: string): boolean {
    const regex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[A-Za-z\d]{8,}$/;
    return regex.test(password);
}

function isValidAssetId(assetId: string): boolean {
    const validAssets = ["BTC", "USD"];
    return validAssets.includes(assetId);
}

function isValidUUID(uuid: string): boolean {
    const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;
    return uuidRegex.test(uuid);
}

async function isValidAccount(accountId: string): Promise<boolean> {
    const response = await db.oneOrNone("select * from ccca.users where account_id = $1", [accountId]);
    return response?.account_id === accountId;
}

function isValidQuantity(quantity: number): boolean {
    const num = Number(quantity);
    return !isNaN(num) && num > 0 && isFinite(num);
}    

app.post("/signup", async (req: Request, res: Response) => {
    const account = req.body;
    const accountId = crypto.randomUUID();
    if (!isNameValid(account.name)) {
        res.status(400).json({ error: "Invalid name format. Name must contain first and last name." });
        return;
    }
    if (!isEmailValid(account.email)) {
        res.status(400).json({ error: "Invalid Email format." });
        return;
    }
    if (await isEmailDuplicated(account.email)) {
        res.status(400).json({ error: "Duplicated email" })
        return ;
    }
    if (!isDocumentValid(account.document)) {
        res.status(400).json({ error: "Invalid document format." });
        return;
    }
    if (!isValidPassword(account.password)) {
        res.status(400).json({ error: "Invalid password format." });
        return;
    }
    await db.query("insert into ccca.users (account_id, name, email, document, password) values ($1, $2, $3, $4, $5)", [accountId, account.name, account.email, account.document, account.password]);
    res.status(201).json({ accountId });
});

app.get("/accounts/:accountId", async (req: Request, res: Response) => {
    const { accountId } = req.params;
    const account = await db.oneOrNone("select * from ccca.users where account_id = $1", [accountId]);
    if (!account) {
        res.status(404).json({ error: "Account not found" });
        return;
    }
    res.status(200).json({
        accountId: account.account_id,
        name: account.name,
        email: account.email,
        document: account.document,
    });
});

app.post("/deposit", async (req: Request, res: Response) => {
    const { accountId, assetId, quantity } = req.body;
    if (!isValidUUID(accountId) || !(await isValidAccount(accountId))) {
        res.status(400).json({ error: "Invalid accountId" });
        return;
    }
    if (!isValidQuantity(quantity)) {
        res.status(400).json({ error: "Invalid quantity" });
        return;
    }
    if (!isValidAssetId(assetId)) {
        res.status(400).json({ error: "Invalid assetId" });
        return;
    }
    const assetExists = await db.oneOrNone("select * from ccca.assets where asset_id = $1 and account_id = $2", [assetId, accountId]);
    if (assetExists) {
        await db.query("update ccca.assets set amount = amount + $1 where asset_id = $2 and account_id = $3", [quantity, assetId, accountId]);
    } else {
        await db.query("insert into ccca.assets (asset_id, account_id, amount) values ($1, $2, $3)", [assetId, accountId, quantity]);
    }
    res.status(201).json({ message: "Deposit successful" });
});

app.post("/withdraw", async (req: Request, res: Response) => {
    const { accountId, assetId, quantity } = req.body;
    
});

app.listen(3000);