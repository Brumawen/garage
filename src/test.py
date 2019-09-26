import telepot
from telepot.loop import MessageLoop
from time import sleep

def handle(msg):
    print(msg)

bot = telepot.Bot('709856064:AAG9yLMnuOIRAKJmLssRG69R7GNJdNdYidA')
print (bot.getMe())

MessageLoop(bot, handle).run_as_thread()

while True:
    sleep(10)