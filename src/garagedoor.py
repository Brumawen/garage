from gpiozero import Button
from gpiozero import LED
import argparse
import os
import requests

# Command line arguments
parser = argparse.ArgumentParser(description='Control a Garagae Door')
parser.add_argument('-n', type=str, required=True, help='The name of the Garage Door')
parser.add_argument('-d', type=int, required=True, help='The pin number for the Garage Door closed sensor switch')
parser.add_argument('-l', type=int, required=True, help='The pin number for the Garage Door LED')
args = parser.parse_args()

currDir = os.path.dirname(os.path.abspath(__file__))
print(currDir)
dataDir = os.path.join(currDir, "data")
if not os.path.exists(dataDir):
    os.mkdir(dataDir)
    print("Created directory '" + dataDir + "'")
updFile = os.path.join(dataDir, args.n + ".state")

led = LED(args.l)
switch = Button(args.d)

def send_update_to_service():
    try:
        f = open(updFile,"w+")
        if switch.is_pressed:
            f.write("closed")
        else:
            f.write("open")
        f.close()
    except Exception as ex:
        print('Exception writing door state to file.', ex)
        return
    
    try:
        requests.post("http://localhost:20515/room/update")
    except Exception as ex:
        print('Exception calling update web method.', ex)

print ("Controlling door '" + args.n + "'")

while True:
    if switch.is_pressed:
        led.on()
        switch.wait_for_release()
        led.off()
    else:
        led.off()
        switch.wait_for_press()
        led.on()
    send_update_to_service()

