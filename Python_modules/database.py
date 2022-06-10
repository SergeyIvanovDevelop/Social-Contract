from psycopg2 import connect, OperationalError
from psycopg2.extras import execute_values
from psycopg2.sql import SQL, Identifier

from config import Config

class Database:
    def __init__(self) -> None:
        __conn: dict = Config.get('postgresql', 'connect')
        __dsn = ' '.join([
            f'dbname={__conn.get("database", "")}',
            f'user={dict(__conn.get("user", "")).get("name", "")}',
            f'password={dict(__conn.get("user", "")).get("password", "")}',
            f'host={__conn.get("host", "")}',
            f'port={__conn.get("port", "")}',
        ])
        try:
            self.__conn = connect(dsn=__dsn)
        except OperationalError as err:
            print(f'[ERROR]: Ошибка при подключении к БД - {err}')
            self.__conn = None
            self.__cur = None
        else:
            self.__cur = self.__conn.cursor()
        finally:
            del(__conn, __dsn)

    def __bool__(self) -> bool:
        return True if self.__conn and self.__cur else False

    def insert(self, table: str, fields: tuple, values: list) -> None:
        if self:
            query = SQL('INSERT INTO {table} ({fields}) VALUES %s;').format(
                table=Identifier(table),
                fields=SQL(', ').join([Identifier(field) for field in fields]),
            )
            try:
                execute_values(self.__cur, query, values)
                self.__conn.commit()
            except:
                print(f'[WARNING]: Ошибка при вставке данных. SQL-запрос:\n{query.as_string(self.__conn)}')

    def select(self, table: str, fields: tuple) -> list:
        if self:
            query = SQL('SELECT {fields} FROM {table};').format(
                table=Identifier(table),
                fields=SQL(',').join([Identifier(field) for field in fields]),
            )
            try:
                self.__cur.execute(query)
                return self.__cur.fetchall()
            except:
                print(f'[WARNING]: Ошибка при выборке данных. SQL-запрос:\n{self.__cur.query}')
        return []

    def execute(self, query):
        if self:
            query = SQL(query)
            try:
                self.__cur.execute(query)
                return self.__cur.fetchall()
            except:
                print(f'[WARNING]: Ошибка при выборке данных. SQL-запрос:\n{self.__cur.query}')
        return []

    def __del__(self):
        if self:
            self.__cur.close()
            self.__conn.close()

if __name__ == '__main__':
    db = Database()
    print(db.select('users', ('login', 'hash')))
