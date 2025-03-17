## How to run on development

1. Clone the repository
2. Install dependencies using `go install`
3. Create a `.env` file in the root directory and add the values following `.env.example`
4. Run the development server using `go run main.go`

## How to build docker image

1. Run `docker build -t IMAGENAME .` to build the image

### Where to use

- This image will be used in the `scn-deployment` repository to deploy the `tusd-server`
