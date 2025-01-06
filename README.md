# ffmate

## The Rest API for FFMpeg with queue support

**ffmate** is a wrapper for **ffmpeg** that adds a Rest API, a queue, webhooks and other features.

ffmates **featureset** (and roadmap) looks like this:

- [x] Rest API
    - [x] Tasks
    - [x] Webhooks
    - [x] Swagger documentation
- [x] Tasks
    - [x] Status (Queue, Running, Successful/Error/Canceled)
    - [x] Progress
    - [x] Wildcards
- [x] Queue
    - [x] Max concurrent tasks
- [x] Webhooks
    - [x] Create task
    - [x] Update task status
- [x] Update service    
- [ ] Presets
- [ ] Dashboard
- [ ] and more..

If you feel like the above description fits your needs, **welcome**! 

## Quick start guide

ffmate comes with a documented commandline.
Using **ffmate server** will start the server on the default port, 3000. This port can be overriden using the **-p** flag.

Use the API to add new Tasks. A Task includes the ffmpeg command to run together with the input and output file. \
One can make use of various wildcards:
- ${INPUT_FILE}
- ${OUTPUT_FILE}
- ${DATE_YEAR}
- ${DATE_SHORTYEAR}
- ${DATE_MONTH}
- ${DATE_DAY}
- ${TIME_HOUR}
- ${TIME_MINUTE}
- ${TIME_SECOND}
- ${TIMESTAMP_SECONDS}
- ${TIMESTAMP_MILLISECONDS}
- ${TIMESTAMP_MICROSECONDS}
- ${TIMESTAMP_NANOSECONDS}

## Swagger / Openapi

The swagger documentation can be found here **http://localhost:3000/swagger/index.html** once the server is running.

## Contributing

Thank you for considering contributing to ffmate! A contribution guide will be released in the future.

## License

ffmate is MIT Licensed. Copyright © 2025 by We love media
