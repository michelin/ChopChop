import http.server
from threading import Thread, Barrier
from functools import partial

class _WebMockServer(http.server.SimpleHTTPRequestHandler):
    def log_message(self, format, *args):
        pass # avoid logging

    def do_GET(self):
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(bytes("Example", "utf8"))

class MockServer(object):
    def __init__(self):
        self.__server = None
        self.__thread = None
    
    def start_mock_server(self, port: int):
        print("Starting server on port", port)
        self.__server = http.server.HTTPServer(("", port), _WebMockServer)

        def serve_forever(server):
            with server:
                server.serve_forever()

        self.__thread = Thread(target=serve_forever, args=(self.__server, ))
        self.__thread.setDaemon(True)
        self.__thread.start()

    def stop_mock_server(self):
        print("Stopping server")
        if self.__server is not None:
            self.__server.shutdown()
            self.__thread.join()

            self.__server = None
            self.__thread = None
