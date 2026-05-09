# README-AI

Este documento descreve a estrutura interna do projeto arcpy2go para uso por humanos e agentes. O objetivo do projeto é fazer web parsing da documentação de ferramentas do ArcPy e gerar structs em Go que representam essas classes, com métodos para produzir a chamada em texto no formato do ArcPy.

## Visão geral

O fluxo principal do projeto é:

1. Receber uma URL de documentação de uma ferramenta ArcPy.
2. Fazer scraping do HTML da página.
3. Extrair nome da função, comentário, parâmetros, obrigatoriedade e enums.
4. Renderizar um arquivo Go a partir de um template.
5. Gerar structs que conseguem imprimir a chamada da ferramenta em texto.

Em termos práticos, o projeto transforma documentação web em código Go tipado, reduzindo o trabalho manual de escrever wrappers para ArcPy.

## Estrutura do repositório

### Arquivos na raiz

- [main.go](main.go): ponto de entrada da aplicação.
- [go.mod](go.mod): módulo Go e dependências.
- [README.md](README.md): documentação de uso do projeto.
- [README-AI.md](README-AI.md): documentação técnica interna do código.
- [LICENSE](LICENSE): licença do projeto.

### Arquivo principal

- [main.go](main.go): contém a CLI principal com a biblioteca padrão `flag`, lê argumentos, configura cache e coordena scraping e geração.

### Diretório arcpy-scraper/web

- [arcpy-scraper/web/scraper.go](arcpy-scraper/web/scraper.go): faz o scraping da documentação ArcPy com Colly e GoQuery.
- [arcpy-scraper/web/scraper_test.go](arcpy-scraper/web/scraper_test.go): valida a extração básica de dados de uma página real.

### Diretório gen

- [gen/gen.go](gen/gen.go): define as estruturas de dados usadas no gerador e executa o template.
- [gen/template.go.tmpl](gen/template.go.tmpl): template que vira código Go gerado.
- [gen/gen_test.go](gen/gen_test.go): valida a renderização do template.

### Diretório test

- [test/CreateTable.go](test/CreateTable.go): exemplo de saída já gerada para a ferramenta CreateTable.

## Responsabilidade de cada parte

### main.go

O arquivo [main.go](main.go) é o orquestrador do projeto.

Ele faz o seguinte:

- recebe 1 ou 2 argumentos na linha de comando;
- trata a primeira posição como URL da documentação ArcPy;
- cria um collector do Colly com cache local;
- chama o scraper em [arcpy-scraper/web/scraper.go](arcpy-scraper/web/scraper.go);
- decide se a saída vai para stdout ou para arquivo;
- aplica a opção `--package` para salvar o arquivo em uma pasta com o nome da função;
- aplica a opção `--package-name` para definir o nome do package gerado.

Também cria a pasta de cache em `%AppData%` ou no diretório de configuração do usuário, usando `os.UserConfigDir`.

### arcpy-scraper/web/scraper.go

O scraper lê a página da documentação ArcPy e monta uma estrutura do tipo [gen.Generator](gen/gen.go).

Ele extrai:

- o comentário/título da ferramenta;
- o nome da função a partir da assinatura;
- a string completa do comando ArcPy;
- a lista de parâmetros;
- os comentários de cada parâmetro;
- os valores enumerados, quando a documentação fornece opções fechadas;
- se o parâmetro é obrigatório ou opcional.

Os seletores usados são específicos da estrutura atual do site da documentação ArcGIS, então a extração depende do HTML manter o mesmo padrão.

### gen/gen.go

O pacote [gen/gen.go](gen/gen.go) define os tipos centrais do gerador:

- [gen.Generator](gen/gen.go): contém nome da função, comando, comentário e parâmetros.
- [gen.Parameter](gen/gen.go): representa um parâmetro de uma ferramenta ArcPy.
- [gen.Enum](gen/gen.go): representa uma opção enumerada de um parâmetro.

Esse pacote também faz a renderização do template embutido e fornece a função de conversão de snake_case para CamelCase usada na geração do código.

### gen/template.go.tmpl

