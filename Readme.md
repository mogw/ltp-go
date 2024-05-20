# API for the Last Traded Price (ltp) of bitcoin


## Prerequisites

- Go 1.16 or later
- Docker (optional)

## Building the App

### Without Docker

1. Clone the repository:

    ```bash
    git clone https://github.com/mogw/ltp-go.git
    ```

2. Build the app:

    ```bash
    cd ltp-go
    go build -o app
    ```

### With Docker

1. Clone the repository:

    ```bash
    git clone https://github.com/mogw/ltp-go.git
    ```

2. Build the Docker image:

    ```bash
    docker build -t your-image-name .
    ```

## Running the App

### Without Docker:

1. Make sure you have Go installed and the app is built (see "Building the App" section above).

2. Run the app:

    ```bash
    ./app
    ```

    The app will start listening on http://localhost:8080.

### With Docker

1. Make sure you have Docker installed (see "Prerequisites" section above).

2. Run the Docker container:

    ```bash
    docker run -p 8080:8080 your-image-name
    ```

    Replace your-image-name with the name of your Docker image. The app will start listening on http://localhost:8080.


## Testing the App

To run the tests, execute the following command:

```bash
go test
```