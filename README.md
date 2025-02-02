# kusmala  
**masih dalam pengembangan*  

id:  
Sebuah bahasa pemrograman interpretasi dengan penegasan pada Bahasa Indonesia. Ditulis dalam Go.  


en:  
A interpreted programming language with emphasis on Indonesian language. Written in Go.  

## Pemasangan  
### Build from source  
Clone repo ini:
```
$ git clone https://github.com/Vricap/kusmala.git  
$ cd kusmala
```  

Build dengan *make*:  
```
$ make
```  
Pastikan dalam komputer mu sudah terpasang [make](https://www.gnu.org/software/make/) dan [Go](https://go.dev/). Binary akan berada di `./bin/kusmala`  

Alternatif tanpa *make*:  
```
$ go build -o ./bin/ .
```  
Lakukan ini jika tidak mempunyai *make*.  

## Penggunaan  
Menjalankan kusmala tanpa argumen akan membuka REPL:  
![screenshot 1](./resource/screenshot/1.png)  

### Contoh Kode  
```
buat a = 1; // pemberian nilai ke variabel
buat b = 2;

buat f = fungsi (x, y) { // fungsi literal
	jika (x > y) { // kondisional 
		kembalikan x;
	} lainnya {
		kembalikan y;
	}
}

buat c = f(a, b);
cetak("Nilai terbesar ialah:", c); // cetak akan mengeluarkan hasil ke stdout
```  

Menjalankan kode dari file:  
```
$ ./bin/kusmala ./example/fungsi_dan_jika.km  
Nilai terbesar ialah: 2
```  
Tempatkan lokasi file pada argumen ke dua.  

Gunakan argumen `-tree` untuk mencetak pohon AST dari kode:  
```
$ ./bin/kusmala ./example/fungsi_dan_jika.km -tree  
AST_TREE:
  BUAT_STATEMENT:
    IDENT: a
    INTEGER_LITERAL: 1

  BUAT_STATEMENT:
    IDENT: b
    INTEGER_LITERAL: 2

  BUAT_STATEMENT:
    IDENT: f
    FUNGSI_EXPRESSION: 
      PARAMS: 
        IDENT: x
        IDENT: y
      FUNGSI_BODY: 
        JIKA_STATEMENT:
          CONDITION:
            INFIX_EXPRESSION:
              IDENT: x
              OEPERATOR: >
              IDENT: y
            JIKA_BLOCK: 
              KEMBALIKAN_STATEMENT:
                IDENT: x
            LAINNYA_BLOCK: 
              KEMBALIKAN_STATEMENT:
                IDENT: y

  BUAT_STATEMENT:
    IDENT: c
    CALL_EXPRESSION: 
      IDENT: f
      ARGUMENTS: 
        IDENT: a
        IDENT: b

  CETAK_STATEMENT: 
    EXPRESSION: 
      INFIX_EXPRESSION:
        INFIX_EXPRESSION:
          INTEGER_LITERAL: 28
          OEPERATOR: /
          INTEGER_LITERAL: 7
        OEPERATOR: *
        INTEGER_LITERAL: 2
      INFIX_EXPRESSION:
        INTEGER_LITERAL: 15
        OEPERATOR: *
        INTEGER_LITERAL: 3
```  
