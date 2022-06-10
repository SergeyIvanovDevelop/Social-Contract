from datetime import datetime
from get_mark import parse_page_marks, get_mark
from database import Database as DB

db = DB()

def get_summ(child_login, timestamp) -> int:
    '''Получение суммы, заработанной за оценки'''
    diary_credential = db.execute(f'SELECT "diary_login", "diary_password" FROM "view_diary_credentials" WHERE "login" = \'{child_login}\';')[0]
    print('diary_credential = ',diary_credential)
    grades = list(parse_page_marks(get_mark(
        year=datetime.strptime(timestamp, '%Y-%m-%d').year,
        mounth=datetime.strptime(timestamp, '%Y-%m-%d').month,
        day=datetime.strptime(timestamp, '%Y-%m-%d').day,
        login=diary_credential[0],
        password=diary_credential[1]
    )).values())
    print('grades = ',grades)
    grades_list = []
    for grade in grades:
        if isinstance(grade, list):
            for gr in grade:
                grades_list.append(int(gr))
        else:
            grades_list.append(int(grade))
    print('grades_list = ',grades_list)
    price = db.execute(f'SELECT * FROM "view_prices" WHERE "recipient_login" = \'{child_login}\';')
    gr_to_pr = {}
    for pr in price:
        gr_to_pr[pr[2]] = pr[3]
    summ = 0
    for grade in grades_list:
        summ += gr_to_pr[grade]
    return summ
