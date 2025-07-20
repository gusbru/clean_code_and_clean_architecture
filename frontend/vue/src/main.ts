import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import { AccountGatewayHttp } from './AccountGateway';

const app = createApp(App);
app.provide("accountGateway", new AccountGatewayHttp());
app.mount('#app')
