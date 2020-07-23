# Concurrent fetch cancelation

Example to demonstrate how to use `context.Context` to cancel/halt processes in goroutine.

## Example

Given following endpoint:

`GET https://jsonplaceholder.typicode.com/photos/{id}` (`{id}` can be 1-5000)

We will try to concurrently fetch these 5000 data (1 API call for 1 ID), but with a timeout (set via `timeout` URL query during call).

When the timeout exceeded, the background fetches will be canceled/halted and the already-fetched-data will be returned.

## Running

1. Clone this repo
2. Run `go get`
3. Run `go run *.go`, it will start on `localhost:8000`

You can start testing by making API call:

```
curl --request GET \
  --url 'http://localhost:8000/?timeout=1000'
```

The response should be as follow:

```json
{
  "cancelled": 4999,
  "data": [
    {
      "id": 188,
      "albumId": 4,
      "title": "quae accusamus voluptas aperiam est amet",
      "url": "https://via.placeholder.com/600/40bdc9",
      "thumbnailUrl": "https://via.placeholder.com/150/40bdc9"
    }
  ],
  "error": 0,
  "success": 1,
  "total": 5000
}
```

Explanation:
- `cancelled` indicates how many calls were canceled/halted
- `error` indicates how many calls were failed with error (mostly HTTP error)
- `success` speaks for itself
- `total` is the number of data stored in the data source which we try to fetch concurrently
- `data` contains the successfully fetched data 