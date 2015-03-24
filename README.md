# ninjasphere-limitlessled
Limitlessled driver for NinjaSphere


**How to start the driver:**

1. Download the sources and do **'make target'**. This will generate **'driver-limitlessled'** folder with driver binary compatible with sphere architecture along with some other necessary files

2. Make sure these directories are available in the sphere --> **/data/sphere/user-autostart/{drivers,apps}**. If not --> **sudo mkdir -p /data/sphere/user-autostart/{drivers,apps} && sudo chown -R ninja.ninja /data/sphere** in the sphere
3. Change the ip field in milight.xml present in the 'driver-limitlessled' folder to your bridge ip

3. **scp -r driver-limitlessled ninja@ninjasphere.local:/data/sphere/user-autostart/drivers**

4. reboot the sphere or ssh into sphere & **nservice driver-limitlessled start**
