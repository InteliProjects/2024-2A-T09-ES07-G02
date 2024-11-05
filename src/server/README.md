# Como rodar a solução?

## Introdução

Este projeto utiliza Docker e Docker Compose para gerenciar e orquestrar os serviços necessários. Este documento fornece instruções para configurar e executar os serviços utilizando o Docker Compose.

## Requisitos

Certifique-se de ter os seguintes requisitos instalados:

- Docker: [Instalar Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Instalar Docker Compose](https://docs.docker.com/compose/install/)

## Estrutura do Projeto

O projeto é composto pelos seguintes serviços, cada um definido em um arquivo Dockerfile:

- **WebScraping Service**
- **TagService**
- **NLP Service**
- **SpeechToText Service**
- **CoreService**

Além disso, o projeto utiliza Kafka para comunicação entre os serviços. A configuração do Docker Compose para os serviços e o Kafka está no arquivo `docker-compose.yml`.

## Comandos Docker Compose

### Iniciar os Serviços

Para iniciar todos os serviços definidos no arquivo `docker-compose.yml`, execute o seguinte comando:

```sh
docker-compose up -d

