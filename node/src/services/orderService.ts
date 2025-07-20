import IOrderDAO from "../DAO/orderDAO";
import { IAccountAssetService } from "./accountAssetService";

export default interface IOrderService {
  accountAssetService: IAccountAssetService;
  executeOrder(order: any): Promise<string>;
  getOrders(accountId: string, status?: string): Promise<any[]>;
}

export class OrderService implements IOrderService {
  constructor(
    private orderDAO: IOrderDAO,
    public accountAssetService: IAccountAssetService
  ) {}

  async executeOrder(order: any): Promise<string> {
    const account = await this.accountAssetService.accountService.getAccountById(order.account_id);
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
