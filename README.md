# Image Processing API Golang

## Features 
    1. Png to Jpg Images Converter 
       * Endpoint : POST /png-to-jpeg
       * Request Form-Data : images []File
       * Response Form-Data : file Zip
    2. Resize Images
       * Endpoint : POST /resize-image
       * Request Form-Data : 
         - images []File
         - image_dimensions []string{`"width": string, "height": string`}
       * Response Form-Data : file Zip
    3. Compress Images
       * Endpoint : POST /compress-image
       * Request Form-Data : 
         - images []File
         - qualities []string{}
       * Response Form-Data : file Zip

## Installation 
    1. Install golang [Guide](https://go.dev/doc/install)
    2. Install OpenCV
       > Use guide from GoCV since we are using it as binding to use opencv API [GoCV](https://gocv.io/getting-started/)
    3. Install dependencies
       > go mod download && go mod verify
    3. Run the app 
       > go run ./cmd/main.go 
    4. Then access the api via port 3000
    
### Dependencies
    1. Gin - https://github.com/gin-gonic/gin
    2. Viper - https://github.com/spf13/viper
    3. GoCv - https://github.com/hybridgroup/gocv

## Tests
    * Run test.sh script or
    * Run command 
      > go test -timeout 300s github.com/dzulfiikar/image-processing-api-golang/cmd/tests

## Docker Deployment
    * Run command
      > docker compose up -d 
    * Then access the api via port 3000
    * Adjust config according to your preferences

## Postman Collection
    Import file "Golang Image Processing API.postman_collection.json" into your postman