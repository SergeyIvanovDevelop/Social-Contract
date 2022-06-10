import qrcode
from database import Database as DB

db = DB()

def qr_code_save(login: str, purpuse: str, summ: int, filename='qr.png'):
    for card in db.select('view_cards', ('login','Name','PersonalAcc','BankName','BIC','CorrespAcc','PayeeINN','KPP')):
        if card[0] == login:
            img = qrcode.make(f'ST00012|Name={card[1]}|PersonalAcc={card[2]}|BankName={card[3]}|BIC={card[4]}|CorrespAcc={card[5]}|PayeeINN={card[6]}|KPP={card[7]}|Purpuse={purpuse}|PaymPeriod=1221|Sum={summ}|')
            img.save(filename)
            break
