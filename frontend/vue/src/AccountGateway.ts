export default interface AccountGateway {
    save(account: any): Promise<any>;
}

export class AccountGatewayHttp implements AccountGateway {
    async save(account: any): Promise<any> {
        const response = await fetch('http://localhost:3000/signup', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(account)
        });

        return await response.json();
    }
}

export class AccountGatewayMemory implements AccountGateway {
    async save(account: any): Promise<any> {
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    accountId: crypto.randomUUID(),
                    message: 'Account created successfully!'
                });
            }, 1000);
        });
    }
}