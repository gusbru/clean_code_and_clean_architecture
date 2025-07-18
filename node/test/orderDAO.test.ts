import { OrderDAODatabase } from "../src/DAO/orderDAO";

test("Should persist an order", async () => {
    const orderDAO = new OrderDAODatabase();
    const accountId = crypto.randomUUID();
    const orderIdInput = crypto.randomUUID();
    const inputOrder = {
        order_id: orderIdInput,
        market_id: "BTC/USD",
        account_id: accountId,
        side: "buy",
        quantity: 10.0,
        price: 100.0,
        fill_quantity: 3,
        fill_price: 99.0,
        status: "open",
        timestamp: new Date().toISOString()
    }
    const orderId = await orderDAO.save(inputOrder);
    expect(orderId).not.toBeNull();
    expect(orderId).toBeDefined();
    const orders = await orderDAO.getOrders(accountId);
    expect(orders.length).toBe(1);
    expect(orders[0].order_id).toBe(orderId);
    expect(orders[0].market_id).toBe(inputOrder.market_id);
    expect(orders[0].account_id).toBe(inputOrder.account_id);
    expect(orders[0].side).toBe(inputOrder.side);
    expect(orders[0].quantity).toBe(inputOrder.quantity);
    expect(orders[0].price).toBe(inputOrder.price);
    expect(orders[0].fill_quantity).toBe(inputOrder.fill_quantity);
    expect(orders[0].fill_price).toBe(inputOrder.fill_price);
    expect(orders[0].status).toBe(inputOrder.status);
});