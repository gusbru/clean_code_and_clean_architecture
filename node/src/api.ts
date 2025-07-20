import express, { Request, Response } from "express";
import cors from "cors";
import { AccountDAODatabase } from "./DAO/accountDAO";
import { AccountService } from "./services/accountService";
import { AssetDAODatabase } from "./DAO/assetDAO";
import { AssetService } from "./services/assetService";
import { AccountAssetService } from "./services/accountAssetService";

const app = express();
app.use(express.json());
app.use(cors());
 
const accountDAO = new AccountDAODatabase();
const assetDAO = new AssetDAODatabase()
const accountService = new AccountService(accountDAO);
const assetService = new AssetService(assetDAO);
const accountAssetService = new AccountAssetService(accountService, assetService);

app.post("/signup", async (req: Request, res: Response) => {
    const account = req.body;
    try {
        const response = await accountAssetService.accountService.signup(account);
        res.status(201).json({ accountId: response.accountId });
    } catch (error: any) {
        res.status(422).json({ error: error.message });
    }
});

app.get("/accounts/:accountId", async (req: Request, res: Response) => {
    const { accountId } = req.params;
    try {
        const account = await accountAssetService.getAccountWithAssets(accountId);
        res.status(200).json(account);
    } catch (error: any) {
        res.status(404).json({ error: error.message });
    }
});

app.post("/deposit", async (req: Request, res: Response) => {
    const { accountId, assetId, quantity } = req.body;
    try {
        const response = await accountAssetService.deposit({ accountId, assetId, quantity });
        res.status(201).json(response);
    } catch (error: any) {
        res.status(422).json({ error: error.message });
    }
});

app.post("/withdraw", async (req: Request, res: Response) => {
    const { accountId, assetId, quantity } = req.body;
    try {
        const response = await accountAssetService.withdraw({ accountId, assetId, quantity });
        res.status(200).json(response);
    } catch (error: any) {
        res.status(422).json({ error: error.message });
    }
});

app.listen(3000);