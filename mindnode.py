# This script controls the daily task automation I use for mindnode
import os
import calendar
import shutil
from datetime import datetime

daily_directory_path = '/Users/matthewgleich/Library/Mobile Documents/W6L39UYL6Z~com~mindnode~MindNode/Documents/personal/daily'
# Mindnode files are folders (in package form)
template_folder_path = daily_directory_path + '/template.mindnode'

# Creating month folder
today = datetime.now()
month_name = (calendar.month_name[today.month]).lower()
os.chdir(daily_directory_path)
if not os.path.exists(month_name) and not os.path.isdir(month_name):
    os.mkdir(month_name)
    print(f'Created the folder {month_name}')
os.chdir(month_name)

# Creating file from template
new_folder_path = os.getcwd() + '/{}-{}-{}.mindnode'.format(
    month_name.title(), today.day, today.year)
if os.path.exists(new_folder_path):
    print('File already exists!')
    exit(1)
shutil.copytree(template_folder_path, new_folder_path)
print('Created the file')

os.system('open ' + new_folder_path.replace(' ', r'\ '))
