import { OrderDAODatabase, OrderDAOMemory } from "../src/DAO/orderDAO";

describe.each([
    ['Database', () => new OrderDAODatabase()],
    ['Memory', () => new OrderDAOMemory()],
])('%s DAO', (daoType, createDAO) => {
    test("Should persist an order", async () => {
        const orderDAO = createDAO();
        const accountId = crypto.randomUUID();
        const orderIdInput = crypto.randomUUID();
        const inputOrder = {
            orderId: orderIdInput,
            marketId: "BTC/USD",
            accountId: accountId,
            side: "buy",
            quantity: 10.0,
            price: 100.0,
            fillQuantity: 3,
            fillPrice: 99.0,
            status: "open",
            timestamp: new Date().toISOString()
        }
        const orderId = await orderDAO.save(inputOrder);
        expect(orderId).not.toBeNull();
        expect(orderId).toBeDefined();
        const orders = await orderDAO.getOrders(accountId);
        expect(orders.length).toBe(1);
        expect(orders[0].orderId).toBe(orderId);
        expect(orders[0].marketId).toBe(inputOrder.marketId);
        expect(orders[0].accountId).toBe(inputOrder.accountId);
        expect(orders[0].side).toBe(inputOrder.side);
        expect(orders[0].quantity).toBe(inputOrder.quantity);
        expect(orders[0].price).toBe(inputOrder.price);
        expect(orders[0].fillQuantity).toBe(inputOrder.fillQuantity);
        expect(orders[0].fillPrice).toBe(inputOrder.fillPrice);
        expect(orders[0].status).toBe(inputOrder.status);
    });

    test("Should retrieve orders by status", async () => {
        const orderDAO = createDAO();
        const accountId = crypto.randomUUID();
        const inputOrder1 = {
            orderId: crypto.randomUUID(),
            marketId: "ETH/USD",
            accountId: accountId,
            side: "sell",
            quantity: 5.0,
            price: 200.0,
            fillQuantity: 2,
            fillPrice: 199.0,
            status: "filled",
            timestamp: new Date().toISOString()
        };
        const inputOrder2 = {
            orderId: crypto.randomUUID(),
            marketId: "ETH/USD",
            accountId: accountId,
            side: "buy",
            quantity: 3.0,
            price: 210.0,
            fillQuantity: 1,
            fillPrice: 209.0,
            status: "open",
            timestamp: new Date().toISOString()
        };
        await orderDAO.save(inputOrder1);
        await orderDAO.save(inputOrder2);
        
        const filledOrders = await orderDAO.getOrders(accountId, "filled");
        expect(filledOrders.length).toBe(1);
        expect(filledOrders[0].orderId).toBe(inputOrder1.orderId);

        const openOrders = await orderDAO.getOrders(accountId, "open");
        expect(openOrders.length).toBe(1);
        expect(openOrders[0].orderId).toBe(inputOrder2.orderId);
    });
});