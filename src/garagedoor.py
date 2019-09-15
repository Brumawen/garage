from gpiozero import Button
from gpiozero import LED
import argparse
import os
import requests
import logging

# Command line arguments
parser = argparse.ArgumentParser(description='Control a Garagae Door')
parser.add_argument('-name', type=str, required=True, help='The name of the Garage Door')
parser.add_argument('-switch', type=int, required=True, help='The pin number for the Garage Door closed sensor switch')
parser.add_argument('-led', type=int, required=True, help='The pin number for the Garage Door LED')
parser.add_argument('-log', type=str, required=True, help='The log file')
args = parser.parse_args()

# Create a custom logger
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

# Create handlers
f_handler = logging.FileHandler(args.log)
f_handler.setLevel(logging.DEBUG)

# Create formatters and add it to handlers
f_format = logging.Formatter('%(asctime)s - %(levelname)s - %(message)s')
f_handler.setFormatter(f_format)

# Add handlers to the logger
logger.addHandler(f_handler)

currDir = os.path.dirname(os.path.abspath(__file__))
logger.debug("Current directory is '" + currDir + "'")
dataDir = os.path.join(currDir, "data")
if not os.path.exists(dataDir):
    os.mkdir(dataDir)
    logger.debug("Created directory '" + dataDir + "'")
updFile = os.path.join(dataDir, args.name + ".state")

led = LED(args.led)
pwrLed = LED(9)
switch = Button(args.switch)

logger.debug("Turning on power led")
pwrLed.on()

def send_update_to_service(post):
    try:
        f = open(updFile,"w+")
        if switch.is_pressed:
            logger.debug("Setting '" + updFile + "' to closed")
            f.write("closed")
        else:
            logger.debug("Setting '" + updFile + "' to open")
            f.write("open")
        f.close()
    except Exception as ex:
        logger.error("Exception writing door state to file. " + str(ex))
        return
    
    if not post:
        return

    try:
        requests.post("http://localhost:20515/room/update")
    except Exception as ex:
        logger.error("Exception calling update web method. " + str(ex))

logger.debug("Controlling door '" + args.name + "'")
send_update_to_service(False)

try:
    while True:
        if switch.is_pressed:
            led.on()
            send_update_to_service(True)
            switch.wait_for_release()
            led.off()
        else:
            led.off()
            send_update_to_service(True)
            switch.wait_for_press()
            led.on()
        
except Exception as ex:
    logger.error("Exception running. " + str(ex))

logger.debug("Turning off power led")
pwrLed.off()
