# bambulab-exporter

Prometheus exporter for [bambulab](https://bambulab.com) 3D printers

## Getting Started

You can run from a prebuilt docker container `evanofslack/bambulab-exporter:latest`

### Credentials

#### Cloud

If connecting to bambulab cloud, you will need to know your user-id (`uid`) and access token. For more information about see [this](https://github.com/Doridian/OpenBambuAPI/blob/main/cloud-http.md):

For a one liner to get your access token (substitute in your bambulab cloud username and password):
```bash
curl -v -X POST -H 'Content-Type: application/json' -d '{"account":"YOUR_USER_NAME","password":"YOUR_PASSWORD"}' https://bambulab.com/api/sign-in/form 2>&1 | grep token= | awk '{print$3}'
```
This should be provided as `BAMBU_PASSWORD`.

Then to obtain your user-id:
```bash
curl -X GET https://api.bambulab.com/v1/iot-service/api/user/project -H "Authorization: Bearer YOUR_TOKEN"
```
In the returned json, look for `user_id` entry which is a number, then your user-id is just the prefix `u_` followed by that number. This should be provided as `BAMBU_USERNAME` env var.

### Local

If you printer is running in local mode, `BAMBU_ENDPOINT` will be the printer's IP address, `BAMBU_USERNAME` is `bblp` by default, and `BAMBU_PASSWORD` is printer password (found on printer under network settings). In both modes `BAMBU_DEVICE_ID` is your printer's serial number. 

#### Container

```bash
services:
  bambulab-exporter:
    image: evanofslack/bambulab-exporter:latest
    container_name: bambulab-exporter
    restart: unless-stopped
    ports:
      - 9091:9091
    environment:
      HTTP_PORT:"9091" # port metrics server served from
      LOG_LEVEL:"info" # exporter logs, info or debug
      BAMBU_DEVICE_ID:"serial_number" # serial number of printer
      BAMBU_ENDPOINT:"us.mqtt.bambulab.com" # connect to bambulab cloud mqtt server (printer ip address for local mode)
      BAMBU_USERNAME:"u_0000000" # bambulab user-id (bblp for local mode)
      BAMBU_PASSWORD:"token" # bambulab access token (printer code for local mode)
```

## Metrics

This is a sample of the metrics available

```
# HELP bambulab_ams_enabled_state
# TYPE bambulab_ams_enabled_state gauge
bambulab_ams_enabled_state{device="01P09C461602411"} 0
# HELP bambulab_ams_humidity
# TYPE bambulab_ams_humidity gauge
bambulab_ams_humidity{device="01P09C461602411",unit="0"} 5
# HELP bambulab_ams_powered_state
# TYPE bambulab_ams_powered_state gauge
bambulab_ams_powered_state{device="01P09C461602411"} 0
# HELP bambulab_ams_temperature
# TYPE bambulab_ams_temperature gauge
bambulab_ams_temperature{device="01P09C461602411",unit="0"} 0
# HELP bambulab_camera_enabled_state
# TYPE bambulab_camera_enabled_state gauge
bambulab_camera_enabled_state{device="01P09C461602411"} 0
# HELP bambulab_camera_timelapse_state
# TYPE bambulab_camera_timelapse_state gauge
bambulab_camera_timelapse_state{device="01P09C461602411"} 0
# HELP bambulab_chamber_light_state
# TYPE bambulab_chamber_light_state gauge
bambulab_chamber_light_state{device="01P09C461602411"} 1
# HELP bambulab_fan_speed_percent
# TYPE bambulab_fan_speed_percent gauge
bambulab_fan_speed_percent{device="01P09C461602411",fan="auxiliary"} 0
bambulab_fan_speed_percent{device="01P09C461602411",fan="chamber"} 0
bambulab_fan_speed_percent{device="01P09C461602411",fan="hotend"} 0
bambulab_fan_speed_percent{device="01P09C461602411",fan="part"} 0
# HELP bambulab_gcode_state
# TYPE bambulab_gcode_state gauge
bambulab_gcode_state{device="01P09C461602411",state="FINISH"} 1
# HELP bambulab_layer_number
# TYPE bambulab_layer_number gauge
bambulab_layer_number{device="01P09C461602411"} 99
# HELP bambulab_layer_number_target
# TYPE bambulab_layer_number_target gauge
bambulab_layer_number_target{device="01P09C461602411"} 99
# HELP bambulab_nozzle_diameter
# TYPE bambulab_nozzle_diameter gauge
bambulab_nozzle_diameter{device="01P09C461602411"} 0.4
# HELP bambulab_nozzle_speed_level
# TYPE bambulab_nozzle_speed_level gauge
bambulab_nozzle_speed_level{device="01P09C461602411",level="standard"} 1
# HELP bambulab_nozzle_speed_magnitude
# TYPE bambulab_nozzle_speed_magnitude gauge
bambulab_nozzle_speed_magnitude{device="01P09C461602411"} 100
# HELP bambulab_nozzle_temperature
# TYPE bambulab_nozzle_temperature gauge
bambulab_nozzle_temperature{device="01P09C461602411"} 34.15625
# HELP bambulab_nozzle_temperature_target
# TYPE bambulab_nozzle_temperature_target gauge
bambulab_nozzle_temperature_target{device="01P09C461602411"} 0
# HELP bambulab_nozzle_type_state
# TYPE bambulab_nozzle_type_state gauge
bambulab_nozzle_type_state{device="01P09C461602411",type="stainless_steel"} 1
# HELP bambulab_print_percent
# TYPE bambulab_print_percent gauge
bambulab_print_percent{device="01P09C461602411",model=""} 100
# HELP bambulab_print_time_remaining_minutes
# TYPE bambulab_print_time_remaining_minutes gauge
bambulab_print_time_remaining_minutes{device="01P09C461602411",model=""} 0
# HELP bambulab_prints_total
# TYPE bambulab_prints_total counter
bambulab_prints_total{device="01P09C461602411",result="finish"} 1
# HELP bambulab_wifi_signal
# TYPE bambulab_wifi_signal gauge
bambulab_wifi_signal{device="01P09C461602411"} -63
```
