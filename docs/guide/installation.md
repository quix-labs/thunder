# Installation

You have several options for installing **Thunder**:


## Using Docker

1. Pull the Docker image from the GitHub Container Registry:

    ```bash
    docker pull ghcr.io/quix-labs/thunder:latest
    ```

2. Start the container:

    ```bash
    touch config.json
    docker run --rm -p "3000:3000" -v "./config.json:/config.json" --name thunder ghcr.io/quix-labs/thunder:latest
    ```

3. Access the tool at `http://localhost:3000` and configure it via the web interface.

:::tip NOTE
It is recommended to specify a tag (version) when pulling the Docker image, for example:

```bash
docker pull ghcr.io/quix-labs/thunder:v1.0.0
```

This ensures you're using a specific, stable version of Thunder rather than the latest version, which may be in
development.
:::

## Using Prebuilt Assets

1. Download the appropriate package from the [Releases page](https://github.com/quix-labs/thunder/releases).

2. Follow the installation instructions provided for your platform (e.g., `.deb`, `.apk`, etc.).

## Building from Source

1. Set up your Go project:

    ```bash
    cd /path/to/your/project
    go mod init your_app_name
    go get -u github.com/quix-labs/thunder/app@main
    ```

2. Create a `main.go` file and add the following code:

    ```go
    package main
    
    import (
        "github.com/quix-labs/thunder"
        _ "github.com/quix-labs/thunder/exporters/csv"
        _ "github.com/quix-labs/thunder/exporters/json"
        _ "github.com/quix-labs/thunder/exporters/yaml"
        _ "github.com/quix-labs/thunder/modules/api"
        _ "github.com/quix-labs/thunder/modules/frontend"
        _ "github.com/quix-labs/thunder/modules/http_server"
        _ "github.com/quix-labs/thunder/source-drivers/postgresql_flash"
        _ "github.com/quix-labs/thunder/target-drivers/elastic"
    )
    
    func main() {
        err := thunder.Start()
        if err != nil {
            panic(err)
        }
    }
    ```

3. Compile the app:

    ```bash
    go mod tidy
    CGO_ENABLED=0 go build -ldflags="-s -w" -o your_app_name
    ```

4. Run your app:

    ```bash
    ./your_app_name
    ```

5. The configuration file (`config.json`) will be placed in the current directory. Access the tool at
   `http://localhost:3000` to configure it.

## Next Steps

Once you have Thunder up and running, follow these steps:

1. **[Configure Sources](./sources)**
2. **[Define Targets](./targets)**
3. **[Set Up Processors](./processors)**
