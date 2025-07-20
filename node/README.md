# Clean Code and Clean Architecture

Clean Code e Clean Architecture - Turma 22

## Class 1

Project: Trade platform
* Venda mais barata.
* Compra mais cara.
* Modelo de dominio: encaixa bem com a plataform de trade.
* Vamos priorizar a questao da performance.
  * 1000 orders por segundo.
  * vamos usar bots.

Code smell:
* nome estranho (de variable, metodo, classe, etc)
* ruido: comentarios, codigo morto, linhas em branco, 


Signup
    Endpoint: POST /signup
    Input: name, email, documents, password
    Output: account_id

* Name deve ser composto de nome e sobrenome
* document

## Class 2

### Teste de integracao: 

* teste a integracao entre duas ou mais camadas do sistema.
* mais lentos devido a integracao com o sistema externo.
