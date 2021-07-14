# jco (Jonathan's Converter), for working with hex and binary numbers

Command line utility to quickly get info about the hex, binary and decimal representations of numbers, and the most common operations on the bits. This is particularly helpful when looking
at hex dumps from some source when you're not sure which bits go where.

## Installation

* Download the latest release for your OS/arch from https://github.com/jonathangjertsen/jco-go/releases/latest
* Unpack the tar.gz to get the `jco` executable
* Move it to some directory that is on the PATH
* Test it by running `jco -v`, it should print the version of the latest release.

### Usage

Here is the output of `jco -h`:

```
jco (Jonathan's converter) v1.0.1

Usage:

        Show information about <number>
                jco <number>

        Show information about how <number1> and <number2> relate
                jco <number1> <number2>

        Like the above, but treat numbers as 16-bit
                jco <number1> <number2> -b 16

        Show this help screen
                jco --help

        Show one-liner version
                jco --version

Below is a list of the operations when running jco <number>:

        twos_complement:        Two's complement (depends on bit width)
        popcount:               Number of bits that are 1
        clz:                    Number of leading zeros
        nbits:                  Number of bits needed to represent the number
        reverse_bitorder        Reverses the bit order within each byte    (0b11100011 -> 0b11000111)
        reverse_nibbleorder     Reverses the nibble order within each byte (0xab -> 0xba)
        reverse_byteorder       Reverses the byte order
        reverse_bitstring       Interprets the input as a stream of bits, and reverses them.
                                Equivalent to reverse_bitorder followed by reverse_byteorder.
```

You'll want to supply either 1 or 2 numbers. Here are some examples:

`jco 1877 0x4a5e`

```
          FORMULA   |      DECIMAL   HEXADECIMAL                               BINARY
             1877   |         1877    0x00000755   0b00000000000000000000011101010101
           0x4a5e   |        19038    0x00004a5e   0b00000000000000000100101001011110
   1877  + 0x4a5e   |        20915    0x000051b3   0b00000000000000000101000110110011
   1877  | 0x4a5e   |        20319    0x00004f5f   0b00000000000000000100111101011111
   1877  & 0x4a5e   |          596    0x00000254   0b00000000000000000000001001010100
   1877  ^ 0x4a5e   |        19723    0x00004d0b   0b00000000000000000100110100001011
   1877 ^~ 0x4a5e   |   4294947572    0xffffb2f4   0b11111111111111111011001011110100
   1877  - 0x4a5e   |   4294950135    0xffffbcf7   0b11111111111111111011110011110111
   1877 &~ 0x4a5e   |         1281    0x00000501   0b00000000000000000000010100000001
   1877 >> 0x4a5e   |            0    0x00000000   0b00000000000000000000000000000000
   1877 << 0x4a5e   |            0    0x00000000   0b00000000000000000000000000000000
   0x4a5e  - 1877   |        17161    0x00004309   0b00000000000000000100001100001001
   0x4a5e &~ 1877   |        18442    0x0000480a   0b00000000000000000100100000001010
   0x4a5e >> 1877   |            0    0x00000000   0b00000000000000000000000000000000
   0x4a5e << 1877   |            0    0x00000000   0b00000000000000000000000000000000
```


`jco 1877 0x4a5e -b 16`

```
          FORMULA   |   DECIMAL   HEXADECIMAL               BINARY
             1877   |      1877        0x0755   0b0000011101010101
           0x4a5e   |     19038        0x4a5e   0b0100101001011110
   1877  + 0x4a5e   |     20915        0x51b3   0b0101000110110011
   1877  | 0x4a5e   |     20319        0x4f5f   0b0100111101011111
   1877  & 0x4a5e   |       596        0x0254   0b0000001001010100
   1877  ^ 0x4a5e   |     19723        0x4d0b   0b0100110100001011
   1877 ^~ 0x4a5e   |     45812        0xb2f4   0b1011001011110100
   1877  - 0x4a5e   |     48375        0xbcf7   0b1011110011110111
   1877 &~ 0x4a5e   |      1281        0x0501   0b0000010100000001
   1877 >> 0x4a5e   |         0        0x0000   0b0000000000000000
   1877 << 0x4a5e   |         0        0x0000   0b0000000000000000
   0x4a5e  - 1877   |     17161        0x4309   0b0100001100001001
   0x4a5e &~ 1877   |     18442        0x480a   0b0100100000001010
   0x4a5e >> 1877   |         0        0x0000   0b0000000000000000
   0x4a5e << 1877   |         0        0x0000   0b0000000000000000
```

`jco 0x1877`

```
                       FORMULA   |      DECIMAL   HEXADECIMAL                               BINARY
                       0x1877    |         6263    0x00001877   0b00000000000000000001100001110111
                      ~0x1877    |   4294961032    0xffffe788   0b11111111111111111110011110001000
       twos_complement(0x1877)   |   4294961033    0xffffe789   0b11111111111111111110011110001001
              popcount(0x1877)   |            8    0x00000008   0b00000000000000000000000000001000
                   clz(0x1877)   |           19    0x00000013   0b00000000000000000000000000010011
                 nbits(0x1877)   |           13    0x0000000d   0b00000000000000000000000000001101
     reverse_bitstring(0x1877)   |   3994550272    0xee180000   0b11101110000110000000000000000000
      reverse_bitorder(0x1877)   |         6382    0x000018ee   0b00000000000000000001100011101110
     reverse_byteorder(0x1877)   |   1998061568    0x77180000   0b01110111000110000000000000000000
   reverse_nibbleorder(0x1877)   |        33143    0x00008177   0b00000000000000001000000101110111
```

`jco 1877`

```
                     FORMULA   |      DECIMAL   HEXADECIMAL                               BINARY
                       1877    |         1877    0x00000755   0b00000000000000000000011101010101
                      ~1877    |   4294965418    0xfffff8aa   0b11111111111111111111100010101010
       twos_complement(1877)   |   4294965419    0xfffff8ab   0b11111111111111111111100010101011
              popcount(1877)   |            7    0x00000007   0b00000000000000000000000000000111
                   clz(1877)   |           21    0x00000015   0b00000000000000000000000000010101
                 nbits(1877)   |           11    0x0000000b   0b00000000000000000000000000001011
     reverse_bitstring(1877)   |   2866806784    0xaae00000   0b10101010111000000000000000000000
      reverse_bitorder(1877)   |        57514    0x0000e0aa   0b00000000000000001110000010101010
     reverse_byteorder(1877)   |   1426522112    0x55070000   0b01010101000001110000000000000000
   reverse_nibbleorder(1877)   |        28757    0x00007055   0b00000000000000000111000001010101
```

That's all it does!
