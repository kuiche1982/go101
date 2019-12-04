00 00 49 // length 73
01 04 // header fragment, has all headers in this package
00 00 00 01 
// pkg length from here
83 // padding length 131 = 16 byte 3 bit
86 45 92 62
91 // weight
9a a5 f1 96 
8a 0f b8 b6 77
31 0a c6 32 d1 
41 ff 41 8a 08 
9d 5c 0b 81 70 
dc 78 0f 3b 5f 
8b 1d 75 d0 62 
0d 26 3d 4c 4d 
65 64 7a 8d 9a 
ca c8 b4 c7 60 
2b 89 95 c0 b4 // 95 should be seperated 90, and remove other after 95
85 ef 40 02 74 
65 86 4d 83 35 
05 b1 1f

00 00 0b // length 11
00  // fragment type: Data 0x00
01  // this fragment has all data 0x01
00 00 00 01 
// pkg length from here
00  // padding length
00 // is compressed
00 00 06 // length
0a 04 6b 6b 6b 6b  // serialized obj


0a xxxxxyyy  Field order + Field Type
here it's a reference type, which will is length prefixed. ( default to string)



