import { test, expect } from '@playwright/test';

test('Deve criar uma conta', async ({ page }) => {
  await page.goto('http://localhost:5173/');

  const input = {
    name: "Gustavo B",
    email: `gustavo-${crypto.randomUUID()}@example.com`,
    password: "Test1234",
    document: "11144477735"
  };
  await page.locator('.input-name').fill(input.name);
  await page.locator('.input-email').fill(input.email);
  await page.locator('.input-password').fill(input.password);
  await page.locator('.input-document').fill(input.document);
  await page.locator('.btn-signup').click();
  await page.waitForTimeout(1000); // Espera a mensagem ser atualizada
  await expect(page.locator('.span-message')).toHaveText('Account created successfully!');
});
