# ninjasphere-limitlessled
Limitlessled driver for NinjaSphere


**How to start the driver **

1. Download the source and do **'make target'**. This will generate **'driver-limitlessled'** folder

2. Make sure these directories are available --> **/data/sphere/user-autostart/{drivers,apps}**. If not --> **sudo mkdir -p /data/sphere/user-autostart/{drivers,apps} && sudo chown -R ninja.ninja /data/sphere**

3. **scp -r driver-limitlessled ninjasphere.local:/data/sphere/user-autostart/drivers**

4. reboot the sphere or ssh into sphere & cd into **/data/sphere/user-autostart/drivers/driver-limitlessled** folder and ** nservice driver-limitlessled start **
