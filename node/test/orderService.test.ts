import { AccountDAOMemory } from "../src/DAO/accountDAO";
import { OrderDAOMemory } from "../src/DAO/orderDAO";
import IAccountService, {
  AccountService,
} from "../src/services/accountService";
import IOrderService, { OrderService } from "../src/services/orderService";

let orderService: IOrderService;
let accountService: IAccountService;

beforeEach(() => {
  const accountDAO = new AccountDAOMemory();
  accountService = new AccountService(accountDAO);
  const orderDAO = new OrderDAOMemory();
  orderService = new OrderService(orderDAO, accountService);
});

test("Should execute order with valid data", async () => {
  const newAccountInput = {
    name: "John Doe",
    email: `john-${crypto.randomUUID()}@example.com`,
    password: "Password123",
    document: "11144477735",
  };
  const { accountId } = await accountService.signup(newAccountInput);
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
