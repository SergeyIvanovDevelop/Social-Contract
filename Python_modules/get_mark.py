#!/usr/bin/env python3
from argparse import ArgumentParser, FileType
import requests

from lxml import etree
from time import sleep
from flask import Flask, jsonify
from typing import Dict, List
import json
from bs4 import BeautifulSoup

LOGIN_PAGE = "https://login.dnevnik.ru/login"
SHCOOL_PAGE = "https://schools.dnevnik.ru/school.aspx"

def parse_page_marks(page) -> Dict[str, List[int]]:
    soup = BeautifulSoup(page, 'lxml')
    table_collections = soup.find_all(['table'])
    d = {}
    for table in table_collections:
        for tr in table.find_all('tr'):
            td_collection = tr.find_all('td')
            if (td_collection[2].text != '\xa0'):
                mark = td_collection[2].text
                lesson = td_collection[0].text.split(' ')[0]
                if (d.get(lesson) == None):
                    d[lesson] = [mark]
                else:
                    d[lesson].append(mark)
    return d

app = Flask(__name__)
@app.route("/<username>")
@app.route("/<username>/<password>")
@app.route("/<username>/<password>/<year>")
@app.route("/<username>/<password>/<year>/<mouth>")
@app.route("/<username>/<password>/<year>/<mouth>/<day>")

def get_dump_json(username, password, year, mouth, day):
    marks = parse_page_marks(get_mark(year=year, mounth=mouth, day=day, login=username, password=password))
    return jsonify(marks)

def get_school_id(s):
    url = s.get(SHCOOL_PAGE).url
    print('url = ', url)
    return url.split("school=", 1)[1]

def get_mark(year, mounth, day, login, password):
    print('year = ', year)
    print('mounth = ', mounth)
    print('day = ', day)
    print('login = ', login)
    print('password = ', password)
    data = {
        "exceededAttempts": "False",
        "ReturnUrl": "",
        "login": login,
        "password": password
    }
    s = requests.Session()
    s.post(LOGIN_PAGE, data)
    sleep(1)
    school_id = get_school_id(s)
    sleep(1)
    MARKS_PAGE = f"https://schools.dnevnik.ru/marks.aspx?school={school_id}&index=0&tab=week&year={year}&month={mounth}&day={day}&homebasededucation=False"
    s.headers.update({'referer': "https://schools.dnevnik.ru/marks.aspx?school=49152&tab=week"})
    r = s.get(MARKS_PAGE)
    print('r.text = ', r.text)
    return r.text
