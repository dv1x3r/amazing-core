from enum import Enum
import datetime as dt
from rich.console import Console
from rich.panel import Panel
from rich.pretty import Pretty
from rich.theme import Theme

console = Console(theme=Theme(inherit=False))


class LogLevel(Enum):
    DEBUG = 10
    INFO = 20
    WARN = 30
    ERROR = 40
    FATAL = 50


def log(log_level: LogLevel, message: str, debug_object: any = None):
    ts = dt.datetime.now().strftime('%H:%M:%S')
    if LogLevel(log_level) == LogLevel.DEBUG:
        console.print(f'{ts} [bold green]DEBUG[/]  {message}')
        if debug_object:
            console.print(Panel(Pretty(debug_object), expand=False))
    elif LogLevel(log_level) == LogLevel.INFO:
        console.print(f'{ts} [bold]INFO[/]  {message}')
    elif LogLevel(log_level) == LogLevel.WARN:
        console.print(f'{ts} [bold yellow]WARN[/]  {message}')
    elif LogLevel(log_level) == LogLevel.ERROR:
        console.print(f'{ts} [bold red]ERROR[/] {message}')
        console.print_exception()
    elif LogLevel(log_level) == LogLevel.FATAL:
        console.print(f'{ts} [bold white on red]FATAL[/] {message}')
        console.print_exception()
