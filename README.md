# s3-explorer

A simple "file explorer" API, serving an S3 bucket as if it were a filesystem.

Responses are json objects:

- for files: the content, string encoded
  - for "directories" (keys with a trailing slash): a list of the `children`

### Configuration

Environment variables:

- PORT: the port to serve the HTTP api
  - default: 8080
- ENV: `"development"` or `"production"`
  - when "production", configures gin's "release mode"
  - default: "development"
- FILESYSTEM: "local" or "s3"
  - default: "local"
- BUCKET_NAME: s3 bucket name (required when FILESYSTEM="s3")
- BUCKET_REGION: s3 bucket region (required when FILESYSTEM="s3")
