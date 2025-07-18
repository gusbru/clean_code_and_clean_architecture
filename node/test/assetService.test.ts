import IAssetService, { AssetService } from "../src/services/assetService";
import { AssetDAOMemory } from "../src/DAO/assetDAO";
import { AccountDAOMemory } from "../src/DAO/accountDAO";
import IAccountService, {
  AccountService,
} from "../src/services/accountService";
import { AccountAssetService, IAccountAssetService } from "../src/services/accountAssetService";

let accountAssetService: IAccountAssetService;

beforeEach(() => {
  const accountDAO = new AccountDAOMemory();
  const accountService = new AccountService(accountDAO);
  const assetDAO = new AssetDAOMemory();
  const assetService = new AssetService(assetDAO);
  accountAssetService = new AccountAssetService(accountService, assetService);
});

test("Should allow withdrawal with valid data", async () => {
  const newAccountInput = {
    name: "Gustavo B",
    email: `gustavo-${crypto.randomUUID()}@example.com`,
    password: "Senha123",
    document: "11144477735",
  };
  const { accountId } = await accountAssetService.createAccountWithInitialAssets(newAccountInput, [
    { assetId: "BTC", quantity: 10 },
  ]);
  const inputWithdraw = {
    accountId,
    assetId: "BTC",
    quantity: 10,
  };
  const { message } = await accountAssetService.withdraw(inputWithdraw);
  expect(message).toBe("Withdraw successful");
});

test("Should not allow withdrawal with invalid accountId", async () => {
  const input = {
    accountId: "invalid-account-id",
    assetId: "BTC",
    quantity: 10,
  };
  await expect(() => accountAssetService.withdraw(input)).rejects.toThrow(
    "Invalid accountId"
  );
});

test("Should not allow withdrawal with non-existent account", async () => {
  const input = {
    accountId: "6f813af6-f151-4cbf-a423-6135909daa51",
    assetId: "BTC",
    quantity: 10,
  };
  await expect(() => accountAssetService.withdraw(input)).rejects.toThrow(
    "Account not found"
  );
});

test.each([
  { assetId: "" },
  { assetId: "invalid-asset-id" },
  { assetId: "EUR" },
])("Should not allow withdrawal with invalid assetId", async ({ assetId }) => {
  const newAccountInput = {
    name: "Gustavo B",
    email: `gustavo-${crypto.randomUUID()}@example.com`,
    password: "Senha123",
    document: "11144477735",
  };
  const { accountId } = await accountAssetService.createAccountWithInitialAssets(newAccountInput, [
    { assetId: "BTC", quantity: 10 },
  ]);
  const inputWithdraw = {
    accountId,
    assetId,
    quantity: 10,
  };
  await expect(() => accountAssetService.withdraw(inputWithdraw)).rejects.toThrow(
    "Invalid assetId"
  );
});

test("Should allow deposit with valid data", async () => {
  const newAccountInput = {
    name: "Gustavo B",
    email: `gustavo-${crypto.randomUUID()}@example.com`,
    password: "Senha123",
    document: "11144477735",
  };
  const { accountId } = await accountAssetService.createAccountWithInitialAssets(newAccountInput, [
    { assetId: "BTC", quantity: 10 },
  ]);
  const inputDeposit = {
    accountId,
    assetId: "BTC",
    quantity: 5,
  };
  const { message } = await accountAssetService.deposit(inputDeposit);
  expect(message).toBe("Deposit successful");
});

test.each([{ accountId: "" }, { accountId: "invalid-account-id" }])(
  "Should not allow deposit with invalid accountId",
  async ({ accountId }) => {
    const input = {
      accountId,
      assetId: "BTC",
      quantity: 10,
    };
    await expect(() => accountAssetService.deposit(input)).rejects.toThrow(
      "Invalid accountId"
    );
  }
);

test.each([
  { assetId: "" },
  { assetId: "invalid-asset-id" },
  { assetId: "EUR" },
])("Should not allow deposit with invalid assetId", async ({ assetId }) => {
  const newAccountInput = {
    name: "Gustavo B",
    email: `gustavo-${crypto.randomUUID()}@example.com`,
    password: "Senha123",
    document: "11144477735",
  };
  const { accountId } = await accountAssetService.accountService.signup(newAccountInput);
  const input = {
    accountId,
    assetId,
    quantity: 10,
  };
  await expect(() => accountAssetService.deposit(input)).rejects.toThrow(
    "Invalid assetId"
  );
});

test.each([
  { quantity: 0 },
  { quantity: -1 },
  { quantity: "invalid-quantity" },
])("Should not allow deposit with invalid quantity", async ({ quantity }) => {
  // Given
  const newAccountInput = {
    name: "Gustavo B",
    email: `gustavo-${crypto.randomUUID()}@example.com`,
    password: "Senha123",
    document: "11144477735",
  };
  const { accountId } = await accountAssetService.accountService.signup(newAccountInput);
  const input = {
    accountId,
    assetId: "BTC",
    quantity,
  };
  await expect(() => accountAssetService.deposit(input)).rejects.toThrow(
    "Invalid quantity"
  );
});
