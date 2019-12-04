#include <stdio.h>
#include <stdlib.h>

// typedef struct {
//     unsigned  Field1: 1;
//     unsigned  Field2: 1;
//     unsigned  Field3: 1;
//     unsigned  : 1;
//     unsigned  Field4: 1;
//     unsigned  Field5: 3;
//     unsigned  Field6: 24;
// } someField;


// typedef struct {
//     unsigned light : 1;
//     unsigned fan: 1;
//     int count;            /* 4 bytes */
//     unsigned ac : 4;
//     unsigned : 4;
//     unsigned clock : 1;
//     unsigned : 0;
//     unsigned flag : 1;
// } someField;

// someField newBitmap() {
//     someField bitmap;
//     return bitmap;
// }

// someField setField1(someField m) {
//     m.light = 1;
//     return m;
// }

// someField clearField1(someField m) {
//     m.light = 0;
//     return m;
// }


// someField setField2(someField m) {
//     m.flag = 1;
//     return m;
// }

// someField clearField2(someField m) {
//     m.flag = 0;
//     return m;
// }

// void printsize() {
//     // someField v ;
//     // printf("%d\n", sizeof(v));
//     printf("%ld\n", sizeof(someField));
// }

typedef unsigned bool;
bool true = 1;
bool false = 0;

typedef struct {
    bool light  : 1;
    bool fan    : 1;
    bool count;            /* 4 bytes */
    bool ac     : 4;
    bool        : 4;
    bool clock  : 1;
    bool        : 0;
    bool flag   : 1;
} someField;

someField newBitmap() {
    someField bitmap;
    return bitmap;
}

someField setField1(someField m) {
    m.light = true;
    return m;
}

someField clearField1(someField m) {
    m.light = false;
    return m;
}

bool getField1(someField m) {
    return m.light;
}


someField setField2(someField m) {
    m.flag = true;
    return m;
}

someField clearField2(someField m) {
    m.flag = false;
    return m;
}

bool getField2(someField m) {
    return m.flag;
}

void printsize() {
    // someField v ;
    // printf("%d\n", sizeof(v));
    printf("%ld\n", sizeof(someField));
}