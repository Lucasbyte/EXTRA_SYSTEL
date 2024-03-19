# Integração de Dados com PostgreSQL com Arquivos Toledo

Este repositório contém um script em Python para integração de dados em um banco de dados PostgreSQL. O script realiza análises em arquivos de texto e envia as informações processadas para o banco de dados.

## Funcionalidades Principais

1. **Análise de Arquivos**
   - Funções para analisar arquivos de texto e extrair informações específicas, como conservação, fracionamento, alergias, informações do fornecedor, etc.

2. **Conexão com Banco de Dados**
   - Conexão com banco de dados PostgreSQL, tanto em ambiente local quanto remoto.

3. **Atualização de Dados no Banco**
   - Funções para enviar informações analisadas para o banco de dados PostgreSQL.

4. **Integração de Dados**
   - Integração de informações de produtos com campos específicos do banco de dados.

## Uso e Execução

1. **Configuração do Ambiente**
   - Certifique-se de ter o Python instalado em sua máquina. Você pode baixar o Python em [python.org](https://www.python.org/downloads/).
   - Instale a biblioteca psycopg2 para conexão com o PostgreSQL:
     ```
     pip install psycopg2
     ```

2. **Configuração do Banco de Dados**
   - Tenha um banco de dados PostgreSQL configurado e em execução. Baixe o PostgreSQL em [postgresql.org](https://www.postgresql.org/download/).
   - Anote as credenciais do banco de dados: nome do banco, usuário, senha, host e porta.

3. **Configuração do Script**
   - Abra o arquivo `main.py` e atualize as credenciais do banco de dados nas funções `conectar_banco` e `infoSystel_writer` de acordo com suas configurações.

4. **Execução do Script**
   - Execute o script principal `main.py` para iniciar o processo de integração de dados.

## Detalhes Técnicos

- Linguagem: Python
- Dependências: psycopg2
- Banco de Dados: PostgreSQL

## Observações

Certifique-se de configurar corretamente as credenciais do banco de dados no script antes de executá-lo. Se precisar de mais assistência ou esclarecimentos sobre a configuração, não hesite em entrar em contato.
