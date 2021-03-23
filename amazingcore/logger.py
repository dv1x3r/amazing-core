from enum import Enum
from rich.console import Console
from rich.panel import Panel
from rich.pretty import Pretty

console = Console(highlight=False, log_path=False, markup=True)


class LogLevel(Enum):
    DEBUG = 10
    INFO = 20
    WARN = 30
    ERROR = 40
    FATAL = 50


def log(log_level: LogLevel, message: str, debug_object: any = None):
    if LogLevel(log_level) == LogLevel.DEBUG:
        console.log(f'[bold green]DEBUG[/]  {message}')
        if debug_object:
            console.log(Panel(Pretty(debug_object), expand=False))
    elif LogLevel(log_level) == LogLevel.INFO:
        console.log(f'[bold blue]INFO[/]  {message}')
    elif LogLevel(log_level) == LogLevel.WARN:
        console.log(f'[bold yellow]WARN[/]  {message}')
    elif LogLevel(log_level) == LogLevel.ERROR:
        console.log(f'[bold red]ERROR[/] {message}')
        console.print_exception()
    elif LogLevel(log_level) == LogLevel.FATAL:
        console.log(f'[bold white on red]FATAL[/] {message}')
        console.print_exception()
