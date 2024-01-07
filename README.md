ffufClean
==

cleans up false positive results from route handlers that match based on `startswith` and keeps the shortest match

e.g.
```
/admin.zip
/admin.tar
/admin
/admin.xml
/adminpanel
```

```bash
cat results.json | jq -c '.results[]' | ./ffufClean | jq -c '.[]'
```