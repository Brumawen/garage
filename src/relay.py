import sys
from gpiozero import LED
from time import sleep

if len(sys.argv) == 1:
    print("Relay GPIO pin not included.")
else:    
    ledNo = 0
    if sys.argv[1] == "1":
        ledNo = 27
    elif sys.argv[1] == "2":
        ledNo = 17
    led = LED(ledNo)
    led.off()
    sleep(0.25)
    led.off()
