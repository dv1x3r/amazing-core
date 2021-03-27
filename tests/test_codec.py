import unittest
from unittest.mock import AsyncMock
from amazingcore.codec.bit_protocol import BitProtocol
from amazingcore.codec.bit_stream import BitStream

# flags = 0, service_class = 18, message_type = 566, log_correlator = ''
# client_name = 'AmazingWorld'
AMAZING_WORLD = b'\x15 \xc2\\\x04m\x0c\x0c\x18AmazingWorld\x00'


class TestBitProtocol(unittest.IsolatedAsyncioTestCase):

    async def test_decode_length(self):
        cases = [
            {'stream': [b'\x00'], 'expected': 0},
            {'stream': [b'\x15'], 'expected': 21},
            {'stream': [b'\x81', b'\x00'], 'expected': 128},
            {'stream': [b'\x81', b'\x80', b'\x00'], 'expected': 16384},
            {'stream': [b'\x81', b'\x80', b'\x80', b'\x00'],
             'expected': 2097152}
        ]
        for case in cases:
            reader = AsyncMock()
            reader.read.side_effect = case['stream']
            result = await BitProtocol().__decode_data_length__(reader)
            self.assertEqual(result, case['expected'])
        with self.assertRaises(ValueError):
            reader = AsyncMock()
            reader.read.side_effect = [
                b'\x81', b'\x80', b'\x80', b'\x80', b'\x80']
            result = await BitProtocol().__decode_data_length__(reader)

    def test_encode_length(self):
        cases = [
            {'bytes': bytes(0), 'expected': '01'},
            {'bytes': bytes(20), 'expected': '15'},
            {'bytes': bytes(128), 'expected': '81 01'},
            {'bytes': bytes(16384), 'expected': '81 80 01'},
            {'bytes': bytes(2097152), 'expected': '81 80 80 01'}]
        for case in cases:
            result = BitProtocol().__encode_data_length__(case['bytes'])
            self.assertEqual(result.hex(sep=' '), case['expected'])
        with self.assertRaises(OverflowError):
            BitProtocol().__encode_data_length__(bytes(268435455))


