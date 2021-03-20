from enum import Enum
from colorama import Fore, Back, Style
from datetime import datetime


class LogLevel(Enum):
    TRACE = 1
    DEBUG = 2
    INFO = 3
    WARN = 4
    ERROR = 5
    FATAL = 6


def log(message: str, log_level: LogLevel = 3):
    ts = datetime.now().strftime('%d/%b/%y %H:%M:%S')
    if LogLevel(log_level) == LogLevel.TRACE:
        print(f'{Fore.CYAN}[TRACE] {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.DEBUG:
        print(f'{Fore.GREEN}[DEBUG] {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.INFO:
        print(f'{Fore.WHITE}[INFO]  {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.WARN:
        print(f'{Fore.YELLOW}[WARN]  {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.ERROR:
        print(f'{Fore.RED}[ERROR] {ts} > {message}', end='')
    elif LogLevel(log_level) == LogLevel.FATAL:
        print(f'{Fore.BLACK}{Back.RED}[FATAL] > {ts} {message}', end='')
    print(Style.RESET_ALL)


# log('trace text', LogLevel.TRACE)
# log('debug text', LogLevel.DEBUG)
# log('debug text', LogLevel.INFO)
# log('debug text', LogLevel.WARN)
# log('debug text', LogLevel.ERROR)
# log('debug text', LogLevel.FATAL)
