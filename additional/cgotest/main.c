#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "cgotest.h"
const char* hello = "abcdef";
int main() {
    GoInt a = 1;
    GoString b;
    b.p = hello;
    b.n = strlen(b.p);
    int i = sizeof(b.p);
    for (; i >= 0 ; i--) {
        printf("%d\n", b.p[i]);
    }
    GoInt64 result = MyFunc(a,b);
    printf("%lld\n", result);
    printf("%td\n", b.n);
}

// ldd for linux 
// otool -L cgotest.bin
//// go build -buildmode=c-shared -o cgotest.so cexport.go
//// go build -buildmode=c-archive -o cgotest.a cexport.go
// rm cgotest.bin; gcc -c main.c -o hello.o && gcc hello.o cgotest.so -o cgotest.bin && ./cgotest.bin
// rm cgotest.bin; gcc -c main.c -o hello.o && gcc hello.o cgotest.a -o cgotest.bin && ./cgotest.bin