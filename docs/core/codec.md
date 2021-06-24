# Introduction to bit codec

Message communication between server and client happens using TCP sockets and streaming encoded messages.

## Message internals

All request messages consists of Header (**MessageHeader.cs**) and Base (**GSF\*Svc.cs** type) wrapped into the **GSFRequestMessage.cs** (**GSFMessage.cs**) class. In our example **GSFGetClientVersionInfoSvc** is used.

- [MessageHeader properties](../messages/message_header.md)
- [GSFGetClientVersionInfoSvc Message properties](../messages/get_client_version_info.md)

Serialization starts from the parent to the child class, and Header attributes are always the first.

## Encoding

```cs
// Converting message into bytes, writing message size, message and 0 byte
GSFSession.cs -> class GSFSession -> WriteMessage(); 

// Converting message into bytes array
GSFBitProtocolCodec.cs -> class GSFBitProtocolCodec -> WriteMessage();

// Determining what and how data type is going to be written
GSFBitProtocolCodec.cs -> class BitOutput -> Write();

// Responsible for writing bytes
BitStream.cs -> class BitStream -> Put*();

// Responsible for writing compressed integers
BitStream.cs -> class IntModeX -> Put();
```

- All the GSFIExternalizable objects do start with bit, showing is it null or not.
  - GSFBitProtocolCodec.cs -> class BitOutput -> Write(GSFIExternalizable o)
- Integers are compressed by default, meaning they do start with number of bits required.
  - GSFBitProtocolCodec.cs -> class BitOutput -> Write(int i) -> class IntModeX
  - [0\] if only following 4 bits are needed (5 bits per integer), when i < 8 and i > -9
  - [10\] if 8 bits are needed (10 bits per integer) etc.
- Strings do start with number of required bytes and aligned bytes with characters.
  -  GSFBitProtocolCodec.cs -> class BitOutput -> Write(string s)

## "AmazingWorld" Encoding Process

```js
const net = require('net');
const server = net.createServer();
server.listen(8182, '127.0.0.1');
server.on('connection', socket => {
  socket.on('data', data => {
    console.log(`str: ${data}`);
    console.log('hex:', data, data.byteLength, 'bytes', '\n');
  });
});
```

- The first message from the client is 22 bytes: ```§ �\♦m♀♀↑AmazingWorld```
- hex: <15 20 c2 5c 04 6d 0c 0c 18 41 6d 61 7a 69 6e 67 57 6f 72 6c 64 00>

```py
def print_i(i_arr):
    print('byte  int    hex           bin       char')
    for bit, i in enumerate(i_arr):
        int_str = str(i).ljust(3)
        hex_str = hex(i)[2:].zfill(2)
        bin_str = bin(i)[2:].zfill(8)
        bin_str = str(bit * 8).rjust(3) + ': ' + bin_str[:4] + ' ' + bin_str[4:]
        chr_str = chr(i) if i <= 127 else ''
        bit = str(bit + 1).ljust(2)
        print(bit, int_str, hex_str, bin_str, chr_str, sep='    ')

src_hex = "15 20 c2 5c 04 6d 0c 0c 18 41 6d 61 7a 69 6e 67 57 6f 72 6c 64 00".split()
src_int = list(map(lambda i: int(i, 16), src_hex))
print_i(src_int)
```

```
byte  int    hex           bin       char
1     21     15      0: 0001 0101    
2     32     20      8: 0010 0000     
3     194    c2     16: 1100 0010    
4     92     5c     24: 0101 1100    \
5     4      04     32: 0000 0100    
6     109    6d     40: 0110 1101    m
7     12     0c     48: 0000 1100    
8     12     0c     56: 0000 1100    
9     24     18     64: 0001 1000    
10    65     41     72: 0100 0001    A
11    109    6d     80: 0110 1101    m
12    97     61     88: 0110 0001    a
13    122    7a     96: 0111 1010    z
14    105    69    104: 0110 1001    i
15    110    6e    112: 0110 1110    n
16    103    67    120: 0110 0111    g
17    87     57    128: 0101 0111    W
18    111    6f    136: 0110 1111    o
19    114    72    144: 0111 0010    r
20    108    6c    152: 0110 1100    l
21    100    64    160: 0110 0100    d
22    0      00    168: 0000 0000    
```

The first byte contains message size: 21 bytes  
Last byte contains nothing: 0, meaning this is the end of the message.

- 8: **Message is not null**: bitStream.Put(GSFRequestMessage == null) => False (0)
- 9: **MessageHeader is not null**: bitStream.Put(MessageHeader == null) => False (0)

MessageHeader -> flags = 0

- 10: **Is Compressed**: compressor.Put(this, MessageHeader.flags, 4) => bs.Put(i < w) => True (1)
- 11: **No extra bytes needed**: bs.Put(j < i) => False (0)
- 12-15: **flags** => [0000\]

MessageHeader -> svcClass = 18

- 16: **Is Compressed**: compressor.Put(this, MessageHeader.svcClass, 4) => bs.Put(i < w) => True (1)
- 17: **1 byte needed**: bs.Put(j < i) => True (1)
- 18: **No more extra bytes needed**: bs.Put(j < i) => False (0)
- 19-26: **svcClass** => [00010010\]

MessageHeader -> msgType = 566

- 27: **Is Compressed**: compressor.Put(this, MessageHeader.msgType, 4) => bs.Put(i < w) => True (1)
- 28: **1 byte needed**: bs.Put(j < i) => True (1)
- 29: **2 bytes needed**: bs.Put(j < i) => True (1)
- 30: **No more extra bytes needed**: bs.Put(j < i) => False (0)
- 31-46: **msgType**: [0000001000110110\]

MessageHeader -> requestId = 1

- 47: **Is Compressed**: compressor.Put(this, MessageHeader.requestId, 4) => bs.Put(i < w) => True (1)
- 48: **No extra bytes needed**: bs.Put(j < i) => False (0)
- 49-52: **requestId** => [0001\]

MessageHeader -> logCorrelator = ""

- 53: **Is Compressed**: compressor.Put(this, logCorrelator.UTF8.Length, 4) => bs.Put(i < w) => True (1)
- 54: **No extra bytes needed**: bs.Put(j < i) => False (0)
- 55-58: **length** => [0000\]

MessageBase

- 59: **MessageBase is not null**: bitStream.Put(GSFRequest == null) => False (0)

MessageBase -> clientName = "AmazingWorld" (12)

- 60: **Is Compressed**: compressor.Put(this, clientName.UTF8.Length, 4) => bs.Put(i < w) => True (1)
- 61: **1 byte needed**: bs.Put(j < i) => extra 8 bits => True (1)
- 62: **No more extra bytes needed**: bs.Put(j < i) => False (0)
- 63-70: **length** => [0000110\]
- 71: **align** => 0
- 72-167: **clientName** => "AmazingWorld"
