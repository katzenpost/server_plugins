#!/usr/bin/env python

from concurrent import futures
import sys
import time

import grpc

import kaetzchen_pb2
import kaetzchen_pb2_grpc




class EchoServicer(kaetzchen_pb2_grpc.KaetzchenServicer):

    def __init__(self):
        pass

    def OnRequest(self, request, context):
        return request


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    kaetzchen_pb2_grpc.add_KaetzchenServicer_to_server(
        EchoServicer(), server)
    server.add_insecure_port('unix:/home/user/tmp/echo.sock')
    server.start()
    try:
        while True:
            time.sleep(60 * 5)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
