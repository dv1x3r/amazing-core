from enum import Enum
from datetime import datetime


class LogLevel(Enum):
    TRACE = 1
    DEBUG = 2
    INFO = 3
    WARN = 4
    ERROR = 5
    FATAL = 6


class Color:
    # https://stackoverflow.com/questions/287871/how-to-print-colored-text-to-the-terminal
    RED = '\x1b[1;31;1m'
    GREEN = '\x1b[1;32;1m'
    YELLOW = '\x1b[1;33;1m'
    BLUE = '\x1b[1;34;1m'
    CYAN = '\x1b[1;36;1m'
    FATAL = '\x1b[1;37;41m'
    END = '\x1b[0m'


def log(message: str, log_level: LogLevel = 3):
    ts = datetime.now().strftime('%d/%b/%y %H:%M:%S')
    if LogLevel(log_level) == LogLevel.TRACE:
        print(f'{Color.BLUE}[TRACE] {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.DEBUG:
        print(f'{Color.GREEN}[DEBUG] {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.INFO:
        print(f'[INFO]  {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.WARN:
        print(f'{Color.YELLOW}[WARN]  {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.ERROR:
        print(f'{Color.RED}[ERROR] {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.FATAL:
        print(f'{Color.FATAL}[FATAL] > {ts} {message}', end='')
    print(Color.END)


# log('trace text', LogLevel.TRACE)
# log('debug text', LogLevel.DEBUG)
# log('debug text', LogLevel.INFO)
# log('debug text', LogLevel.WARN)
# log('debug text', LogLevel.ERROR)
# log('debug text', LogLevel.FATAL)
