# Paralisus API

Base URL: `http://localhost:36789`

---

## GET /api/class/getClassInfo

Returns remaining slot information for one or more classes.

**Query Parameters**

| Parameter | Type   | Required | Description                          |
|-----------|--------|----------|--------------------------------------|
| `ids`     | string | Yes      | Comma-separated list of class IDs    |

**Single class**
```
GET /api/class/getClassInfo?ids=CS101
```
```json
{
  "CS101": { "remain_slot": "30" }
}
```

**Multiple classes** — pass a comma-separated list, no spaces
```
GET /api/class/getClassInfo?ids=161001,161002
```
```json
{
  "161001": { "remain_slot": " 40" },
  "161002": { "remain_slot": " 40" }
}
```

> Response keys match the requested IDs. If an ID has no data in Redis, it is omitted from the response.

**Error Responses**

| Status | Description                        |
|--------|------------------------------------|
| `400`  | `missing ids` — `ids` param absent |
| `500`  | Redis script execution failed      |

---

## POST /api/class/register

Register a student for a class.

> ⚠️ Not yet implemented.

---

## POST /api/class/unregister

Unregister a student from a class.

> ⚠️ Not yet implemented.

---

## Architecture Notes

Requests to `/register` and `/unregister` are processed asynchronously via the **IngressRouter** event bus:

- **ResponseBus** — updates remaining slot in Redis (high consistency)
- **LoggingBus** — persists registration status to MongoDB

Each bus runs a pool of 4 workers. Workers drain all pending jobs before shutdown.

### Event Payload

```json
{
  "type": "register | unregister",
  "request": {
    "origin_id": "<request-id>",
    "payload": { }
  }
}
```
