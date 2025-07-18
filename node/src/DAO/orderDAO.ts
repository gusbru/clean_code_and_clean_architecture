import pgPromise from "pg-promise";

export default interface IOrderDAO {
  save(order: any): Promise<string>;
  getOrders(accountId: string, status?: string): Promise<any[]>;
}

export class OrderDAODatabase implements IOrderDAO {
  async save(order: any): Promise<string> {
    const db = pgPromise()({
      host: "db",
      port: 5432,
      database: "app",
      user: "postgres",
      password: "postgres",
    });
    const [{ order_id }] = await db.query(
      "insert into ccca.order (order_id, market_id, account_id, side, quantity, price, fill_quantity, fill_price, status, timestamp) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning order_id",
      [
        order.orderId,
        order.marketId,
        order.accountId,
        order.side,
        order.quantity,
        order.price,
        order.fillQuantity,
        order.fillPrice,
        order.status,
        order.timestamp,
      ]
    );
    await db.$pool.end();
    return order_id;
  }

  async getOrders(accountId: string, status?: string): Promise<any[]> {
    const db = pgPromise()({
      host: "db",
      port: 5432,
      database: "app",
      user: "postgres",
      password: "postgres",
    });
    const query =
      "select * from ccca.order where account_id = $1" +
      (status ? " and status = $2" : "");
    const params = status ? [accountId, status] : [accountId];
    const orders = await db.query(query, params);
    console.log("Retrieved orders:", orders);
    await db.$pool.end();
    return orders.map((order: any) => ({
      orderId: order.order_id,
      marketId: order.market_id,
      accountId: order.account_id,
      side: order.side,
      quantity: parseFloat(order.quantity),
      price: parseFloat(order.price),
      fillQuantity: parseFloat(order.fill_quantity),
      fillPrice: parseFloat(order.fill_price),
      status: order.status,
      timestamp: order.timestamp,
    }));
  }
}

export class OrderDAOMemory implements IOrderDAO {
  private orders: Map<string, any> = new Map();

  async save(order: any): Promise<string> {
    this.orders.set(order.orderId, order);
    return order.orderId;
  }

  async getOrders(accountId: string, status?: string): Promise<any[]> {
    return Array.from(this.orders.values()).filter(
      (order) =>
        order.accountId === accountId && (!status || order.status === status)
    );
  }
}
