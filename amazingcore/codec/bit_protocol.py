from asyncio.streams import StreamReader, StreamWriter


class BitProtocol:

    async def decode_data_length(self, reader: StreamReader):
        data_length = 0
        while True:
            data_byte = await reader.read(1)
            if not data_byte or data_byte == 0:  # client is disconnected
                raise ConnectionError('empty stream')
            byte_int = int.from_bytes(data_byte, 'big')
            data_length <<= 7  # shift the big endian
            data_length |= 0x7F & byte_int  # 0x80 if extra bytes
            if data_length > 0xfffffff:
                raise ValueError('message length is out of range')
            if 0x80 & byte_int == 0:
                break  # finished reading all length bytes
        return data_length

    async def read_data(self, reader: StreamReader, data_length):
        data = await reader.read(data_length)  # message starts with its size
        if data[-1] != 0:  # and ends with 0 byte
            raise ValueError(f'last message byte is not 0: {data}')
        return bytearray(data)

    def __encode_data_length__(self, data: bytearray):
        data_length = len(data) + 1  # + zero byte
        length_bytes = bytearray()
        if data_length > 0xfffffff:
            raise OverflowError('message length is out of range', data)
        if data_length > 0x1fffff:
            length_bytes.append(0x80 | (0x7F & (data_length >> 21)))
        if data_length > 0x3fff:
            length_bytes.append(0x80 | (0x7F & (data_length >> 14)))
        if data_length > 0x7F:
            length_bytes.append(0x80 | (0x7F & (data_length >> 7)))
        length_bytes.append(0x7F & data_length)
        return length_bytes

    async def write_message(self, writer: StreamWriter, data: bytearray):
        send_array = self.__encode_data_length__(data)
        send_array.extend(data)
        send_array.append(0)
        writer.write(send_array)
        await writer.drain()
