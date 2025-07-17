import { describe, it, expect } from 'vitest'

import { mount } from '@vue/test-utils'
import HelloWorld from '../HelloWorld.vue'
import App from '../../App.vue'

const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms));

describe("App.vue", () => {
  it("Deve criar uma conta", async () => {
    const wrapper = mount(App, {});
    const input = {
      name: "Gustavo B",
      email: `gustavo-${crypto.randomUUID()}@example.com`,
      password: "Test1234",
      document: "11144477735"
    }
    await wrapper.get(".input-name").setValue(input.name);
    await wrapper.get(".input-email").setValue(input.email);
    await wrapper.get(".input-password").setValue(input.password);
    await wrapper.get(".input-document").setValue(input.document);
    await wrapper.get(".btn-signup").trigger("click");
    await sleep(1000); 
    expect(wrapper.get(".span-message").text()).toBe("Conta criada com sucesso!");
  })
})