from mitmproxy import http
from time import sleep

def request(flow: http.HTTPFlow) -> None:
    print(f"URL: {flow.request.url}/{flow.request.path}")
    print(f"Method: {flow.request.method}")
    if flow.request.method == "POST" and \
       flow.request.path == "/testsuite/3/test/7":
        print("Intercepted DELETE request:")
        print(f"URL 1: {flow.request.url}")
        print(f"Headers 1: {flow.request.headers}")
        print(f"Method 1: {flow.request.method}")
        sleep(120)
        # Drop the request
        # flow.kill()