import IOrderDAO from "../DAO/orderDAO";
import IAccountService from "./accountService";
import IAssetService from "./assetService";

export default interface IOrderService {
  executeOrder(order: any): Promise<string>;
  getOrders(accountId: string, status?: string): Promise<any[]>;
}

export class OrderService implements IOrderService {
  constructor(
    private orderDAO: IOrderDAO,
    private accountService: IAccountService
  ) {}

  async executeOrder(order: any): Promise<string> {
    const account = await this.accountService.getAccountById(order.account_id);
    order.accountId = account.accountId;
    order.orderId = crypto.randomUUID();
    this.validateOrder(order, account);
    return await this.orderDAO.save(order);
  }

  async getOrders(accountId: string, status?: string): Promise<any[]> {
    return this.orderDAO.getOrders(accountId, status);
  }

  private async validateOrder(order: any, account: any): Promise<void> {
    
  }
}
