<script setup lang="ts">
  import { ref } from 'vue';

  const form = ref({
    name: '',
    email: '',
    password: '',
    document: '',
    accountId: "",
    message: ""
  })

  async function signup() {
    const input = form.value;
    const response = await fetch('http://localhost:3000/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(input)
    });
    const output = await response.json();
    if (output.accountId) {
      form.value.accountId = output.accountId;
      form.value.message = "Conta criada com sucesso!";
    } else {
      form.value.message = output.message
    }
  }
</script>

<template>
  <div>
    <div>
      <input class="input-name" type="text" v-model="form.name" placeholder="Digite seu nome" />
    </div>
  </div>

  <div>
    <div>
      <input class="input-email" type="email" v-model="form.email" placeholder="Digite seu email" />
    </div>
  </div>

  <div>
    <div>
      <input class="input-password" type="password" v-model="form.password" placeholder="Digite sua senha" />
    </div>
  </div>

  <div>
    <div>
      <input class="input-document" type="text" v-model="form.document" placeholder="Digite seu CPF" />
    </div>
    <button class="btn-signup" type="button" @click="signup">Criar conta</button>
    <div class="span-message">{{ form.message }}</div>
  </div>

</template>

<style scoped>
header {
  line-height: 1.5;
}

.logo {
  display: block;
  margin: 0 auto 2rem;
}

@media (min-width: 1024px) {
  header {
    display: flex;
    place-items: center;
    padding-right: calc(var(--section-gap) / 2);
  }

  .logo {
    margin: 0 2rem 0 0;
  }

  header .wrapper {
    display: flex;
    place-items: flex-start;
    flex-wrap: wrap;
  }
}
</style>
