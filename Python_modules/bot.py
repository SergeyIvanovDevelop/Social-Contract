import telebot 
from telebot import types
import schedule
from threading import Thread
from time import sleep
import requests
import re
from translater import get_summ
from gen_qr1_alter import create_qr
from gen_qr import qr_code_save
from config import Config

token=Config.get('bot', 'token')
keyboard1 = telebot.types.ReplyKeyboardMarkup()
keyboard1.row('Проверка связи', 'Авторизация')
keyboard1.row('Отправь мне QR')
bot=telebot.TeleBot(token)

COMMANDS = ['Проверка связи','Авторизация','Отправь мне QR']
IDS=[]  # id пользователей кто запустил бота
USERNAMES = []
LOGINS = dict()

# Запуск бота
@bot.message_handler(commands=['start'])
def start_message(message):
    bot.send_message(message.chat.id, 'Здравствуйте, Вы меня запустили!\nЯ буду напоминать Вам переводить деньги ребенку.\nПожалуйста, введите логин ребенка. Команда /login', reply_markup=keyboard1)

    """Добавляем id если его нет"""
    if message.chat.id not in IDS:
        IDS.append(message.chat.id)
        LOGINS[message.chat.id] = ""

    if message.from_user.username not in USERNAMES:
        USERNAMES.append(message.from_user.username)

    """Запрос логина"""
    pass

@bot.message_handler(commands=['help'])
def author(message):
    text = """/start - Запуск бота
/help - Справка
/login - Авторизация
/qr - Получение qr-кода"""
    bot.send_message(message.chat.id, text)

# Авторизация
@bot.message_handler(commands=['login'])
def save_login(message):
    msg = bot.send_message(message.chat.id, "Введите свой логин")
    bot.register_next_step_handler(msg, user_answer)

def user_answer(message):
    print(message.text)
    if message.text in COMMANDS:
        bot.send_message(message.chat.id, "Ошибка авторизации, попробуйте снова. Команда /login")
    else:
        LOGINS[message.chat.id] = message.text
        bot.send_message(message.chat.id, "Авторизация прошла успешно!")

# Обработка сообщений
@bot.message_handler()
def handle_message(message):
    print(message.text)
    print(message.chat.id)
    print(message.from_user.username)

    """Добавляем id если его нет"""
    if message.chat.id not in IDS:
        IDS.append(message.chat.id)
        LOGINS[message.chat.id] = ""

    if message.from_user.username not in USERNAMES:
        USERNAMES.append(message.from_user.username)
    
    if message.text=="Проверка связи":
        bot.send_message(message.chat.id, f'Привет, {message.from_user.username}, я работаю, все ОК!', reply_markup=keyboard1)

    # base processing
    elif message.text =="Отправь мне QR": 
        if LOGINS[message.chat.id]:
            print("Login Exist")
            msg = bot.send_message(message.chat.id, 'Введите день месяц год через пробел. Пример:\n01 04 2020', reply_markup=keyboard1)
            day = bot.register_next_step_handler(msg, get_date)            
        else: 
            bot.send_message(message.chat.id, 'Сначала авторизуйтесь. Введите логин ребенка. Команда /login', reply_markup=keyboard1)
    # Login
    elif message.text == "Авторизация":
        msg = bot.send_message(message.chat.id, "Введите логин ребенка")
        bot.register_next_step_handler(msg, user_answer)
    else: 
        bot.send_message(message.chat.id, 'Не знаю такую команду', reply_markup=keyboard1)

def get_date(message):
    date = message.text
    print(date)
    # проверка корректности
    template = '[0-9]{2} [0-9]{2} [0-9]{4}'
    if not re.match(template, date):
        bot.send_message(message.chat.id, 'Неверный формат даты. Поробуйте еще раз.', reply_markup=keyboard1)
    else:
        m = date.split()
        day = m[2]+'-'+m[1]+'-'+m[0]
        print(day)
        bot.send_message(message.chat.id, 'Начинаю генерировать QR-код', reply_markup=keyboard1)
        try:
            summ = get_summ(LOGINS[message.chat.id], day )
            qr_code_save(LOGINS[message.chat.id], "Оценки" , summ, filename= f'{LOGINS[message.chat.id]}_qr.png')
        except Exception as e: 
            print(e)
            create_qr(LOGINS[message.chat.id]) # Заглушка, на всякий случай =)
        sleep(1)
        bot.send_photo(message.chat.id, open(f'{LOGINS[message.chat.id]}_qr.png', 'rb'))

"""Функции ежеминутной отправки сообщений в потоке"""
def schedule_checker():
    while True:
        schedule.run_pending()
        sleep(10)

def function_to_run():
    for id in IDS:
        return bot.send_message(id, "Проверьте QR код, может Ваш ребенок что-нибудь заработал)")

if __name__ == "__main__":  
    """Функции ежеминутной отправки сообщений в потоке"""
    schedule.every().minute.do(function_to_run)
    Thread(target=schedule_checker).start()   
    bot.polling()  # запуск бота
    # Отправка всем уведомления о выключении
    for id in IDS:
        bot.send_message(id, 'Я временно выключаюсь...')
    print(LOGINS)
