# jco (Jonathan's Converter), for working with hex and binary numbers

Command line utility to quickly get info about the hex, binary and decimal representations of numbers, and the most common operations on the bits. This is particularly helpful when looking
at hex dumps from some source when you're not sure which bits go where.

## Installation and usage

Installation: TODO

### Usage

You'll want to supply either 1 or 2 numbers. Here are some examples:

`jco 0x1877`

```
                  FORMULA  |  DECIMAL  HEXADECIMAL                              BINARY
                  0x1877   |     6263   0x00001877  0b00000000000000000001100001110111
                 ~0x1877   |    59272   0x0000e788  0b00000000000000001110011110001000
  twos_complement(0x1877)  |    59273   0x0000e789  0b00000000000000001110011110001001
         popcount(0x1877)  |        8   0x00000008  0b00000000000000000000000000001000
              clz(0x1877)  |        3   0x00000003  0b00000000000000000000000000000011
            nbits(0x1877)  |       13   0x0000000d  0b00000000000000000000000000001101
```

`jco 1877`

```
                FORMULA  |  DECIMAL  HEXADECIMAL                              BINARY
                  1877   |     1877   0x00000755  0b00000000000000000000011101010101
                 ~1877   |    63658   0x0000f8aa  0b00000000000000001111100010101010
  twos_complement(1877)  |    63659   0x0000f8ab  0b00000000000000001111100010101011
         popcount(1877)  |        7   0x00000007  0b00000000000000000000000000000111
              clz(1877)  |        5   0x00000005  0b00000000000000000000000000000101
            nbits(1877)  |       11   0x0000000b  0b00000000000000000000000000001011
```

`jco 1877 -b 16`

```
                FORMULA  |  DECIMAL  HEXADECIMAL              BINARY
                  1877   |     1877       0x0755  0b0000011101010101
                 ~1877   |    63658       0xf8aa  0b1111100010101010
  twos_complement(1877)  |    63659       0xf8ab  0b1111100010101011
         popcount(1877)  |        7       0x0007  0b0000000000000111
              clz(1877)  |        5       0x0005  0b0000000000000101
            nbits(1877)  |       11       0x000b  0b0000000000001011
```

`jco 1877 0x4a5e -b 16`

```
         FORMULA  |  DECIMAL  HEXADECIMAL              BINARY
            1877  |     1877       0x0755  0b0000011101010101
          0x4a5e  |    19038       0x4a5e  0b0100101001011110
  1877  + 0x4a5e  |    20915       0x51b3  0b0101000110110011
  1877  | 0x4a5e  |    20319       0x4f5f  0b0100111101011111
  1877  & 0x4a5e  |      596       0x0254  0b0000001001010100
  1877  ^ 0x4a5e  |    19723       0x4d0b  0b0100110100001011
  1877 ^~ 0x4a5e  |    45812       0xb2f4  0b1011001011110100
  1877  - 0x4a5e  |    48375       0xbcf7  0b1011110011110111
  1877 &~ 0x4a5e  |     1281       0x0501  0b0000010100000001
  0x4a5e  - 1877  |    17161       0x4309  0b0100001100001001
  0x4a5e &~ 1877  |    18442       0x480a  0b0100100000001010
```

That's all it does!
