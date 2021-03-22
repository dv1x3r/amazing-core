class BitStream:

    def __init__(self, data: bytearray = None):
        self.data = data or bytearray()
        self.cursor = 0

    def __byte_index__(self):
        return self.cursor >> 3

    def __bit_mask__(self):
        return 0x80 >> (self.cursor % 8)

    def __read_bit__(self):
        byte = self.__byte_index__()
        bit = self.__bit_mask__()
        self.cursor += 1
        return self.data[byte] & bit

    def __read_align_byte__(self):
        while self.__bit_mask__() != 0x80:
            self.cursor += 1  # go to the next byte start
        byte = self.__byte_index__()
        self.cursor += 8
        return self.data[byte]

    def __read_size__(self):
        size_bits = 4  # base size = 4
        while self.__read_bit__() != 0:
            size_bits <<= 1  # extra bits = 4, 8, 16, 32
            if size_bits > 32:
                raise ValueError('invalid integer size')
        return size_bits

    def read_start(self):
        if self.__read_bit__() != 0:  # message starts with 0
            raise ValueError('invalid message object')

    def read_int(self):
        if self.__read_bit__() == 0:  # integer starts with 1
            raise ValueError('invalid integer object')
        size_bits = self.__read_size__()  # is followed by number of bytes
        int_value = 0  # read bit by bit
        for _ in range(size_bits):
            int_value <<= 1  # shift the big endian
            int_value |= self.__read_bit__() != 0
        is_negative_mask = 1 << (size_bits - 1)
        if int_value & is_negative_mask != 0:  # the first bit stands for negative
            int_value |= -is_negative_mask  # signed two’s complement
        return int_value

    def read_str(self):
        size_bytes = self.read_int()
        if size_bytes == 0:  # string starts with size in bytes
            return None
        # characters are aligned on whole bytes
        bytes = [self.__read_align_byte__() for _ in range(size_bytes)]
        str_value = bytearray(bytes).decode('utf-8')
        return str_value

    def __write_bit__(self, active: int):
        byte = self.__byte_index__()
        bit = self.__bit_mask__()
        if bit == 0x80:  # cursor stands on the new byte
            self.data.append(0)
        if active != 0:  # set current bit
            self.data[byte] |= bit
        self.cursor += 1

    def __write_align_byte__(self, byte: int):
        while self.__bit_mask__() != 0x80:
            self.cursor += 1  # go to the next byte start
        self.data.append(byte)
        self.cursor += 8

    def __write_size__(self, int_value: int):
        if int_value is None:
            raise ValueError('int value is empty')
        size_bits = 4  # int_value must fit in size_bits
        if int_value > 0:
            int_max = 7
            while int_value > int_max:
                int_max = ((int_max + 1) << 4) - 1
                size_bits <<= 1  # 7, 127...
                if size_bits > 32:
                    raise ValueError('invalid pos integer size', int_value)
                self.__write_bit__(1)
        else:
            int_min = -8
            while int_value < int_min:
                int_min <<= size_bits
                size_bits <<= 1  # -8, -128...
                if size_bits > 32:
                    raise ValueError('invalid neg integer size', int_value)
                self.__write_bit__(1)
        self.__write_bit__(0)  # done with extra bits
        return size_bits

    def write_start(self):
        self.__write_bit__(0)  # Message starts with 0

    def write_int(self, int_value: int):
        self.__write_bit__(1)
        size_bits = self.__write_size__(int_value)
        write_bit = (1 << (size_bits - 1))  # current write bit mask
        for _ in range(size_bits):
            bit = int_value & write_bit != 0
            self.__write_bit__(bit)
            write_bit >>= 1  # next write bit mask

    def write_str(self, str_value: str):
        if not str_value:
            self.write_int(0)
            return  # string starts with size (0 if empty)
        str_bytes = str_value.encode('utf-8')
        self.write_int(len(str_bytes))
        for char_byte in str_bytes:
            self.__write_align_byte__(char_byte)  # write to the byte start
