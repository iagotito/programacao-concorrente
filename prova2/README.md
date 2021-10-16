UFCG/UASC Programação concorrente 2020.2e - Prova 2  
Professor: Thiago Emmanuel Pereira

# Q1 - Go fsck yourself (10.0)

O programa fsck é um clássico de sistemas UNIX usado para verificar a
consistência de sistemas de arquivos. Considere que você precisará
reimplementar uma nova versão concorrente deste programa que verifica a
consistência de um sistema de arquivo com base em um caminho para um
diretório passado como argumento:

```
fsck /home/thiagoepdc
```

No exemplo acima, a verificação será feita para toda a sub-árvore do
sistema de arquivo tendo o diretório /home/thiagoepdc como raiz. Leve
em conta que a execução do programa demora bastante e o usuário deseja
que seja reportado o andamento da execução parcial do programa. Isso
deve ser feito escrevendo o progresso na saída padrão, a cada um
segundo, até a finalização do programa:

$  go run fsck.go /home/thiagoepdc
$ fscked\_files 10 damaged\_files 3 fscked\_dirs 5 damaged\_dirs 1
$ fscked\_files 15 damaged\_files 4 fscked\_dirs 6 damaged\_dirs 1
$ fscked\_files 27 damaged\_files 7 fscked\_dirs 7 damaged\_dirs 2
$ fscked\_files 30 damaged\_files 7 fscked\_dirs 8 damaged\_dirs 2


Considere a seguinte API:

//retorna um vetor de Files representando arquivos e diretórios  
//contidos no diretório de nome dirname  
io.ReadDir(dirname string) []File

//indica se o File é um diretório (considere que se um File não é um  
//diretório é um arquivo)  
io.IsDir(file File) bool

//retorna o nome de um File  
io.Path(file File) string

//retorna o nome completo (partindo da raiz do sistema de arquivos)  
//de um File  
io.AbsolutePath(file File)

//verifica se um arquivo está danificado.  
io.fsckFile(file File) bool

//verifica se um diretório está danificado.  
io.fsckDir(file File) bool

//retorna o Diretório pai de um File  
io.parent(file File) File

Leve em consideração as seguintes diretivas:  
As funções io.fsckDir e io.fsckFile devem ser executadas por goroutines
(diferentes da main goroutine);  
O número máximo goroutines executando io.fsckDir concorrentemente é 8;  
O número máximo goroutines executando io.fsckFile concorrentemente é 8;  
io.fsckDir não deve ser executada para todos os diretório de uma
sub-árvore. Ao invés disso, deve ser executada somente caso um dos
"filhos" diretos do diretório ter sido detectado como danificado (seja
o filho arquivo ou diretório);  
io.fsckFile será executa para todos os arquivos da árvore;  
No report de progresso, as strings fscked\_files ou fscked\_dirs
indicam a quantidade de arquivos ou diretórios para os quais chamamos
as funções io.fsckDir e io.fsckFile.

Considere também que neste sistema de arquivos não há links.

obs.1 Crie qualquer função utilitária que achar necessária.  
obs.2 Crie qualquer estrutura de dados que achar necessária.  
obs.3 É pouco provável que você precise de uma construção no estilo de
shared memory. Canais devem ser suficientes. Tente usar somente canais.
Você será penalizado caso tenha usado shared memory constructs
(semáforos, var. cond) em situações nas quais canais seriam mais
adequados.  
obs.4 Corretude é o mais importante. Entretanto, código complicado em
excesso será penalizado.  
obs.5 Você pode usar pseudo-código ou programar direto em golang (eu
recomendo programar na linguagem para ter o auxílio do compilador).  
obs.6 Crie uma função main que trata os argumentos (o root path),
inicializa os objetos bem como cria e chama as funções necessárias.  
obs.7 Note que as funções io.fsckDir e io.fsckFile nao precisam ser
chamadas diretamente como goroutine, no estilo go io.fsckDir. Ao invés
disso, você também pode criar uma função anônima, encapsular a chamada
para io.fsckDir e io.fsckFile através dessa função anônima e chamá-la
com a diretiva go  
obs.8 Erros na contagem de diretórios (por duplicação de execução de
io.fsckDir) implacarão em um pênalty baixo (não maior que 1,0)  

sugestão: Implemente primeiro uma versão serial do programa. Talvez,
essa versão inicial nem precise fazer o report periodicamente.

## Anexo

Abaixo, equivalente em golang da API listada acima

https://pkg.go.dev/io/ioutil#ReadDir 
```go
func ReadDir(dirname string) ([]fs.FileInfo, error)
```

https://pkg.go.dev/io/fs#FileInfo
```go
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}
```

A função abaixo pode ser usada para montar o absolutepath
https://pkg.go.dev/path/filepath#Join 
```go
func Join(elem ...string) string
```

O esqueleto de funcão abaixo pode ser usado para implementar fake das funções de fsck (caso usem golang na prova)

```go
import (
    "math/rand"
    "time"
)

func fsckFile(path string) bool {
//dorme por um tempo aleatório entre 0 e 3 segs. aumente se //preferir
    rSleep := rand.Intn(4)
    time.Sleep(time.Duration(rSleep) * time.Second)
//retorna true ou false com igual probabilidade
    rn := rand.Intn(2)
    return (rn % 2 == 0)
}
```