class TestBitStream(unittest.TestCase):

    def test_read_size(self):
        cases = [
            {'data': b'\x00', 'expected': 4},  # 0000
            {'data': b'\x80', 'expected': 8},  # 1000
            {'data': b'\xC0', 'expected': 16},  # 1100
            {'data': b'\xE0', 'expected': 24},  # 1110
            {'data': b'\xF0', 'expected': 32}]  # 1111
        for case in cases:
            bit_stream = BitStream(case['data'])
            result = bit_stream.__read_size__(4)
            self.assertEqual(result, case['expected'])
        with self.assertRaises(ValueError):
            bit_stream = BitStream(b'\xF8')  # 1111 1
            result = bit_stream.__read_size__(4)

    def test_read_int(self):
        cases = [  # is_integer size data        i s data
            {'data': b'\x80', 'expected': 0},  # 1 0 0000
            {'data': b'\x90', 'expected': 4},  # 1 0 0100
            {'data': b'\x9C', 'expected': 7},  # 1 0 0111
            {'data': b'\xC1\x00', 'expected': 8},  # 1 10 0000 1000
            {'data': b'\xC8\x00', 'expected': 64},  # 1 10 0100 0000
            {'data': b'\xCF\xE0', 'expected': 127},  # 1 10 0111 1111
            {'data': b'\xE0\x08\x00', 'expected': 128},  # 1 110 0000 0000 1000
            {'data': b'\xB0', 'expected': -4},  # 1 0 1100
            {'data': b'\xA0', 'expected': -8},  # 1 0 1000
            {'data': b'\xDE\xE0', 'expected': -9},  # 1 10 1111 0111
            {'data': b'\xD8\x00', 'expected': -64},  # 1 10 1100 0000
            {'data': b'\xD0\x00', 'expected': -128},  # 1 10 1000 0000
            {'data': b'\xEF\xF7\xF0', 'expected': -129}]  # 1 110 1111 1111 0111 1111
        for case in cases:
            bit_stream = BitStream(case['data'])
            result = bit_stream.read_int()
            self.assertEqual(result, case['expected'])

    def test_read_str(self):
        bit_stream = BitStream(AMAZING_WORLD[1:-1])
        bit_stream.cursor = 52  # client_name
        result = bit_stream.read_str()
        self.assertEqual(result, 'AmazingWorld')

    def test_write_size(self):
        cases = [
            {'int': 0, 'expected': b'\x00'},  # 4
            {'int': 4, 'expected': b'\x00'},  # 4
            {'int': 7, 'expected': b'\x00'},  # 4
            {'int': 8, 'expected': b'\x80'},  # 8
            {'int': 64, 'expected': b'\x80'},  # 8
            {'int': 127, 'expected': b'\x80'},  # 8
            {'int': 128, 'expected': b'\xC0'},  # 16
            {'int': -4, 'expected': b'\x00'},  # 4
            {'int': -8, 'expected': b'\x00'},  # 4
            {'int': -9, 'expected': b'\x80'},  # 8
            {'int': -64, 'expected': b'\x80'},  # 8
            {'int': -128, 'expected': b'\x80'},  # 8
            {'int': -129, 'expected': b'\xC0'}]  # 16
        for case in cases:
            bit_stream = BitStream()
            bit_stream.__write_size__(case['int'], 4)
            self.assertEqual(bit_stream.data, case['expected'], case['int'])
        with self.assertRaises(ValueError):
            BitStream().__write_size__(4294967296, 4)

    def test_write_int(self):
        cases = [  # is_integer size data        i s data
            {'expected': b'\x80', 'int': 0},  # 1 0 0000
            {'expected': b'\x90', 'int': 4},  # 1 0 0100
            {'expected': b'\x9C', 'int': 7},  # 1 0 0111
            {'expected': b'\xC1\x00', 'int': 8},  # 1 10 0000 1000
            {'expected': b'\xC8\x00', 'int': 64},  # 1 10 0100 0000
            {'expected': b'\xCF\xE0', 'int': 127},  # 1 10 0111 1111
            {'expected': b'\xE0\x08\x00', 'int': 128},  # 1 110 0000 0000 1000
            {'expected': b'\xB0', 'int': -4},  # 1 0 1100
            {'expected': b'\xA0', 'int': -8},  # 1 0 1000
            {'expected': b'\xDE\xE0', 'int': -9},  # 1 10 1111 0111
            {'expected': b'\xD8\x00', 'int': -64},  # 1 10 1100 0000
            {'expected': b'\xD0\x00', 'int': -128},  # 1 10 1000 0000
            {'expected': b'\xEF\xF7\xF0', 'int': -129}]  # 1 110 1111 1111 0111 1111
        for case in cases:
            bit_stream = BitStream()
            bit_stream.write_int(case['int'])
            self.assertEqual(bit_stream.data, case['expected'])

    def test_write_str(self):
        bit_stream = BitStream()
        bit_stream.write_str('ЯAmazing')
        self.assertEqual(bit_stream.data, b'\xC1\x20\xD0\xAFAmazing')


class TestBitStreamRW(unittest.TestCase):

    def test_read_write(self):
        int_values = [-32769, -32768, -128, -127, -8, -7,
                      0, 7, 8, 127, 128, 32768, 32769]
        bit_stream = BitStream()

        for i in int_values:
            bit_stream.write_int(i)
            bit_stream.write_long(i)
            bit_stream.write_str(str(i))
        bit_stream.cursor = 0
        for i in int_values:
            self.assertEqual(bit_stream.read_int(), i)
            self.assertEqual(bit_stream.read_long(), i)
            self.assertEqual(bit_stream.read_str(), str(i))


if __name__ == '__main__':
    unittest.main()
