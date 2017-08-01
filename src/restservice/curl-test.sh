echo "Putting data:"
curl -X PUT localhost:8001/entry/first/hello; echo
curl -X PUT localhost:8001/entry/second/hi; echo

echo "
Calling list:"
curl localhost:8001/list; echo

echo "
Calling specific entry:"
curl localhost:8001/entry/first; echo
curl localhost:8001/entry/second; echo

echo "
Calling the main webpage"
curl localhost:8001; echo
