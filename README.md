# go-worker-template
New go worker template apply with gRPC

## Example consume job
```
{
    "id": "123456",
    "text": "Test text",
    "task1": null,
    "task2": null,
    "task3": null
}
```

## Example publish job
```
{
    "id": "123456",
    "text": "Test text",
    "task1": {
        "processed_text": "Test text",
        "entities": [
            {
                "start": 1,
                "end": 2,
                "entity": "Hospital",
                "label": "PLACE"
            },
            {
                "start": 3,
                "end": 10,
                "entity": "Doctor",
                "label": "PERSON"
            }
        ]
    },
    "task2": "POSITIVE",
    "task3": 100
}
```