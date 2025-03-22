# Example Temporal io Workflow

## Sample code

This shows how to use a temporalio code

## Admin code

```bash
temporal operator namespace create -n default # or another namespace
temporal operator namespace update --visibility-archival-state enabled -n default # or another namespace
temporal operator namespace update --history-archival-state enabled default # or another namespace
```
