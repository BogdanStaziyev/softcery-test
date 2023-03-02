## 👋  Summary

The goal is to create an HTTP API for uploading, optimizing, and serving images.

1. The API should expose an endpoint for uploading an image. After uploading, the image should be sent for optimization
   via a queue (e.g. RabbitMQ) to prevent excessive system load in case of many parallel image uploads and increase the
   system durability.
2. Uploaded images should be taken from the queue one by one and optimized using the `github.com/h2non/bimg` go
   package (or `github.com/nfnt/resize` package). For each original image, three smaller-size image variants should be
   generated and saved, with 75%, 50%, and 25% quality.
3. The API should expose an endpoint for downloading an image by ID. The endpoint should allow specifying the image
   optimization level using query parameters (e.g. `?quality=100/75/50/25`).

## 🤔 Evaluation criteria

1. **Functionality.** The developed solution should function as described in the "Summary" section. However, if you
   think that you can create a solution better than described in the "Summary" section, you are welcome to do so.
2. **Code simplicity**. The architecture should be simple and easy to understand, the code should be well-formatted and
   consistent. Usage of code formatters (like gofmt) and linters (like golangci-lint) is encouraged.

## **Getting Started**

1. **Run** Docker-compose file
2. **Use** Postman endpoints to test the functionality of the application

## Test Routes

- **Postman** - You can use Postman. I have uploaded the collection to a folder with the same name.

**Endpoints:**

- Upload image:
- POST: http://localhost:8080/api/v1/image/upload 
- also you will need to use a form file to send the image along with the
  request: key = image, value File from you device


- Download image:
- GET: http://localhost:8080/api/v1/image/download?id=1&quantity=100