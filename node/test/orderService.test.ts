import { AccountDAOMemory } from "../src/DAO/accountDAO";
import { AssetDAOMemory } from "../src/DAO/assetDAO";
import { OrderDAOMemory } from "../src/DAO/orderDAO";
import { AccountAssetService } from "../src/services/accountAssetService";
import { AccountService } from "../src/services/accountService";
import { AssetService } from "../src/services/assetService";
import IOrderService, { OrderService } from "../src/services/orderService";

let orderService: IOrderService;

beforeEach(() => {
  const accountDAO = new AccountDAOMemory();
  const accountService = new AccountService(accountDAO);
  const assetDAO = new AssetDAOMemory();
  const assetService = new AssetService(assetDAO);
  const orderDAO = new OrderDAOMemory();
  const accountAssetService = new AccountAssetService(
    accountService,
    assetService
  );
  orderService = new OrderService(orderDAO, accountAssetService);
});

test("Should execute order with valid data", async () => {
  const newAccountInput = {
    name: "John Doe",
    email: `john-${crypto.randomUUID()}@example.com`,
    password: "Password123",
    document: "11144477735",
  };
  const { accountId } = await orderService.accountAssetService.accountService.signup(newAccountInput);
  const inputOrder = {
    market_id: "BTC/USD",
    account_id: accountId,
    side: "buy",
    quantity: 1.0,
    price: 50000.0,
    fill_quantity: 0,
    fill_price: 0,
    status: "open",
    timestamp: new Date().toISOString(),
  };
  const orderId = await orderService.executeOrder(inputOrder);
  expect(orderId).toBeDefined();
});

test("Should not execute order with nonexistent accountId", async () => {
  const inputOrder = {
    market_id: "BTC/USD",
    account_id: crypto.randomUUID(),
    side: "buy",
    quantity: 1.0,
    price: 50000.0,
    fill_quantity: 0,
    fill_price: 0,
    status: "open",
    timestamp: new Date().toISOString(),
  };
  await expect(() => orderService.executeOrder(inputOrder)).rejects.toThrow(
    "Account not found"
  );
});
