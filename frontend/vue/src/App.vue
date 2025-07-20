<script setup lang="ts">
  import { inject, ref } from 'vue';
import type AccountGateway from './AccountGateway';

  const form = ref({
    name: '',
    email: '',
    password: '',
    document: '',
    accountId: "",
    message: ""
  })
  const accountGateway = inject("accountGateway") as AccountGateway;

  async function signup() {
    const input = form.value;
    if (!input.name || !input.email || !input.password || !input.document) {
      form.value.message = "All fields are required!";
      return;
    }
    try {
      const output = await accountGateway.save(input);
      if (output.accountId) {
        form.value.accountId = output.accountId;
        form.value.message = "Account created successfully!";
      } else {
        form.value.message = output.message || output.error || "An unexpected error occurred.";
      }
    } catch (error: any) {
      form.value.message = error instanceof Error ? error.message : "An unexpected error occurred.";
      console.error(error);
      return;
    }
  }

  function getMessageClass() {
    if (form.value.message.includes('successfully')) {
      return 'success';
    } else if (form.value.message && !form.value.accountId) {
      return 'error';
    }
    return '';
  }
</script>

<template>
  <div class="container">
    <div class="form-group">
      <label for="name">Name</label>
      <input class="input-name" type="text" v-model="form.name" placeholder="Enter your name" required />
    </div>

    <div class="form-group">
      <label for="email">Email</label>
      <input class="input-email" type="email" v-model="form.email" placeholder="Enter your email" required />
    </div>

    <div class="form-group">
      <label for="password">Password</label>
      <input class="input-password" type="password" v-model="form.password" placeholder="Enter your password" required />
    </div>

    <div class="form-group">
      <label for="document">Document</label>
      <input class="input-document" type="text" v-model="form.document" placeholder="Enter your document" required />
    </div>

    <button class="btn-signup" type="button" @click="signup">Create Account</button>
    <div class="span-message" :class="getMessageClass()">{{ form.message }}</div>
  </div>
</template>

<style scoped>
.container {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
  background: #ffffff;
  border-radius: 12px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group:last-of-type {
  margin-bottom: 2rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  color: #374151;
  letter-spacing: 0.025em;
}

.input-name,
.input-email,
.input-password,
.input-document {
  width: 100%;
  padding: 0.875rem 1rem;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  font-size: 1rem;
  line-height: 1.5;
  color: #374151;
  background-color: #ffffff;
  transition: all 0.2s ease-in-out;
  box-sizing: border-box;
}

.input-name:focus,
.input-email:focus,
.input-password:focus,
.input-document:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.input-name::placeholder,
.input-email::placeholder,
.input-password::placeholder,
.input-document::placeholder {
  color: #9ca3af;
}

.btn-signup {
  width: 100%;
  background-color: #3b82f6;
  color: #ffffff;
  border: none;
  border-radius: 8px;
  padding: 0.875rem 1rem;
  font-size: 1rem;
  font-weight: 600;
  line-height: 1.5;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  margin-bottom: 1rem;
}

.btn-signup:hover {
  background-color: #2563eb;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.btn-signup:active {
  background-color: #1d4ed8;
  transform: translateY(0);
  box-shadow: 0 2px 4px rgba(59, 130, 246, 0.3);
}

.btn-signup:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.span-message {
  padding: 0.75rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  text-align: center;
  min-height: 1rem;
  transition: all 0.2s ease-in-out;
}

.span-message:empty {
  padding: 0;
  min-height: 0;
}

.span-message:not(:empty) {
  background-color: #f0f9ff;
  color: #0369a1;
  border: 1px solid #bae6fd;
}

/* Success state for message */
.span-message.success {
  background-color: #f0fdf4;
  color: #166534;
  border: 1px solid #bbf7d0;
}

/* Error state for message */
.span-message.error {
  background-color: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
}

/* Form validation states */
.input-name:invalid,
.input-email:invalid,
.input-password:invalid,
.input-document:invalid {
  border-color: #ef4444;
}

.input-name:invalid:focus,
.input-email:invalid:focus,
.input-password:invalid:focus,
.input-document:invalid:focus {
  box-shadow: 0 0 0 3px rgba(239, 68, 68, 0.1);
}

/* Responsive design */
@media (max-width: 480px) {
  .container {
    margin: 1rem;
    padding: 1.5rem;
  }
  
  .form-group label {
    font-size: 0.8125rem;
    margin-bottom: 0.375rem;
  }
  
  .input-name,
  .input-email,
  .input-password,
  .input-document,
  .btn-signup {
    padding: 0.75rem;
    font-size: 0.875rem;
  }
}

/* Focus states for better accessibility */
.form-group label:focus-within {
  color: #3b82f6;
}

/* Loading state for button */
.btn-signup:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}
</style>
