# DICOM Viewer Service

This repository contains a simple microservice to upload, view and search DICOM files. <br> <br>

## Getting Started

### Running Locally

-   Pre-requisites:
    1. Go 1.22+ installed locally
-   To run:
    ```
    go run ./cmd/dicomviewer
    ```
-   To build
    ```
    go build ./cmd/dicomviewer
    ```

### Running With Docker

-   Pre-requisites:
    1. Docker installed and running. See https://docs.docker.com/get-started/ for help
-   To build image:
    ```
    docker build --tag "dicomviewer" .
    ```
-   To run container:
    ```
    docker run -it -p 3000:3000 "dicomviewer"
    ```
    _By default, the service listens on port 3000. This command binds localhost:3000 to the
    container port for easy access_

<br>

### Specifying a custom port

-   Locally:
    ```
    go run ./cmd/dicomviewer -port=4000
    ```
-   Docker:

    ```
    docker run -it -p 4000:4000 "dicomviewer" -port=4000
    ```

## Service API

1. Create a DICOM files

    - Request:

        ```
        Content-Type: multipart/form-data;
        Content-Length: ...;
        Content-Disposition: form-data; name="file"; filename="<filename>"

        POST /api/v1/files
        ```

    - Response:

        ```
        Content-Type: application/json

        {
            fileId: "<new file id>"
        }
        ```

2. List all available DICOM files

    - Request:
        ```
        GET /api/v1/files
        ```
    - Response:

        ```
        Content-Type: application/json

        {
            fileIds: [
                "<fileId>",
                ...
            ]
        }

        ```

3. Retrieve a raw DICOM file
    - Request:
        ```
        GET /api/v1/files/<fileId>
        ```
    - Response:
        ```
        Content-Type: application/octet-stream
        ```
4. Retrieve a DICOM file as a PNG
    - Request:
        ```
        GET /api/v1/files/<fileId>/png?remap=true
        ```
        - remap (default _true_): optionally remap image pixel values from their default range to
          0-255. This can result in better image quality in some cases
    - Response:
        ```
        Content-Type: image/png
        Transfer-Encoding: chunked
        ```
5. Search a DICOM file's elements/attributes by tag

    - Request:

        ```
        GET /api/v1/files/<fileId>/attributes?tag=(0001,0002)&tag=(0002,0004)
        ```

        - tag: a DICOM attribute tag in the form `(<group>, <element>)`. The query can contain many
          tags

    - Response:

        ```
        Content-Type: application/json

        {
            elementsByTag: {
                "<tag>": {
                    "tag": {
                        "Group": "<group>",
                        "Element: "<element>",
                    },
                    "VR": <value representation>,
                    "rawVR": "<raw value representation>",
                    "valueLength": <value length>,
                    "value": <attribute value>
                }
            }
        }
        ```

## Coming Soon

-   Better logging
-   Telemetry, Metrics, Traces
-   Authz + Authn
-   More and better tests
-   Pagination for `GET /api/v1/files/<fileId>/attributes` and `GET /api/v1/files`
-   More server configuration options
-   Proper persistence adapters for DICOM files
-   Style cleanup here and there
