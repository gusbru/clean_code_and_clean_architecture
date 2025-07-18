import pgPromise from "pg-promise";

export default interface IOrderDAO {
  save(order: any): Promise<string>;
  getOrders(accountId: string, status?: string): Promise<any[]>;
  update(orderId: string, orderData: any): Promise<void>;
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
        order.order_id,
        order.market_id,
        order.account_id,
        order.side,
        order.quantity,
        order.price,
        order.fill_quantity,
        order.fill_price,
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
    await db.$pool.end();
    return orders.map((order: any) => ({
      ...order,
      quantity: parseFloat(order.quantity),
      price: parseFloat(order.price),
      fill_quantity: parseFloat(order.fill_quantity),
      fill_price: parseFloat(order.fill_price),
    }));
  }

  async update(orderId: string, orderData: any): Promise<void> {
    // Implementation for updating an existing order
    return;
  }
}

export class OrderDAOMemory implements IOrderDAO {
  private orders: Map<string, any> = new Map();

  async save(order: any): Promise<string> {
    this.orders.set(order.order_id, order);
    return order.order_id;
  }

  async getOrders(accountId: string, status?: string): Promise<any[]> {
    return Array.from(this.orders.values()).filter(
      (order) =>
        order.account_id === accountId && (!status || order.status === status)
    );
  }

  async update(orderId: string, orderData: any): Promise<void> {
    if (this.orders.has(orderId)) {
      this.orders.set(orderId, { ...this.orders.get(orderId), ...orderData });
    }
  }
}
