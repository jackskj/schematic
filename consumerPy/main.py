import threading
import signal
import sys
from kafka import KafkaConsumer
from pb import biometrics_pb2
from google.protobuf.json_format import MessageToJson


class WatchConsumer(threading.Thread):
    daemon = True
    broker = 'localhost:9092',
    topic = "apple_watch"
    group = "consumer_py"
    offset = 'latest'
    consumer = None

    def __init__(self):
        super(WatchConsumer, self).__init__()
        signal.signal(signal.SIGTERM, self.terminate)

    def run(self):
        self.consumer = KafkaConsumer(
            self.topic,
            bootstrap_servers=self.broker,
            auto_offset_reset=self.offset,
            enable_auto_commit=True,
            group_id=self.group
        )
        for message in self.consumer:
            msg = biometrics_pb2.Biometrics()
            msg.ParseFromString(message.value)
            print(MessageToJson(msg))

    #  terminate gracefully with kill,
    #  let kafka broker know that this consumer is no longer listening
    def terminate(self, sigterm, frame):
        self.consumer.close()
        #  exit main thread
        print("Exiting thread with sigterm")
        sys.exit()


def main():
    consumer = WatchConsumer()
    consumer.start()
    consumer.join()


if __name__ == "__main__":
    main()