O template gera uma struct Go com:

- campos para cada parâmetro extraído;
- tipos obrigatórios como string;
- parâmetros opcionais sem enum como ponteiros para string;
- parâmetros com enum como tipos nomeados gerados pelo próprio template;
- método `Name()` para retornar o comando ArcPy;
- método `Command()` para montar a chamada completa;
- método `Args()` para montar os argumentos no formato `param='value'`;
- método `String()` como atalho para `Command()`;
- método `CommentsString()` para juntar comentários adicionais.

O template também gera constantes para enums, permitindo usar valores tipados em vez de strings soltas.

## Fluxo de dados

```
URL da documentação ArcPy
        |
        v
main.go
        |
        v
arcpy-scraper/web/scraper.go
        |
        v
gen.Generator + Parameter + Enum
        |
        v
gen/template.go.tmpl
        |
        v
arquivo Go gerado ou stdout
```

### Detalhe do processamento

1. A CLI recebe a URL.
2. O Colly visita a página e aplica os seletores HTML.
3. O scraper monta a estrutura de geração.
4. O gerador calcula os nomes em CamelCase e aplica o template.
5. O resultado final é um arquivo Go com uma struct que representa a ferramenta ArcPy.

## Modelo de tipos gerados

O comportamento atual da geração é este:

- parâmetro obrigatório sem enum: campo do tipo string;
- parâmetro opcional sem enum: campo do tipo *string;
- parâmetro obrigatório com enum: tipo nomeado gerado;
- parâmetro opcional com enum: ponteiro para o tipo nomeado gerado.

No arquivo de exemplo [test/CreateTable.go](test/CreateTable.go), isso aparece claramente em campos como `OutPath`, `OutName`, `Template` e `OidType`.

## Exemplo do que é gerado

A saída gerada para CreateTable segue este padrão conceitual:

```go
type CreateTable struct {
    OutPath string
    OutName string
    Template *string
    OidType string
}

func (p CreateTable) Command() string {
    return "arcpy.management.CreateTable(...)"
}
```

Na prática, o método `Args()` monta a lista de argumentos em sequência, e `Command()` encapsula isso com o nome completo da função ArcPy.

## Testes existentes

### scraper_test.go

O teste em [arcpy-scraper/web/scraper_test.go](arcpy-scraper/web/scraper_test.go) valida se o scraper consegue ler a documentação real de CreateTable e extrair o comentário principal esperado.

### gen_test.go

O teste em [gen/gen_test.go](gen/gen_test.go) verifica se a renderização do template produz uma string não vazia a partir de uma estrutura de exemplo.

Esses testes confirmam o caminho principal do projeto: scraping funcional e renderização funcional.

## Dependências principais

As dependências centrais do projeto são:

- `github.com/gocolly/colly`: scraping HTTP e parsing orientado a seletores.
- `github.com/PuerkitoBio/goquery`: consultas ao DOM HTML dentro do scraper.
- A biblioteca padrão `flag`: parsing dos argumentos de linha de comando.

As demais dependências em [go.mod](go.mod) são transitivas ou de suporte.

## Limitações atuais

- Os seletores HTML são sensíveis a mudanças no site da documentação ArcGIS.
- O fluxo atual processa uma URL por execução.
- O modelo de tipos cobre bem strings, ponteiros e enums, mas ainda é específico para o formato da documentação atual.
- A geração assume que a assinatura da ferramenta pode ser reconstruída diretamente da página scrapeada.

## Resumo por pasta

- Raiz: ponto de entrada, módulo e documentação.
- main.go: CLI e fluxo de execução.
- arcpy-scraper/web: scraping e extração de dados da documentação.
- gen: estrutura de dados e geração de código.
- test: exemplo de saída gerada.

## Objetivo final do projeto

O propósito do arcpy2go é permitir que a documentação do ArcPy seja convertida em structs Go reutilizáveis, com saída textual compatível com a chamada da ferramenta original. Isso facilita criar código Go que emite comandos ArcPy sem precisar reescrever manualmente a estrutura de cada ferramenta.