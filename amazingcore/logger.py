from enum import Enum
from datetime import datetime
import rich


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
        rich.print(f'[bold cyan][TRACE] {ts} > {message}[/]')
    elif LogLevel(log_level) == LogLevel.DEBUG:
        rich.print(f'[bold green][DEBUG] {ts} > {message}[/]')
    elif LogLevel(log_level) == LogLevel.INFO:
        rich.print(f'[bold white][INFO]  {ts} > {message}[/]')
    elif LogLevel(log_level) == LogLevel.WARN:
        rich.print(f'[bold yellow][WARN]  {ts} > {message}[/]')
    elif LogLevel(log_level) == LogLevel.ERROR:
        rich.print(f'[bold red][ERROR] {ts} > {message}[/]')
    elif LogLevel(log_level) == LogLevel.FATAL:
        rich.print(f'[bold white on red][FATAL] > {ts} {message}[/]')


# log('trace text', LogLevel.TRACE)
# log('debug text', LogLevel.DEBUG)
# log('debug text', LogLevel.INFO)
# log('debug text', LogLevel.WARN)
# log('debug text', LogLevel.ERROR)
# log('debug text', LogLevel.FATAL)
