# Law Hunter

# Inteli - Instituto de Tecnologia e Liderança

<p align="center">
<a href= "https://www.inteli.edu.br/"><img src="images/logo-Inteli.png" alt="Inteli - Instituto de Tecnologia e Liderança" border="0" width=40% height=40%></a>
</p>

<br>

# Automação com Reconhecimento de Voz

## Grupo 2 - Law Hunter

## 👨‍🎓 Integrantes:
- <a href="https://www.linkedin.com/in/anna-aragao/">Anna Aragão</a>
- <a href="https://www.linkedin.com/in/davi-ferreira-arantes/">Davi Ferreira Arantes</a>
- <a href="https://www.linkedin.com/in/kaiane-souza/">Kaiane Souza</a>
- <a href="https://www.linkedin.com/in/kaleb-carvalho/">Kaleb Carvalho</a>
- <a href="https://www.linkedin.com/in/omatheusrsantos/">Matheus R. dos Santos</a>
- <a href="https://www.linkedin.com/in/rafael-coutinho2004/">Rafael Coutinho</a>

## 👩‍🏫 Professores:
### Orientador(a)
- Hermano Peixoto

### Instrutores
- Reginaldo Arakaki
- Victor Hayashi
- Hermano Peixoto
- Francisco Escobar
- Geraldo Magela
- Lisane Valdo
- Ana Cristina dos Santos

## 📜 Descrição

O monitoramento constante das leis financeiras é um desafio crucial para empresas e instituições que operam no setor financeiro. Com a rápida evolução das regulamentações, a conformidade legal torna-se não apenas uma obrigação, mas também um fator estratégico para garantir a continuidade dos negócios e evitar penalidades. Nesse contexto, a utilização de técnicas de Processamento de Linguagem Natural (NLP, na sigla em inglês) surge como uma solução inovadora e eficaz para automatizar e otimizar a análise e o acompanhamento de mudanças legislativas.

Este projeto visa desenvolver um sistema baseado em NLP que automatiza o processo de monitoramento, interpretação e categorização de novas leis e regulamentos financeiros. O objetivo principal é permitir que o LawHunter possa receber alertas sobre mudanças relevantes, facilitando a tomada de decisão e garantindo a conformidade contínua com as exigências legais. O sistema buscará e analisará informações de fontes confiáveis, classificando e resumindo as leis de acordo com seu impacto e relevância para o negócio.


## 📁 Estrutura de pastas

|--> images<br>
|--> docs<br>
  &emsp;| --> images <br>
  &emsp;|--> index.md<br>
|--> src<br>
  &emsp;|--> notebooks<br>
  &emsp;|--> server<br>
  &emsp;|--> web<br>
| README.md<br>

Dentre os arquivos e pastas presentes na raiz do projeto, definem-se:

- <b>images</b>: aqui estão os arquivos relacionados a parte gráfica do projeto, ou seja, as imagens e links de vídeos que os representam (o logo do grupo pode ser adicionado nesta pasta).

- <b>docs</b>: aqui estão todos os documentos do projeto. Há também um arquivo README para o grupo registrar a localização de cada artefato.

- <b>src</b>: Todo o código fonte criado para o desenvolvimento do projeto, incluindo backend e frontend se aplicáveis.

- <b>README.md</b>: arquivo que serve como guia e explicação geral sobre o projeto (o mesmo que você está lendo agora).

## 🔧 Instalação

&emsp;&emsp;Nesta seção, apresentaremos as instruções necessárias para configurar o ambiente de desenvolvimento do nosso projeto, que utiliza TypeScript junto com Tailwind CSS para o front-end e integrações entre o back-end e o front-end para funcionalidades como Speech to Text. Para essa integração, usamos contêineres Docker, o que permite maior consistência e facilidade na gestão dos serviços, além de facilitar o desenvolvimento local.

**Instalação de Dependências e Ferramentas Básicas**

**Docker**

1. Certifique-se de que o Docker esteja em execução no seu sistema:

```
docker --version
```

2. No diretório onde está o arquivo `docker-compose.yml` (no diretório server), execute o seguinte comando para levantar o contêiner de Speech to Text:

```
docker-compose up -d
```

## 🗃 Histórico de lançamentos

* 0.1.0 - 16/08/2024
    * Sprint 1 : Entendimento de negócio e arquitetura proposta.
* 0.2.0 - 30/08/2024
    * Sprint 2 : Primeira versão do serviço de NPL e classificação de documentos.
* 0.3.0 - 13/10/2024
    * Sprint 3 : Construção do Backend da solução
* 0.4.0 - 27/10/2024
    * Sprint 4 : Construção do Frontend da solução
* 0.5.0 - 10/11/2024
    * Sprint 5 : Finalização do projeto


## 📋 Licença/License

<img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/cc.svg?ref=chooser-v1"><img style="height:22px!important;margin-left:3px;vertical-align:text-bottom;" src="https://mirrors.creativecommons.org/presskit/icons/by.svg?ref=chooser-v1"><p xmlns:cc="http://creativecommons.org/ns#" xmlns:dct="http://purl.org/dc/terms/"><a property="dct:title" rel="cc:attributionURL" href="https://github.com/Inteli-College/2024-T0010-SI05-G01">LawHunter <a> by </a> <a rel="cc:attributionURL dct:creator" property="cc:attributionName" href="https://github.com/InteliProjects/.github/blob/main/profile/README.md">Inteli, <a href="https://www.linkedin.com/in/anna-aragao/">Anna Aragão</a>, <a href="https://www.linkedin.com/in/davi-ferreira-arantes/">Davi Ferreira Arantes</a>,  <a href="https://www.linkedin.com/in/kaiane-souza/">Kaiane Souza</a>, <a href="https://www.linkedin.com/in/kaleb-carvalho/">Kaleb Carvalho</a>, <a href="https://www.linkedin.com/in/omatheusrsantos/">Matheus Ribeiro dos Santos</a>, <a href="https://www.linkedin.com/in/rafael-coutinho2004/">Rafael Coutinho</a> is licensed under <a href="http://creativecommons.org/licenses/by/4.0/?ref=chooser-v1" target="_blank" rel="license noopener noreferrer" style="display:inline-block;">Attribution 4.0 International</a>.</p>
