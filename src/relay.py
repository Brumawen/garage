import sys
from gpiozero import LED
from time import sleep

if len(sys.argv) == 1:
    print("Relay GPIO pin not included.")
else:    
    led = LED(sys.argv[1])
    led.off()
    sleep(0.5)
    led.off()
