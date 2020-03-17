# Schematic
High performance message collection and shema managment

## Infrastructure

[![N|Solid](diagram.png)](diagram.png)

## Notes

- Collector is a high performance generic message collection service. It receives messages from a watch through a client-streaming grpc connection. 
- Schematic Layer registers and maintains message schemas. It allows services to decode unknown messages at run time and verify their integrity.
- ConsumerPy is using generated code to consume the messages.
- ConsumeGo is not aware of the schema of consumed messages. It is able to deserialize and manipulate them by fetching their file descriptor from the schematic layer.


## Demo
![](demo.gif)

> left: consumerGo, middle: watchExporter App, right: consumerPy

### Flow

1) Apple Watch app establishes client-sctreaming grpc connection with collector.
2) Collector, consumer, and schematic layer establish connection with the kafka broker.
3) Biometric message schema (file descriptor) is registered with the schematic layer.
4) ConsumeGo fetches biometric schema.
5) Apple Watch app begins streaming biometric data.
6) ConsumerPy receives messages and deserializes them with pre-generated code.
7) ConsumerGo receives messages and deserializes them with previously fetched schema.

### Stress Test

Test Details: 
1) kafka: 1 partition, 2 replicas, 
2) 1 collector
3) 0.5 sec collector flush rate
4) collector waits for leader acknowledgement
5) each kafka broker sits on independent n1-standard node
6) there are 10 clients, each on their independent n1-standard-2, each with up to 2500 connections, thus in total there are 25 000 connections to collector
7) collector runs on n1-standard-4  (4 vCPU, 15 GB memory)  and n1-standard-8  (8 vCPU, 30 GB memory)  
8) time is measured to collect 10 million ~126-254 byte message

#### Resuls


[![N|Solid](stress_results.png)](stress_results.png)

#### Raw stress test data


1 collector running on n1-standard-4 (GKE)


| Concurrent Connections | 2.5e2         | 2.5e3         | 2.5e4         |
| ---------------------- | ------------- | ------------- | ------------- |
| Test completion times  | 45.2117       | 53.2022       | 63.6975       |
| (seconds)              | 44.8457       | 52.2528       | 62.9623       |
|                        | 45.3871       | 51.9387       | 61.9285       |
|                        | 46.3244       | 52.8396       | 63.6561       |
| Average                | 45.442225     | 52.558325     | 63.0611       |

1 collector running on n1-standard-8 (GKE)


| Concurrent Connections | 2.5e2         | 2.5e3         | 2.5e4         |
| ---------------------- | ------------- | ------------- | ------------- |
| Test completion times  | 31.8510       | 32.1720       | 38.1253       |
| (seconds)              | 30.3364       | 33.3108       | 36.5358       |
|                        | 30.2292       | 32.1427       | 38.3180       |
|                        | 30.5378       | 32.3915       | 36.7096       |
| Average                | 30.7386       | 32.5042       | 37.4221       |


    
