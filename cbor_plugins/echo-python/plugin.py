#!/usr/bin/env python

import sys
import time
import os
import os.path
import random
import string
import logging
import argparse
from http.server import BaseHTTPRequestHandler, HTTPServer
import socketserver


class UnixSocketHttpServer(socketserver.UnixStreamServer):
    def get_request(self):
        request, client_address = super(UnixSocketHttpServer, self).get_request()
        return (request, ["local", 0])

class NoSURBException(Exception):
    """
    This exception is raised when a request
    is received that does not contain a SURB.
    """

class EchoServicer():

    def __init__(self, logger):
        self.logger = logger

    def OnRequest(self, request):
        if not request.HasSURB:
            self.logger.error("error, request %s without SURB" % request.ID)
            raise NoSURBException
        self.logger.info("received request ID %s" % request.ID)
        return request.Payload

    def Parameters(self, empty):
        params = {
            "name":"python_echo_server",
            "version":"0.0.0",
        }
        return params


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("-l", required=True, help="log directory")
    args = vars(ap.parse_args())
    log_dir = args["l"]

    if not os.path.exists(log_dir) or not os.path.isdir(log_dir):
        print("log dir doesn't exist or is not a directory")
        os.exit(1)

    # setup logging
    logger = logging.getLogger('echo-python')
    logger.setLevel(logging.DEBUG)
    log_path = os.path.join(log_dir, "echo_python_%s.log" % os.getpid())
    fh = logging.FileHandler(log_path)
    fh.setLevel(logging.DEBUG)
    logger.addHandler(fh)
    logger.setLevel(logging.DEBUG)

    # start service
    logger.info("starting echo-python service")
    rand = ''.join(random.choice(string.digits) for _ in range(10))
    socket_file = "/tmp/pyecho_plugin_%s.sock" % rand

    print("%s\n" % socket_file)
    sys.stdout.flush()

    server = UnixSocketHttpServer((sock_file), EchoHttpHandler)
    server.serve_forever()


if __name__ == '__main__':
    main()
