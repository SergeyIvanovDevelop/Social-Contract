from typing import Any
from json import load
from pathlib import Path

CONFIG_FILENAME = 'Backend/PostgreSQL/config.json'
CONFIG_FILE_PATH = Path(__file__).parents[1].joinpath(CONFIG_FILENAME)

class Config:
    def __init__(self) -> None:
        with open(CONFIG_FILE_PATH, 'r') as config_file:
            self.data: dict = load(config_file)

    def get(*keys: Any) -> Any:
        value = None
        if keys:
            config: dict = Config().data.get(keys[0], None)
            for key in keys[1:]:
                if config is not None:
                    config: dict = config.get(key, None)
            value = config
        return value

if __name__ == "__main__":
    pass
